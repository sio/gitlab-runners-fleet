ifdef NON_INTERACTIVE
PULUMI_ARGS+=--yes --non-interactive
endif

GIT?=git
PULUMI?=pulumi

PULUMI_PROJECT=pulumi
PULUMI_STACK=main
PULUMI:=$(PULUMI) -C $(PULUMI_PROJECT)
PULUMI_STATE_DIRECTORY=$(PULUMI_PROJECT)/state
PULUMI_BACKEND_URL=file://$(realpath $(PULUMI_STATE_DIRECTORY))
export PULUMI_BACKEND_URL

VENVDIR=$(PULUMI_PROJECT)/venv
REQUIREMENTS_TXT=$(PULUMI_PROJECT)/requirements.txt


include makefiles/*.mk
include Makefile.venv


up: venv state-backend stack
	$(PULUMI) $(PULUMI_ARGS) $@


stack:
	-$(PULUMI) stack init --non-interactive $(PULUMI_STACK)


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
