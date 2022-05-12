PULUMI?=pulumi
.DEFAULT_GOAL=up


SETUP_PY=setup.cfg
include Makefile.venv

include makefiles/*.mk


# Add pulumi binary to path if it's not there yet
ifneq (./.,$(dir $(PULUMI)))
PATH:=$(dir $(PULUMI)):$(PATH)
endif


.PHONY: up down destroy
up down destroy: | venv check-software state-backend
	$(VENV)/fleet-manager $@


.PHONY: test
test: $(VENV)/tox
	$(VENV)/tox $(TOX_ARGS)


.PHONY: check-software
check-software:
	$(PULUMI) version
	$(PY) --version
	@$(PY) -c 'import venv; import ensurepip'  # check for python3-venv package
	ssh -V


.PHONY: clean
clean:
	git clean -idx


.PHONY: loc
loc: $(VENV)/pygount
	$(VENV)/pygount --format=summary .


PULUMI_STATE_DIRECTORY=pulumi-stack
export PULUMI_BACKEND_URL=file://$(realpath $(PULUMI_STATE_DIRECTORY))


$(PULUMI_STATE_DIRECTORY):
	mkdir -p "$@"


.PHONY: state-backend
state-backend: $(PULUMI_STATE_DIRECTORY)
	$(PULUMI) login
