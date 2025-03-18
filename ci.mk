include Makefile

.PHONY: .push
push:
	docker push --all-tags ${IMG_NAME}

.PHONY: pr-check
pr-check: main-check check-modified-migrations

.PHONY: main-check
main-check: tidy-check check-duplicated-migrations check-generated-code

.PHONY: check-duplicated-migrations
check-duplicated-migrations:
	@set -e; \
	echo "Performing duplicated migration check"; \
	output="$$(ls -1 store/migrations/ | cut -d "_" -f1 | uniq -D)"; \
	if [ -n "$$output" ]; then \
		echo "Found duplicate migration versions:"; \
		echo "$$output"; \
		exit 1; \
	fi; \
	echo "No duplicated migrations found"

.PHONY: check-modified-migrations
check-modified-migrations:
	@set -e; \
	if test -z "$$BASE_REF"; then \
		echo "BASE_REF must be set"; \
		exit 1; \
	fi; \
	if test -z "$$HEAD_REF"; then \
		echo "HEAD_REF must be set"; \
		exit 1; \
	fi; \
	echo "Performing migration verification on PR against $$BASE_REF"; \
	git fetch origin $$BASE_REF; \
	git fetch origin $$HEAD_REF; \
	git diff --exit-code --name-only --diff-filter=a origin/$$BASE_REF origin/$$HEAD_REF -- store/migrations/ || (echo "Migrations files are out of sync, please rebase" && exit 1); \
	echo "No modified migrations found"

.PHONY: check-generated-code
check-generated-code: generate
	@git diff --exit-code --name-only || (printf "Generated code isn't up to date.\nRun make generate before commiting" && exit 1)

.PHONY: tidy-check
tidy-check:
	go mod tidy
	git diff --exit-code --name-status -- go.mod go.sum
