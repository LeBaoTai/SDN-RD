# Makefile
# compile yang file

YANG_REPO := ../yang-repo
OUTPUT    := ./internal/oc/oc.go

.PHONY: generate-yang
generate-yang:
	go run github.com/openconfig/ygot/generator \
		-path=$(YANG_REPO) \
		-output_file=./internal/oc/oc.go \
		-package_name=oc \
		-generate_fakeroot \
		-fakeroot_name=Device \
		-shorten_enum_leaf_names \
		-typedef_enum_with_defmod \
		-exclude_modules=ietf-interfaces \
		-compress_paths=true \
		-generate_simple_unions \
		-generate_populate_defaults=true \
		-generate_getters \
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
