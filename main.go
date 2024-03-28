package main

import (
	"context"
	"net/http"
	"os"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-legacy-redirector/config"
	server "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

var redir = http.StatusTemporaryRedirect
var landingPage = "https://www.ons.gov.uk/help/localstatistics"
var apiResponse = "This service is no longer available. Please visit https://www.ons.gov.uk/help/localstatistics for more information."

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
		log.Fatal(ctx, "unable to retrieve service configuration", err)
		os.Exit(1)
	}

	log.Info(ctx, "config on startup", log.Data{"config": cfg, "build_time": BuildTime, "git-commit": GitCommit})

	// Health check
	versionInfo, err := healthcheck.NewVersionInfo(BuildTime, GitCommit, Version)
	if err != nil {
		log.Fatal(ctx, "Failed to obtain VersionInfo for healthcheck", err)
		os.Exit(1)
	}
	hc := healthcheck.New(versionInfo, cfg.HealthckeckCriticalTimeout, cfg.HealthckeckInterval)

	srv := server.NewServer(cfg.BindAddr, getRouter(hc))

	log.Info(ctx, "starting http server", log.Data{"bind_addr": cfg.BindAddr})
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(ctx, "error starting server", err)
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
	log.Info(req.Context(), "redirecting to landing page", log.Data{
		"host": req.Host,
		"path": req.URL.Path,
		"dest": landingPage,
	})
	w.Header().Set("Location", landingPage)
	w.WriteHeader(redir)
}

func dataVisHandler(w http.ResponseWriter, req *http.Request) {
	dest := "https://www.ons.gov.uk/visualisations/nesscontent/" + mux.Vars(req)["uri"]
	log.Info(req.Context(), "redirecting visualisation", log.Data{
		"host": req.Host,
		"path": req.URL.Path,
		"dest": dest,
	})
	w.Header().Set("Location", dest)
	w.WriteHeader(redir)
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	log.Info(req.Context(), "returning api help text", log.Data{
		"host": req.Host,
		"path": req.URL.Path,
	})
	w.WriteHeader(410)
	_, err := w.Write([]byte(apiResponse))
	if err != nil {
		log.Error(req.Context(), "error writing response", err)
	}
}

func visualAssetHandler(w http.ResponseWriter, req *http.Request) {
	dest := "https://static.ons.gov.uk/visual/" + mux.Vars(req)["uri"]
	log.Info(req.Context(), "redirecting visual.ons.gov.uk wp-content", log.Data{
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
		log.Info(req.Context(), "redirecting visual request to ONS", log.Data{
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
		log.Info(req.Context(), "redirecting visual request to ONS", log.Data{
			"article": article,
			"uri":     uri,
			"host":    req.Host,
			"path":    req.URL.Path,
		})

		w.Header().Set("Location", dest)
		w.WriteHeader(redir)
		return
	}

	log.Info(req.Context(), "redirecting visual request to national archives", log.Data{
		"article": article,
		"uri":     uri,
		"host":    req.Host,
		"path":    req.URL.Path,
	})
	w.Header().Set("Location", "http://webarchive.nationalarchives.gov.uk/20171102124620/https://visual.ons.gov.uk/"+article+uri)
	w.WriteHeader(redir)
}
