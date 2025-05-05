OLD_MODULE  := github.com/heppu/go-template
OLD_NAME	:= demo
OLD_CMD_DIR	:= cmd/demo
NEW_URL     ?= $(shell git config --get remote.origin.url)
NEW_MODULE  ?= $(subst ssh://,,$(subst .git,,$(subst https://,,${NEW_URL})))
# Remove protocol prefixes one by one.
NEW_MODULE := $(patsubst https://github.com/%,github.com/%,${NEW_MODULE})
NEW_MODULE := $(patsubst ssh://git@github.com:%,github.com/%,${NEW_MODULE})
NEW_MODULE := $(patsubst git@github.com:%,github.com/%,${NEW_MODULE})
# Remove the trailing ".git"
NEW_MODULE := $(patsubst %.git,%, ${NEW_MODULE})
NEW_NAME    ?= $(notdir ${NEW_MODULE})
NEW_CMD_DIR	:= cmd/${NEW_NAME}

ifeq ($(findstring GNU,$(shell strings $$(which sed))),)
    SED_INPLACE_ARG := ''
endif

.PHONY: rename
rename:
	@printf "Renaming project using following configuration:\n\n"
	@printf "URL:    ${NEW_URL}\n"
	@printf "MODULE: ${NEW_MODULE}\n"
	@printf "NAME:   ${NEW_NAME}\n\n"
	@go mod edit -module ${NEW_MODULE}
	@find . -type f -name '*.go' -exec sed -i ${SED_INPLACE_ARG} 's|${OLD_MODULE}|${NEW_MODULE}|g' {} \;
	@mv ${OLD_CMD_DIR} ${NEW_CMD_DIR}
	@sed -i ${SED_INPLACE_ARG} 's|${OLD_CMD_DIR}|${NEW_CMD_DIR}|g' ./applicationtest/application_test.go
	@sed -i ${SED_INPLACE_ARG} 's|${OLD_MODULE}|${NEW_MODULE}|g' ./README.md
	@sed -i ${SED_INPLACE_ARG} 's|${OLD_NAME}|${NEW_NAME}|g' ./Makefile
	@sed -i ${SED_INPLACE_ARG} 's|${OLD_NAME}|${NEW_NAME}|g' ./Dockerfile
	@printf "\nProject renamed succesfully, deleting rename.mk\n"
	@rm rename.mk
