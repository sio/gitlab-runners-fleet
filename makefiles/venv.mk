Makefile.venv:
	curl \
		-o Makefile.fetched \
		-L "https://github.com/sio/Makefile.venv/raw/v2022.04.13/Makefile.venv"
	echo "bb2e61acbc3a8ea83011a43e0b010b331eddec5ac24d4edce2b7b427362460c9 *Makefile.fetched" \
		| sha256sum --check - \
		&& mv Makefile.fetched Makefile.venv
