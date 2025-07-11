VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.0.0)

# Расчёт следующей патч-версии (vX.Y.Z => vX.Y.(Z+1))
NEXT_PATCH := $(shell echo $(VERSION) | awk -F. '{printf "v%d.%d.%d", $$1, $$2, $$3+1}')

release:
	@echo "Current version: $(VERSION)"
	@echo "Releasing version: $(NEXT_PATCH)"
	git tag $(NEXT_PATCH)
	git push origin $(NEXT_PATCH)
