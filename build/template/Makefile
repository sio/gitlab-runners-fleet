# Configure this machine


export ANSIBLE_STDOUT_CALLBACK?=yaml
export PY_COLORS=1
export ANSIBLE_FORCE_COLOR=1


PLAYBOOKS=runner gateway
.PHONY: $(PLAYBOOKS)
$(PLAYBOOKS): venv
	$(VENV)/ansible-galaxy install -r requirements.yml
	$(VENV)/ansible-playbook --inventory inventory.yml $@.yml --extra-vars "ansible_python_interpreter=/usr/bin/python3"


include Makefile.venv
Makefile.venv:
	curl \
		-o Makefile.fetched \
		-L "https://github.com/sio/Makefile.venv/raw/v2022.07.20/Makefile.venv"
	echo "147b164f0cbbbe4a2740dcca6c9adb6e9d8d15b895be3998697aa6a821a277d8 *Makefile.fetched" \
		| sha256sum --check - \
		&& mv Makefile.fetched Makefile.venvk
