PULUMI_STATE_DIRECTORY=pulumi/state
PULUMI_BACKEND_URL=file://$(realpath $(PULUMI_STATE_DIRECTORY))
export PULUMI_BACKEND_URL


GIT?=git
PULUMI?=pulumi


include makefiles/*.mk
include Makefile.venv


$(PULUMI_STATE_DIRECTORY):
	mkdir -p "$@"


state-backend: $(PULUMI_STATE_DIRECTORY)
	$(PULUMI) login



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
