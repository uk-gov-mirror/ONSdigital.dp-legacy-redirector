dp-legacy-redirector
====================

Handles redirects for legacy ONS services including Neighbourhood Statistics, Web Data Access
and Visual.ONS.

### Getting started

To run locally:

* Clone this repo
* `go run main.go`

To build for release:

* `make docker`

### Configuration

Configuration for the redirector.

| Environment variable | Default | Description
| -------------------- | ------- | -----------
| BIND_ADDR            | :8080   | The host and port to bind to

### License

Copyright ©‎ 2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
