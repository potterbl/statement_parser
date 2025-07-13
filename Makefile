VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.0.0)
NEXT_PATCH := $(shell echo $(VERSION) | awk -F. '{printf "v%d.%d.%d", $$1, $$2, $$3+1}')

changelog:
	@echo "## $(NEXT_PATCH)" > .changelog.tmp
	@git log $(VERSION)..HEAD --pretty=format:"* %s" --no-merges >> .changelog.tmp
	@echo "" >> .changelog.tmp
	@cat changelog.md >> .changelog.tmp 2>/dev/null || true
	@mv .changelog.tmp changelog.md

release: changelog
	@echo "Current version: $(VERSION)"
	@echo "Releasing version: $(NEXT_PATCH)"
	git add changelog.md
	git commit -m "Update changelog for $(NEXT_PATCH)"
	git tag $(NEXT_PATCH)
	git push origin $(NEXT_PATCH)
	git push origin HEAD