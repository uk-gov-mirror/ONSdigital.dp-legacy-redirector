VERSION=1.1.0

build:
	GOOS=linux GOARCH=amd64 go build -o build/dp-legacy-redirector .

docker: build
	curl -o ca-certificates.crt https://raw.githubusercontent.com/bagder/ca-bundle/master/ca-bundle.crt
	docker build -t onsdigital/dp-legacy-redirector:$(VERSION) .

release: build
	zip dp-legacy-redirector-$(VERSION).zip build/dp-legacy-redirector Dockerfile ca-certificates.crt

.PHONY: build docker
