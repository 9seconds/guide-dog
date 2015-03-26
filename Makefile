# Cross platorm build with docker

# ----------------------------------------------------------------------------

ROOT_DIR        := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_PROG      := $(ROOT_DIR)/guide-dog
CROSS_BUILD_DIR := $(ROOT_DIR)/build
PROJECT         := github.com/9seconds/guide-dog

LINUX_ARCH      := amd64 386 arm
DARWIN_ARCH     := amd64 386
FREEBSD_ARCH    := amd64 386 arm

DOCKER_PROG     := docker
DOCKER_GOPATH   := /go
DOCKER_WORKDIR  := $(DOCKER_GOPATH)/src/$(PROJECT)
DOCKER_IMAGE    := golang:1.4.2-cross

PROC_COUNT      := $(shell cat /proc/cpuinfo | awk '/processor/ {n++}; END {print n}')
UID             := $(shell id -u)
GID             := $(shell id -g)

# ----------------------------------------------------------------------------

define crosscompile
	GOOS=$(1) GOARCH=$(2) go build -a -o $(CROSS_BUILD_DIR)/$(1)-$(2) $(PROJECT) && \
	chown $(DEV_UID):$(DEV_GID) $(CROSS_BUILD_DIR)/$(1)-$(2)
endef

# ----------------------------------------------------------------------------

all: tools prog-build
tools: fix lint
cross: cross-linux cross-darwin cross-freebsd
clean: prog-clean cross-clean repo-clean
ci: tools cross

# ----------------------------------------------------------------------------

fix:
	@go fix $(PROJECT)/...

lint:
	@golint $(PROJECT)/...

godep:
	@go get -u github.com/kr/godep

fmt:
	@go fmt $(PROJECT)/...

save: godep
	@godep save

restore: godep
	@godep restore

prog-build: restore prog-clean
	@go build -a -o $(BUILD_PROG) $(PROJECT)

install: restore
	@go install -a $(PROJECT)

prog-clean:
	@rm -f $(BUILD_PROG)

update:
	@grep -v $(PROJECT) $(ROOT_DIR)/Godeps/Godeps.json \
		| awk '/ImportPath/ {gsub(/"|,/, ""); print $$2}' \
		| xargs -n 1 godep update

upgrade_deps:
	@grep -v $(PROJECT) $(ROOT_DIR)/Godeps/Godeps.json \
		| awk '/ImportPath/ {gsub(/"|,/, ""); print $$2}' \
		| xargs -n 1 -P 4 go get -u

# ----------------------------------------------------------------------------

cross-linux: $(addprefix cross-linux-,$(LINUX_ARCH))
cross-freebsd: $(addprefix cross-freebsd-,$(FREEBSD_ARCH))
cross-darwin: $(addprefix cross-darwin-,$(DARWIN_ARCH))

cross-clean:
	@rm -rf $(CROSS_BUILD_DIR)

cross-build-directory: cross-clean
	@mkdir -p $(CROSS_BUILD_DIR) && chown -R $(UID):$(GID) $(CROSS_BUILD_DIR)

cross-linux-%: restore cross-build-directory
	$(call crosscompile,linux,$*)

cross-darwin-%: restore cross-build-directory
	$(call crosscompile,darwin,$*)

cross-freebsd-%: restore cross-build-directory
	$(call crosscompile,freebsd,$*)

cross-docker:
	@$(DOCKER_PROG) run \
		--rm=true \
		-e DEV_UID=$(UID) \
		-e DEV_GID=$(GID) \
		-i -t \
		-v "$(ROOT_DIR)":$(DOCKER_WORKDIR) \
		-w $(DOCKER_WORKDIR) \
		$(DOCKER_IMAGE) \
	make -j $(PROC_COUNT) cross

# ----------------------------------------------------------------------------

repo-clean:
	@git clean -xfd
