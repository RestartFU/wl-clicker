GO ?= go
APP ?= clicker
WINDOWS_APP ?= $(APP).exe
CMD ?= ./cmd/clicker
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
DESTDIR ?=
INSTALL ?= install
WINDOWS_GOARCH ?= amd64
WINDOWS_CC ?= x86_64-w64-mingw32-gcc
WINDOWS_CXX ?= x86_64-w64-mingw32-g++

.PHONY: help build build-windows run test fmt vet clean rebuild install

help:
	@echo "Targets:"
	@echo "  build          Build $(APP) for the current platform"
	@echo "  build-windows  Build $(WINDOWS_APP) for Windows"
	@echo "  install  Install $(APP) to $(DESTDIR)$(BINDIR)"
	@echo "  run      Run the app from source"
	@echo "  test     Run tests"
	@echo "  fmt      Format Go code"
	@echo "  vet      Run go vet"
	@echo "  clean    Remove build artifacts"
	@echo "  rebuild  Clean then build"

build:
	$(GO) build -o $(APP) $(CMD)

build-windows:
	@if command -v $(WINDOWS_CC) >/dev/null 2>&1; then \
		CGO_ENABLED=1 GOOS=windows GOARCH=$(WINDOWS_GOARCH) CC=$(WINDOWS_CC) CXX=$(WINDOWS_CXX) \
			$(GO) build -ldflags="-H windowsgui" -o $(WINDOWS_APP) $(CMD); \
	else \
		echo "error: $(WINDOWS_CC) not found"; \
		echo "Fyne desktop Windows builds need CGO_ENABLED=1 plus a MinGW cross-compiler."; \
		echo "Install mingw-w64 (or override WINDOWS_CC/WINDOWS_CXX), then rerun make build-windows."; \
		exit 1; \
	fi

install: build
	mkdir -p $(DESTDIR)$(BINDIR)
	$(INSTALL) -m 0755 $(APP) $(DESTDIR)$(BINDIR)/$(APP)
	    rm -rf $(APP)

run:
	$(GO) run $(CMD)

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

clean:
	rm -f $(APP) $(WINDOWS_APP)

rebuild: clean build
