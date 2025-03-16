OLD_MODULE  := github.com/heppu/go-template
OLD_CMD_DIR	:= cmd/demo
NEW_URL     ?= $(shell git config --get remote.origin.url)
NEW_MODULE  ?= $(subst ssh://,,$(subst .git,,$(subst https://,,${NEW_URL})))
NEW_NAME    ?= $(notdir ${NEW_MODULE})
NEW_CMD_DIR	:= cmd/${NEW_NAME}

.PHONY: rename
rename:
	@printf "Renaming project using following configuration:\n\n"
	@printf "URL:    ${NEW_URL}\n"
	@printf "MODULE: ${NEW_MODULE}\n"
	@printf "NAME:   ${NEW_NAME}\n\n"
	@go mod edit -module ${NEW_MODULE}
	@find . -type f -name '*.go' -exec sed -i 's|${OLD_MODULE}|${NEW_MODULE}|g' {} \;
	@mv ${OLD_CMD_DIR} ${NEW_CMD_DIR}
	@sed -i 's|${OLD_CMD_DIR}|${NEW_CMD_DIR}|g' ./applicationtest/application_test.go
	@printf "\nProject renamed succesfully, deleting rename.mk\n"
	@rm rename.mk
