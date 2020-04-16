package main

import (
	"context"
	"net/http"
	"os"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-legacy-redirector/config"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

var redir = http.StatusTemporaryRedirect
var landingPage = "https://www.ons.gov.uk/help/localstatistics"
var apiResponse = "This service is no longer available. Please visit https://www.ons.gov.uk/help/localstatistics for more information."
var visualResponse = "The article you have requested is no longer available."

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = "dp-legacy-redirector"
	ctx := context.Background()

	cfg, err := config.Get()
	if err != nil {
		log.Event(nil, "unable to retrieve service configuration", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	log.Event(ctx, "config on startup", log.INFO, log.Data{"config": cfg, "build_time": BuildTime, "git-commit": GitCommit})

	// Health check
	versionInfo, err := healthcheck.NewVersionInfo(BuildTime, GitCommit, Version)
	if err != nil {
		log.Event(ctx, "Failed to obtain VersionInfo for healthcheck", log.FATAL, log.Error(err))
		os.Exit(1)
	}
	hc := healthcheck.New(versionInfo, cfg.HealthckeckCriticalTimeout, cfg.HealthckeckInterval)

	srv := server.New(cfg.BindAddr, getRouter(hc))

	log.Event(ctx, "starting http server", log.INFO, log.Data{"bind_addr": cfg.BindAddr})
	if err := srv.ListenAndServe(); err != nil {
		log.Event(ctx, "error starting server", log.FATAL, log.Error(err))
		os.Exit(1)
	}
}

func getRouter(hc healthcheck.HealthCheck) *mux.Router {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", hc.Handler)

	// NeSS website
	router.Host("neighbourhood.statistics.gov.uk").Path("/HTMLDocs/{uri:.*}").HandlerFunc(dataVisHandler)
	router.Host("{subdomain:[a-z]+}.neighbourhood.statistics.gov.uk").Path("/HTMLDocs/{uri:.*}").HandlerFunc(dataVisHandler)
	// NeSS API
	router.Host("neighbourhood.statistics.gov.uk").Path("/NDE2/{uri:.*}").HandlerFunc(apiHandler)
	router.Host("{subdomain:[a-z]+}.neighbourhood.statistics.gov.uk").Path("/NDE2/{uri:.*}").HandlerFunc(apiHandler)
	// WDA website
	router.Host("web.ons.gov.uk").Path("/ons/apiservice/web/{uri:.*}").HandlerFunc(defaultHandler)
	// WDA API
	router.Host("web.ons.gov.uk").Path("/ons/apiservice/{uri:.*}").HandlerFunc(apiHandler)
	router.Host("web.ons.gov.uk").Path("/ons/api/{uri:.*}").HandlerFunc(apiHandler)
	router.Host("data.ons.gov.uk").Path("/{uri:.*}").HandlerFunc(apiHandler)
	// Visual.ONS
	router.Host("visual.ons.gov.uk").Path("/wp-content/uploads/{uri:.*}").HandlerFunc(visualAssetHandler)
	router.Host("visual.ons.gov.uk").Path("/{article:[^/]*}{uri:/?.*}").HandlerFunc(visualArticleHandler)
	// Catch-all
	router.Path("/{uri:.*}").HandlerFunc(defaultHandler)

	return router
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	log.Event(req.Context(), "redirecting to landing page", log.INFO, log.Data{
		"host": req.Host,
		"path": req.URL.Path,
		"dest": landingPage,
	})
	w.Header().Set("Location", landingPage)
	w.WriteHeader(redir)
}

func dataVisHandler(w http.ResponseWriter, req *http.Request) {
	dest := "https://www.ons.gov.uk/visualisations/nesscontent/" + mux.Vars(req)["uri"]
	log.Event(req.Context(), "redirecting visualisation", log.INFO, log.Data{
		"host": req.Host,
		"path": req.URL.Path,
		"dest": dest,
	})
	w.Header().Set("Location", dest)
	w.WriteHeader(redir)
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	log.Event(req.Context(), "returning api help text", log.INFO, log.Data{
		"host": req.Host,
		"path": req.URL.Path,
	})
	w.WriteHeader(410)
	w.Write([]byte(apiResponse))
}

func visualAssetHandler(w http.ResponseWriter, req *http.Request) {
	dest := "https://static.ons.gov.uk/visual/" + mux.Vars(req)["uri"]
	log.Event(req.Context(), "redirecting visual.ons.gov.uk wp-content", log.INFO, log.Data{
		"host": req.Host,
		"path": req.URL.Path,
		"dest": dest,
	})
	w.Header().Set("Location", dest)
	w.WriteHeader(redir)
}

func visualArticleHandler(w http.ResponseWriter, req *http.Request) {
	article := mux.Vars(req)["article"]
	uri := mux.Vars(req)["uri"]

	if len(article) == 0 {
		log.Event(req.Context(), "redirecting visual request to ONS", log.INFO, log.Data{
			"article": article,
			"uri":     uri,
			"host":    req.Host,
			"path":    req.URL.Path,
		})

		w.Header().Set("Location", "https://www.ons.gov.uk")
		w.WriteHeader(redir)
		return
	}

	if dest, ok := visualRedirects[article]; ok {
		log.Event(req.Context(), "redirecting visual request to ONS", log.INFO, log.Data{
			"article": article,
			"uri":     uri,
			"host":    req.Host,
			"path":    req.URL.Path,
		})

		w.Header().Set("Location", dest)
		w.WriteHeader(redir)
		return
	}

	log.Event(req.Context(), "redirecting visual request to national archives", log.INFO, log.Data{
		"article": article,
		"uri":     uri,
		"host":    req.Host,
		"path":    req.URL.Path,
	})
	w.Header().Set("Location", "http://webarchive.nationalarchives.gov.uk/20171102124620/https://visual.ons.gov.uk/"+article+uri)
	w.WriteHeader(redir)
}
