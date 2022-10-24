DOCKER?=docker
DOCKER_FILE?=container/Dockerfile
DOCKER_CONTEXT=.
DOCKER_REGISTRY?=ghcr.io
DOCKER_REGISTRY_PASSWD?=
DOCKER_USER?=sio
DOCKER_REPO?=ghcr.io/sio/gitlab-runners-fleet
DOCKER_TAG?=latest
export DOCKER_REGISTRY_PASSWD


.PHONY: docker-build
docker-build:
	$(DOCKER) build --pull --tag "$(DOCKER_REPO):$(DOCKER_TAG)" --file $(DOCKER_FILE) $(DOCKER_CONTEXT)


.PHONY: docker-push
docker-push:
	echo $$DOCKER_REGISTRY_PASSWD | $(DOCKER) login -u $(DOCKER_USER) --password-stdin $(DOCKER_REGISTRY)
	$(DOCKER) push "$(DOCKER_REPO):$(DOCKER_TAG)"
