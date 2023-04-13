DOCKER?=docker
DOCKER_FILE?=Dockerfile
DOCKER_CONTEXT=..
DOCKER_REGISTRY?=ghcr.io
DOCKER_REGISTRY_PASSWD?=
DOCKER_USER?=sio
DOCKER_REPO?=ghcr.io/sio/gitlab-runners-fleet
DOCKER_TAG?=v3
export DOCKER_REGISTRY_PASSWD


# Enable BuildKit for --chmod
export DOCKER_BUILDKIT=1


.PHONY: docker-build
docker-build:
	$(DOCKER) build --pull --tag "$(DOCKER_REPO):$(DOCKER_TAG)" --file $(DOCKER_FILE) $(DOCKER_CONTEXT)


.PHONY: docker-push
docker-push:
	echo $$DOCKER_REGISTRY_PASSWD | $(DOCKER) login -u $(DOCKER_USER) --password-stdin $(DOCKER_REGISTRY)
	$(DOCKER) push "$(DOCKER_REPO):$(DOCKER_TAG)"
