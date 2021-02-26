GIT?=git
PULUMI?=pulumi


include makefiles/*.mk
include Makefile.venv


.PHONY: pull
pull:
ifndef GIT_PULL_DISABLE
	$(GIT) pull --ff-only
endif
