# Telemetry Agent

The Telemetry Agent simplifies the process of creating daemon processes that feed data into one or more [Telemetry](http://telemetryapp.com) flows.

Typical use-case scenarios include:

  - Feeding data from existing infrastructure (e.g.: a MySQL database, Excel sheet, custom script written in your language of choice) to one or more Telemetry data flows
  - Automatically creating boards for your customers
  - Interfacing third-party APIs with Telemetry

The Agent is written in Go and has been built to run on most Linux distros, OS X, and Windows. We offer a list of downloadable binaries on our [releases page](https://github.com/telemetryapp/gotelemetry_agent/releases). The Agent is designed to run on your infrastructure, and its only requirement is that it be able to reach the Telemetry API endpoint (https://api.telemetryapp.com) on port 443 via HTTPS. It can therefore happily live behind firewalls without posing a security risk.

Full documentation is available on the [Telemetry Documentation website](https://telemetry.readme.io/docs/telemetry-agent).

## Building

You will need a working install of Go 1.5 and GIT on your local platform in order to build the Agent from source.

The Agent is designed to cross compile using [goxc](https://github.com/laher/goxc/wiki/config) and includes a `.goxc.json` config file to produce a validated build for all compatible platorms.

You can also compile for your current platform by using the `go build` command. Ensure that your packages are up to date by running a `go get -u` prior to building.
