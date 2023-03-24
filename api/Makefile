.PHONY: protolint
protolint:
	@(protolint proto/)

.PHONY: gen-docs
gen-docs:
	@(./gen-docs.sh)

.PHONY: gen-java
gen-java:
	@(./gen-java.sh)

.PHONY: gen-go
gen-go:
	@(./gen-go.sh)

.PHONY: go-vendor
go-vendor:
	go mod tidy && go mod vendor

.PHONY: gen-all
gen-all: gen-docs gen-go go-vendor