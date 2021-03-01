# Initialize GitLab runner instance on localhost


.DEFAULT_GOAL=runner
REQUIREMENTS_TXT=requirements.txt


include Makefile.venv
Makefile.venv:
	curl \
		-o Makefile.fetched \
		-L "https://github.com/sio/Makefile.venv/raw/v2020.08.14/Makefile.venv"
	echo "5afbcf51a82f629cd65ff23185acde90ebe4dec889ef80bbdc12562fbd0b2611 *Makefile.fetched" \
		| sha256sum --check - \
		&& mv Makefile.fetched Makefile.venv


FILES=playbook.yml requirements.yml requirements.txt
$(FILES):
	curl -o $@ -L "https://github.com/sio/gitlab-runners-fleet/raw/master/$@"


.PHONY: runner
runner: venv $(FILES)
	$(VENV)/ansible-galaxy install -r requirements.yml
	$(VENV)/ansible-playbook playbook.yml
