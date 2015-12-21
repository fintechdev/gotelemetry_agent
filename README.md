# Telemetry Agent

The Telemetry Agent simplifies the process of creating daemon processes that feed data into one or more [Telemetry](http://telemetryapp.com) flows.

Typical use-case scenarios include:

  - Feeding data from existing infrastructure (e.g.: a MySQL database, Excel sheet, custom script written in your language of choice) to one or more Telemetry data flows
  - Automatically creating boards for your customers
  - Interfacing third-party APIs with Telemetry

The Agent is written in Go and runs fine on most Linux distros, OS X, and Windows. It is designed to run on your infrastructure, and its only requirement is that it be able to reach the Telemetry API endpoint (https://api.telemetryapp.com) on port 443 via HTTPS. It can therefore happily live behind firewalls without posing a security risk.

Full documentation is available on the [Telemetry Documentation website](http://telemetry.readme.io/v1.0/docs/agents).

## Building

In order to build the Agent, you will need a working install of Go 1.5 for each of the platforms you wish to target. You will also require a C/C++ build environment (GCC on Linux and BSD systems, and MinGW on Windows), since many of the Agent's dependencies require building external libraries.

The Agent uses [gom](https://github.com/mattn/gom) for building a directory of vendored dependencies that are pinned to specific commits to ensure that builds are reproducible. Gom is also used to build the binaries.

If you're building for Linux, you will also need to have a working copy of Jordan Sissel's [fpm](https://github.com/jordansissel/fpm) if you want to generate DEB or RPM files.

### Using the Makefile

The project's makefile supports the following commands:

- `make build_osx` builds the Agent for OS X. The resulting binary is statically linked and placed in `bin/darwin-amd64/telemetry_agent`.
- `make build_win` builds a 32-bit version fo the Agent for Windows. The resulting binary goes to `bin/windows-386/telemetry_agent.exe`.
- `make build_linux` builds for AMD64 versions of Linux. The resulting binary is placed in `bin/linux-amd64/usr/bin/telemetry_agent`, and RPM and DEB files are automatically built and signed.
- `make build_linux_386` builds an i386 binary for Linux, and is otherwise equal to `build_linux`

Please note that you must run `gom install` before starting make if you wish to build against pinned dependencies.

### Building on Linux

The Linux build commands automatically uses `fpm` to generate both DEB and RPM packages for the Agent. RPMs, in particular, are automatically signed using your current designated GPG key.

Note that, in keeping with Go's philosophy of using a single executable for easy deployment, Linux builds are completely static and may require static versions of several libraries (e.g.: `yum install glibc-static` on RedHat systems).