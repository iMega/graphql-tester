TAG = 0.0.1

release:
	@docker login --username $(DOCKER_USER) --password $(DOCKER_PASS)
	@docker push imega/graphql-tester:$(TAG)
	@docker push imega/graphql-tester:latest

build: buildfs
	@docker build -t imega/graphql-tester:$(TAG) .
	@docker tag imega/graphql-tester:$(TAG) imega/graphql-tester:latest

buildfs: WD = /go/src/github.com/imega/graphql-tester
buildfs:
	@docker run -v $(CURDIR):$(WD) -w $(WD) golang:1.10-alpine go build -v -o src/graphql-tester
	@docker run --rm \
		-v $(CURDIR)/runner:/runner \
		-v $(CURDIR)/build:/build \
		-v $(CURDIR)/src:/src \
		-e TAG=$(TAG) \
		imega/base-builder:1.6.1 \
		--packages="musl busybox@main ca-certificates@main"

test:
	@docker run -v $(CURDIR)/github_api:/data imega/graphql-tester:0.0.1 -H '$(HEADERS)' -u https://api.github.com/graphql /data
