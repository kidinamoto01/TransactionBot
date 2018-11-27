PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation' | grep -v '/prometheus' | grep -v '/clitest' | grep -v '/lcd' | grep -v '/protobuf')
PACKAGES_MODULES=$(shell go list ./... | grep 'modules')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')

all: get_tools get_vendor_deps install

COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_FLAGS = -ldflags "-X github.com/irisnet/irishub/version.GitCommit=${COMMIT_HASH}"

########################################
### Tools & dependencies

check_tools:
	cd deps_tools && $(MAKE) check_tools

check_dev_tools:
	cd deps_tools && $(MAKE) check_dev_tools

update_tools:
	cd deps_tools && $(MAKE) update_tools

update_dev_tools:
	cd deps_tools && $(MAKE) update_dev_tools

get_tools:
	cd deps_tools && $(MAKE) get_tools

get_dev_tools:
	cd deps_tools && $(MAKE) get_dev_tools

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v

draw_deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/irisnet/irishub/cmd/iris -d 2 | dot -Tpng -o dependency-graph.png


########################################
### Compile and Install
install:
	go install ./irisrobot