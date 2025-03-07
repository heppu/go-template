URL         ?= $(shell git config --get remote.origin.url)
MODULE      ?= $(subst ssh://,,$(subst .git,,$(subst https://,,$(URL))))
NAME        ?= $(notdir $(MODULE))
CMD_DIR     ?= cmd/$(NAME)
OLD_MODULE  := github.com/heppu/go-template
OLD_CMD_DIR	:= cmd/app

.PHONY: rename
rename:
	@printf "Renaming project using following configuration:\n\n"
	@printf "URL:    $(URL)\n"
	@printf "MODULE: $(MODULE)\n"
	@printf "NAME:   $(NAME)\n\n"
	@go mod edit -module $(MODULE)
	@find . -type f -name '*.go' -exec sed -i 's|${OLD}|${MODULE}|g' {} \;
	@mv cmd/app/ cmd/$(NAME)
	@sed -i 's|${OLD_CMD_DIR}|${CMD_DIR}|g' ./applicationtest/app_test.go
	@printf "\nProject renamed succesfully, deleting rename.mk\n"
	@rm rename.mk
