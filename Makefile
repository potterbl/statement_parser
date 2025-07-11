VERSION ?= $(shell git describe --tags --abbrev=0)
NEXT_PATCH := $(shell echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}')

release:
	@echo "Current version: $(VERSION)"
	@echo "Releasing version: v$(NEXT_PATCH)"
	git tag v$(NEXT_PATCH)
	git push origin v$(NEXT_PATCH)
