TAG = 0.0.8

default:
	go build
	./graphql-tester -H '$(GITHUB_TOKEN)' -u https://api.github.com/graphql github_api

release:
	@docker login --username $(DOCKER_USER) --password $(DOCKER_PASS)
	@docker push imega/graphql-tester:$(TAG)
	@docker push imega/graphql-tester:latest

build: buildfs
	@docker build -t imega/graphql-tester:$(TAG) .
	@docker tag imega/graphql-tester:$(TAG) imega/graphql-tester:latest

WD = /go/src/github.com/imega/graphql-tester
buildfs: unittest
	@sed -i -e 's/0.0.0/$(TAG)/g' $(CURDIR)/cmd/version.go
	@docker run -v $(CURDIR):$(WD) -w $(WD) golang:1.10-alpine go build -v -o src/graphql-tester
	@docker run --rm \
		-v $(CURDIR)/runner:/runner \
		-v $(CURDIR)/build:/build \
		-v $(CURDIR)/src:/src \
		-e TAG=$(TAG) \
		imega/base-builder:1.6.1 \
		--packages="musl busybox@main ca-certificates@main"

unittest:
	@docker run --rm -v $(CURDIR):$(WD) -w $(WD) \
		golang:1.10-alpine sh -c "go list ./... | grep -v 'vendor\|build\|runner\|src' | xargs go test -v"

test:
	@docker run -v $(CURDIR)/github_api:/data imega/graphql-tester:$(TAG) -H '$(HEADERS)' -u https://api.github.com/graphql /data
