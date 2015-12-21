# List building
ALL_LIST = agent.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test -race
GOFMT=gofmt -w

PREBUILD_LIST = $(foreach int, $(ALL_LIST), $(int)_prebuild)
POSTBUILD_LIST = $(foreach int, $(ALL_LIST), $(int)_postbuild)
BUILD_LIST_OSX = $(foreach int, $(ALL_LIST), $(int)_build_osx)
BUILD_LIST_WIN = $(foreach int, $(ALL_LIST), $(int)_build_win)
BUILD_LIST_LINUX = $(foreach int, $(ALL_LIST), $(int)_build_linux)
BUILD_LIST_LINUX_386 = $(foreach int, $(ALL_LIST), $(int)_build_linux_386)
TEST_LIST = $(foreach int, $(ALL_LIST), $(int)_test)
FMT_TEST = $(foreach int, $(ALL_LIST), $(int)_fmt)
RUN_LIST = $(foreach int, $(ALL_LIST), $(int)_run)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(BUILD_LIST)

all: build
build: prebuild $(BUILD_LIST_OSX) $(BUILD_LIST_WIN) $(BUILD_LIST_LINUX) postbuild
build_osx: prebuild $(BUILD_LIST_OSX) postbuild
build_win: prebuild $(BUILD_LIST_WIN) postbuild
build_linux: prebuild $(BUILD_LIST_LINUX) postbuild
build_linux_386: prebuild $(BUILD_LIST_LINUX_386) postbuild
clean: $(CLEAN_LIST)
test: $(TEST_LIST)
fmt: $(FMT_TEST)
run: $(RUN_LIST)

prebuild:
	@if [ -f ./prebuild ]; then \
		echo "Running prebuild script in release mode..." ; \
		./prebuild --release ; \
	else  \
		echo "No pre-build script found in pre-build phase; skipping." ; \
	fi

postbuild:
	@if [ -f ./prebuild ]; then \
		echo "Running prebuild script in debug mode..." ; \
		./prebuild --debug ; \
	else  \
		echo "No pre-build script found in post-build phase; skipping." ; \
	fi

$(BUILD_LIST_OSX): %_build_osx: %_fmt
	@echo "Building Darwin AMD64..."
	@GO15VENDOREXPERIMENT=1 GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) -tags release -o bin/darwin-amd64/telemetry_agent
	@echo "Building complete."

$(BUILD_LIST_WIN): %_build_win: %_fmt
	@echo "Building Windows 386..."
	@GO15VENDOREXPERIMENT=1 GOARCH=386 CGO_ENABLED=1 GOOS=windows CC="i686-w64-mingw32-gcc -fno-stack-protector -D_FORTIFY_SOURCE=0 -lssp -D_localtime32=localtime" $(GOBUILD)  -tags release -o bin/windows-386/telemetry_agent.exe
	@echo "Building complete."

$(BUILD_LIST_LINUX): %_build_linux: %_fmt
	@echo "Building Linux AMD64..."
	@GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC="gcc" $(GOBUILD) --ldflags '-extldflags "-static"' -tags release -o bin/linux-amd64/usr/bin/telemetry_agent
	@echo "Building complete."
	@echo Building DEB and RPM files
	@rm -Rf /tmp/telemetry_agent
	@mkdir -p /tmp/telemetry_agent/usr/bin
	@cp bin/linux-amd64/usr/bin/telemetry_agent /tmp/telemetry_agent/usr/bin
	@chmod 755 /tmp/telemetry_agent/usr/bin/telemetry_agent
	@cp VERSION /tmp/TELEMETRY_AGENT_VERSION
	@cd /tmp/telemetry_agent && fpm -s dir -t deb -n "telemetry_agent" -v `cat ../TELEMETRY_AGENT_VERSION` usr
	@cd /tmp/telemetry_agent/ && fpm -s dir -t rpm -n "telemetry_agent" -v `cat ../TELEMETRY_AGENT_VERSION` --rpm-sign usr
	@cp /tmp/telemetry_agent/*rpm bin/linux-amd64
	@cp /tmp/telemetry_agent/*deb bin/linux-amd64

$(BUILD_LIST_LINUX_386): %_build_linux_386: %_fmt
	@echo "Building Linux 386..."
	@GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=386 CGO_ENABLED=1 CC="gcc" $(GOBUILD) --ldflags '-extldflags "-static"' -tags release -o bin/linux-386/usr/bin/telemetry_agent
	@echo "Building complete."
	@echo Building DEB and RPM files
	@rm -Rf /tmp/telemetry_agent
	@mkdir -p /tmp/telemetry_agent/usr/bin
	@cp bin/linux-386/usr/bin/telemetry_agent /tmp/telemetry_agent/usr/bin
	@chmod 755 /tmp/telemetry_agent/usr/bin/telemetry_agent
	@cp VERSION /tmp/TELEMETRY_AGENT_VERSION
	@cd /tmp/telemetry_agent && fpm -s dir -t deb -a i386 -n "telemetry_agent" -v `cat ../TELEMETRY_AGENT_VERSION` usr
	@cd /tmp/telemetry_agent/ && fpm -s dir -t rpm -a i386 -n "telemetry_agent" -v `cat ../TELEMETRY_AGENT_VERSION` --rpm-sign usr
	@cp /tmp/telemetry_agent/*rpm bin/linux-386
	@cp /tmp/telemetry_agent/*deb bin/linux-386



$(TEST_LIST): %_test:
	@echo "Running go test..."
	@$(GOTEST) ./...

$(FMT_TEST): %_fmt:
	@echo "Running go fmt..."
	@$(GOFMT) agent.go agent plugin
