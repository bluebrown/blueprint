# export these in your shell and run
# make with -e to override them
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64

# build and dist output
dist_dir=.dist
bin_dir=.bin

# project nam and version
name=blueprint
version=0.1.0-alpha


# local build

install: build
	sudo cp -f $(bin_dir)/$(name)-$(version)-$(GOARCH)-static /usr/local/bin/$(name)

build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a \
		-ldflags '-X main.version=$(version) -X main.buildDate=$(shell date +%Y-%m-%d)' \
		-o $(bin_dir)/$(name)-$(version)-$(GOARCH)-static \
		-installsuffix static \
		./cmd/

dist:
	mkdir -p $(dist_dir)
	tar -czf $(dist_dir)/$(name)-$(GOARCH)-static.tar.gz $(bin_dir)/$(name)-$(version)-$(GOARCH)-static

clean:
	rm -rf $(bin_dir)
	rm -rf $(dist_dir)


# docker image

image.build:
	docker build \
		--tag bluebrown/$(name):latest \
		--tag bluebrown/$(name):$(version) \
		--build-arg version=$(version) \
		--build-arg buildDate=$(shell date +%Y-%m-%d) \
		--file assets/docker/Dockerfile .

image.push:
	docker push bluebrown/$(name):latest
	docker push bluebrown/$(name):$(version)


# example

example.local: example.clean
	go run ./cmd/blueprint/ assets/example/repo/ assets/myrepo --set service.enabled=false -f assets/example/v.yaml --values assets/example/v2.yaml

example.upstream: example.clean
	go run ./cmd/blueprint/ https://github.com/bluebrown/$(name)-example assets/myrepo \
	-f https://gist.githubusercontent.com/bluebrown/be3637e410f80b6c0eeeae5dc95ca662/raw/8be2a8e561463ac97448930a146bc0a9095616d4/$(name)-example-values.yaml

example.clean:
	rm -rf assets/myrepo
