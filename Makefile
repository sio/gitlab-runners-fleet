GIT?=git
PULUMI?=pulumi


include makefiles/*.mk
include Makefile.venv


.PHONY: pull
pull:
ifndef GIT_PULL_DISABLE
	$(GIT) pull --ff-only
endif


.PHONY: check-software
check-software:
	$(PULUMI) version
	$(PY) --version
	$(PY) -c 'import venv; import ensurepip'  # check for python3-venv package
