OS=linux
ARCH=amd64

build:
	env GOOS=$(OS) GOARCH=$(ARCH) go build -o dumper.$(OS)-$(ARCH)
