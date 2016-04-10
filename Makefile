default: help

## Make bin for smsmail.
bin:
	./control build

## Get vet go tools.
vet:
	go get golang.org/x/tools/cmd/vet

## Validate this go project.
validate: vet
	bash script/validate-gofmt
	go list ./... | grep -v 'vendor' | xargs -L1 go vet

## Run golint for this go project.
#lint:
#	go list ./... | grep -v /vendor/ | xargs -L1 fgt golint

## Run test case for this go project.
test:
	go list ./... | grep -v 'vendor' | xargs -L1 go test -v

## Clean everything (including stray volumes).
clean:
	-rm -rf var
	-rm -f smsmail

help: # Some kind of magic from https://gist.github.com/rcmachado/af3db315e31383502660
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9]+:/ {                                   \
		nb = sub( /^## /, "", helpMsg );                             \
		if(nb == 0) {                                                \
			helpMsg = $$0;                                             \
			nb = sub( /^[^:]*:.* ## /, "", helpMsg );                  \
		}                                                            \
		if (nb)                                                      \
			printf "\033[1;31m%-" width "s\033[0m %s\n", $$1, helpMsg; \
	}                                                              \
	{ helpMsg = $$0 }'                                             \
	width=$$(grep -o '^[a-zA-Z_0-9]\+:' $(MAKEFILE_LIST) | wc -L)  \
	$(MAKEFILE_LIST)
