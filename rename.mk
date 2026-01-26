OLD_PROJECT := heppu/go-template
OLD_MODULE  := github.com/${OLD_PROJECT}
OLD_NAME    := demo
OLD_CMD_DIR := cmd/demo
GH_INSTANCE ?= github.com
NEW_URL     ?= $(shell git config --get remote.origin.url)
NEW_MODULE  ?= $(subst ssh://,,$(subst .git,,$(subst https://,,${NEW_URL})))
# Remove protocol prefixes one by one.
NEW_MODULE := $(patsubst https://${GH_INSTANCE}/%,${GH_INSTANCE}/%,${NEW_MODULE})
NEW_MODULE := $(patsubst ssh://git@${GH_INSTANCE}:%,${GH_INSTANCE}/%,${NEW_MODULE})
NEW_MODULE := $(patsubst git@${GH_INSTANCE}:%,${GH_INSTANCE}/%,${NEW_MODULE})
# Remove the trailing ".git"
NEW_MODULE  := $(patsubst %.git,%, ${NEW_MODULE})
NEW_PROJECT := ${NEW_MODULE:${GH_INSTANCE}/%=%}
NEW_NAME    ?= $(notdir ${NEW_MODULE})
NEW_CMD_DIR	:= cmd/${NEW_NAME}
# Sed go alternative
GOSED := go tool github.com/rwtodd/Go.Sed/cmd/sed-go

.PHONY: rename
rename:
	@printf "Renaming project using following configuration:\n\n"
	@printf "URL:     ${NEW_URL}\n"
	@printf "MODULE:  ${NEW_MODULE}\n"
	@printf "PROJECT: ${NEW_PROJECT}\n"
	@printf "NAME:    ${NEW_NAME}\n\n"
	@go mod edit -module ${NEW_MODULE}
	@find . -type f -name '*.go' -exec ${GOSED} -i 's|${OLD_MODULE}|${NEW_MODULE}|g' {} \;
	@mv ${OLD_CMD_DIR} ${NEW_CMD_DIR}
	@${GOSED} -i 's|${OLD_CMD_DIR}|${NEW_CMD_DIR}|g' ./applicationtest/application_test.go
	@${GOSED} -i 's|${OLD_PROJECT}|${NEW_PROJECT}|g' ./README.md
	@${GOSED} -i 's|${OLD_NAME}|${NEW_NAME}|g' ./Makefile
	@${GOSED} -i 's|${OLD_NAME}|${NEW_NAME}|g' ./Dockerfile
	@gofmt -w .
	@printf "\nProject renamed successfully, deleting rename.mk\n"
	@rm rename.mk
