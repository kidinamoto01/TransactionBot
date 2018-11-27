
all: get_tools get_vendor_deps install

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

########################################
### Compile and Install
install:
	go install ./irisrobot

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/irisrobot ./irisrobot
