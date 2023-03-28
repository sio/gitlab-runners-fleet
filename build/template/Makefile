# Configure this machine


export ANSIBLE_STDOUT_CALLBACK?=yaml
export PY_COLORS=1
export ANSIBLE_FORCE_COLOR=1
export LC_ALL=C.UTF-8
export PIP_NO_CACHE_DIR=off


PLAYBOOKS=runner gateway common
.PHONY: $(PLAYBOOKS)
$(PLAYBOOKS): | venv .ansible-galaxy-ready
	$(VENV)/ansible-playbook --inventory inventory/ $@.yml


.ansible-galaxy-ready: | venv
	$(VENV)/ansible-galaxy install -r requirements.yml
	touch $@


include Makefile.venv
Makefile.venv:
	curl \
		-o Makefile.fetched \
		-L "https://github.com/sio/Makefile.venv/raw/v2022.07.20/Makefile.venv"
	echo "147b164f0cbbbe4a2740dcca6c9adb6e9d8d15b895be3998697aa6a821a277d8 *Makefile.fetched" \
		| sha256sum --check - \
		&& mv Makefile.fetched Makefile.venv
