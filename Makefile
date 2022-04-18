PULUMI?=pulumi
.DEFAULT_GOAL=up


SETUP_PY=setup.cfg
include Makefile.venv

include makefiles/*.mk


.PHONY: up destroy
up destroy: | venv check-software
	$(VENV)/fleet-manager $@


.PHONY: check-software
check-software:
	$(PULUMI) version
	$(PY) --version
	@$(PY) -c 'import venv; import ensurepip'  # check for python3-venv package
	ssh -V
