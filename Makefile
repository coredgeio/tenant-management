VERSION ?= latest
REPO ?= coredgeio

GIT_TOKEN ?= ""

IMG = $(REPO)/tenant-management:$(VERSION)

.PHONY: all

all: build-image

build-image:
	go fmt ./...
	go vet ./...
ifeq ($(GIT_TOKEN),"")
	rsync -rupE $$HOME/.ssh/* .ssh/
	sudo docker build -t ${IMG} -f build/Dockerfile.ssh .
else
	sudo docker build --build-arg GIT_TOKEN="${GIT_TOKEN}" \
		-t ${IMG} -f build/Dockerfile .
endif

push-images:
	sudo docker push ${IMG}