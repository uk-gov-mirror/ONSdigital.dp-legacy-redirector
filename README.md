# dp-legacy-redirector

Handles redirects for legacy ONS services including Neighbourhood Statistics, Web Data Access
and Visual.ONS.

## Getting started

To run locally:

* Clone this repo
* `go run main.go`

To build for release:

* `make docker`

## Configuration

Configuration for the redirector.

| Environment variable         | Default | Description                                   |
| ---------------------------- | ------- | --------------------------------------------- |
| BIND_ADDR                    | :24600  | The host and port to bind to                  |
| HEALTHCHECK_INTERVAL         | 60s     | The period of time between health checks      |
| HEALTHCHECK_CRITICAL_TIMEOUT | 5s      | The period of time after which failing checks |

## License

Copyright ©‎ 2017-2025, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
