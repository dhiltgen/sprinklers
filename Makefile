# Super simple makefile

# Base layer dependencies
GO_VERSION=1.13
ALPINE_VERSION=3.10

# TODO wire up versioning consistently...
#VERSION=
TAG=0.2.1 # TODO version based
ORG?=docker.io/dhiltgen
DAEMON_IMAGE=$(ORG)/sprinklerd:$(TAG)
CLIENT_IMAGE=$(ORG)/sprinklers:$(TAG)
DOCKER?=docker

# Using buildx for multi-arch
DOCKER_BUILD=$(DOCKER) buildx build --push \
    --build-arg GO_VERSION=$(GO_VERSION) \
    --build-arg ALPINE_VERSION=$(ALPINE_VERSION) \
    --platform linux/amd64,linux/arm/v7
# Using regular build without multi-arch
#DOCKER_BUILD=$(DOCKER) build --build-arg GO_VERSION=$(GO_VERSION) --build-arg ALPINE_VERSION=$(ALPINE_VERSION)

DOCKER_BUILDKIT=1
export DOCKER_BUILDKIT


build: daemon client

daemon:
	$(DOCKER_BUILD) \
	    -t $(DAEMON_IMAGE) \
	    --target daemon \
	    .

client:
	$(DOCKER_BUILD) \
	    -t $(CLIENT_IMAGE) \
	    --target client \
	    .

# Deploy the sprinkler service
deploy:
	kubectl apply -f ./sprinklerd.yaml

# Show the node port that was assigned to the service
# Example docker usage:
# docker run --rm -it dhiltgen/sprinklers --server $(make showport) list
showport:
	@echo "sprinklers:`kubectl get services sprinklerd -o jsonpath="{.spec.ports[0].nodePort}"`"

# Doesn't benefit from multi-arch at present
test:
	$(DOCKER) build \
	    --target test \
	    .

# Doesn't benefit from multi-arch at present
coverage:
	$(DOCKER) build \
	    --target coverage \
	    --output type=local,dest=./ \
	    .
	xdg-open ./cover.html

# Only used for non-multi-arch builds
push:
	$(DOCKER) push $(DAEMON_IMAGE)
	$(DOCKER) push $(CLIENT_IMAGE)

fmt:
	clang-format -i $$(find . -name \*.proto | grep -v "vendor/")

.PHONY: build daemon client test deploy coverage push fmt
