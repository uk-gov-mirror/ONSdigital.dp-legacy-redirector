VERSION=1.0.0

build:
	GOOS=linux GOARCH=amd64 go build -o build/dp-ness-wda-redirector .

docker: build
	curl -o ca-certificates.crt https://raw.githubusercontent.com/bagder/ca-bundle/master/ca-bundle.crt
	docker build -t onsdigital/dp-ness-wda-redirector:$(VERSION) .

release: build
	zip dp-ness-wda-redirector-$(VERSION).zip build/dp-ness-wda-redirector Dockerfile ca-certificates.crt

.PHONY: build docker
