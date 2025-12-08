MODULE := github.com/conv-project/go-shared

# Версию передаём как аргумент: make release VERSION=v1.2.3
VERSION ?=

.PHONY: tag release list-tags

# Локально поставить аннотированный тег без push
tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: VERSION is empty. Use: make tag VERSION=v1.2.3"; \
		exit 1; \
	fi
	@git diff --quiet || (echo "ERROR: Working tree is dirty. Commit changes first."; exit 1)
	@git tag -a "$(VERSION)" -m "Release $(VERSION)"
	@echo "Tag $(VERSION) created."

# Полный релиз: тег + push
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: VERSION is empty. Use: make release VERSION=v1.2.3"; \
		exit 1; \
	fi
	@git diff --quiet || (echo "ERROR: Working tree is dirty. Commit changes first."; exit 1)
	@git tag -a "$(VERSION)" -m "Release $(VERSION)"
	@git push origin "$(VERSION)"
	@echo "Released $(MODULE) version $(VERSION)."

# Посмотреть теги по версиям
list-tags:
	@git tag --sort=version:refname