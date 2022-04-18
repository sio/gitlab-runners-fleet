PULUMI?=pulumi
.DEFAULT_GOAL=up


include makefiles/*.mk
include Makefile.venv


.PHONY: up destroy
up destroy: | venv check-software
	$(VENV)/fleet-manager $@


.PHONY: check-software
check-software:
	$(PULUMI) version
	$(PY) --version
	@$(PY) -c 'import venv; import ensurepip'  # check for python3-venv package
	ssh -V
