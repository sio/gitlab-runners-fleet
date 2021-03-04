ifdef NON_INTERACTIVE
PULUMI_ARGS+=--skip-preview --yes --non-interactive
endif

GIT?=git
PULUMI?=pulumi

PULUMI_PROJECT=pulumi
PULUMI_STACK=main
PULUMI_STATE_DIRECTORY=$(PULUMI_PROJECT)/state
PULUMI_BACKEND_URL=file://$(realpath $(PULUMI_STATE_DIRECTORY))
export PULUMI_BACKEND_URL
PULUMI_SNAPSHOT_OBJECT?=snapshot
export PULUMI_SNAPSHOT_OBJECT
HCLOUD_USERDATA_TEMPLATE?=$(abspath instance/cloudinit.yml.j2)
export HCLOUD_USERDATA_TEMPLATE
GITLAB_RUNNER_SSHKEY?=$(abspath $(PULUMI_PROJECT)/keys/runner)
export GITLAB_RUNNER_SSHKEY

VENVDIR=$(PULUMI_PROJECT)/venv
REQUIREMENTS_TXT=$(PULUMI_PROJECT)/requirements.txt
.DEFAULT_GOAL=up


include makefiles/*.mk
include Makefile.venv


PULUMI:=$(PULUMI) -C $(PULUMI_PROJECT)
.PHONY: up destroy
up destroy: pull venv check-software stack assign $(GITLAB_RUNNER_SSHKEY)
	$(PULUMI) $(PULUMI_ARGS) $@
	$(PULUMI) config set $(PULUMI_SNAPSHOT_OBJECT) "$$($(PULUMI) stack output $(PULUMI_SNAPSHOT_OBJECT) --json)"


.PHONY: assign
assign: venv
	$(VENV)/python pulumi/assign_runners.py


.PHONY: stack
stack: state-backend
	-$(PULUMI) stack init --non-interactive $(PULUMI_STACK)
	$(PULUMI) stack output $(PULUMI_SNAPSHOT_OBJECT) >/dev/null 2>/dev/null \
	|| $(PULUMI) config set $(PULUMI_SNAPSHOT_OBJECT) '{}'


$(PULUMI_STATE_DIRECTORY):
	mkdir -p "$@"


.PHONY: state-backend
state-backend: $(PULUMI_STATE_DIRECTORY)
	$(PULUMI) login


$(GITLAB_RUNNER_SSHKEY):
	mkdir -p "$(dir $@)"
	ssh-keygen -q -t ed25519 -a 10 -f "$(GITLAB_RUNNER_SSHKEY)" -C ci-runner-key -N ""


.PHONY: show
show:
	$(PULUMI) stack


.PHONY: pull
pull:
ifndef GIT_PULL_DISABLE
	-$(GIT) pull --ff-only
endif


.PHONY: check-software
check-software:
	$(PULUMI) version
	$(PY) --version
	@$(PY) -c 'import venv; import ensurepip'  # check for python3-venv package
