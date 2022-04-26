# Initialize GitLab runner instance on localhost


.PHONY: runner router
runner router: venv
	$(VENV)/ansible-galaxy install -r requirements.yml
	$(VENV)/ansible-playbook $@.yml --extra-vars "ansible_python_interpreter=/usr/bin/python3"


include Makefile.venv
Makefile.venv:
	curl \
		-o Makefile.fetched \
		-L "https://github.com/sio/Makefile.venv/raw/v2022.04.13/Makefile.venv"
	echo "bb2e61acbc3a8ea83011a43e0b010b331eddec5ac24d4edce2b7b427362460c9 *Makefile.fetched" \
		| sha256sum --check - \
		&& mv Makefile.fetched Makefile.venv
