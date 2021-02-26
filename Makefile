ifdef NON_INTERACTIVE
PULUMI_ARGS+=--skip-preview --yes --non-interactive
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
.DEFAULT_GOAL=up


include makefiles/*.mk
include Makefile.venv


.PHONY: up destroy
up destroy: venv check-software state-backend stack
	$(PULUMI) $(PULUMI_ARGS) $@


.PHONY: stack
stack:
	-$(PULUMI) stack init --non-interactive $(PULUMI_STACK)


$(PULUMI_STATE_DIRECTORY):
	mkdir -p "$@"


.PHONY: state-backend
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
	@$(PY) -c 'import venv; import ensurepip'  # check for python3-venv package
