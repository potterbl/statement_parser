VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.0.0)

NEXT_PATCH := $(shell echo $(VERSION) | awk -F. '{printf "v%d.%d.%d", $$1, $$2, $$3+1}')

changelog:
	@echo "## $(NEXT_PATCH)" > changelog.md
	@git log $(VERSION)..HEAD --pretty=format:"* %s" --no-merges >> changelog.md
	@echo "" >> changelog.md

release:
	@echo "Current version: $(VERSION)"
	@echo "Releasing version: $(NEXT_PATCH)"
	git tag $(NEXT_PATCH)
	git push origin $(NEXT_PATCH)
