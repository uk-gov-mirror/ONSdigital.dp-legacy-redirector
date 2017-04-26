VERSION=0.1.0

build:
	GOOS=linux GOARCH=amd64 go build -o build/dp-ness-wda-redirector .

docker: build
	curl -o ca-certificates.crt https://raw.githubusercontent.com/bagder/ca-bundle/master/ca-bundle.crt
	docker build -t onsdigital/dp-ness-wda-redirector:$(VERSION) .

.PHONY: build docker
