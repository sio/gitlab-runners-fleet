.PHONY: package build
package build: dist


dist: setup.cfg pyproject.toml Makefile
dist: src
dist: | $(VENV)/build
	-$(RM) -rv dist
	$(VENV)/python -m build


.PHONY: upload
upload: dist | $(VENV)/twine
	$(VENV)/twine upload --repository testpypi $(TWINE_ARGS) dist/*
	$(VENV)/twine upload $(TWINE_ARGS) dist/*
