# Makefile
# compile yang file

YANG_REPO := ./yang-repo
OUTPUT_DIR := ./internal/oc

YANG_PATHS := $(shell find $(YANG_REPO) -type d ! -path '*/.*' | paste -sd, -)

.PHONY: generate-yang
generate-yang:
	go run github.com/openconfig/ygnmi/app/ygnmi generator \
    	--base_package_path=github.com/LeBaoTai/SDN-RD/internal/oc \
    	--output_dir=$(OUTPUT_DIR) \
    	--paths=$(YANG_PATHS) \
    	--compress_paths=true \
		--exclude_modules=ietf-interfaces \
    	--fakeroot_name=root \
    	--split_top_level_packages=true \
		--pathstructs_split_files_count=10 \
		$(YANG_REPO)/openconfig/interfaces/openconfig-interfaces.yang \
		$(YANG_REPO)/openconfig/interfaces/openconfig-if-ip.yang \
		$(YANG_REPO)/openconfig/bgp/openconfig-bgp.yang \
		$(YANG_REPO)/openconfig/acl/openconfig-acl.yang \
		$(YANG_REPO)/openconfig/network-instance/openconfig-network-instance.yang \
		$(YANG_REPO)/openconfig/system/openconfig-system.yang


.PHONY: verify-yang
verify-yang: generate-yang
	go build ./internal/oc/...
	@echo "ocgen package builds successfully"


.PHONY: run-backend
run-backend:
	go run ./cmd/backend

.PHONY: run-controller
run-controller:
	go run ./cmd/controller