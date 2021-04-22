GO111MODULE := on
DOCKER_TAG := $(or ${GITHUB_TAG_NAME}, latest)

all: partition-watchdog

.PHONY: partition-watchdog
partition-watchdog:
	go build -o bin/partition-watchdog
	strip bin/partition-watchdog

.PHONY: dockerimages
dockerimages:
	docker build -t mwennrich/partition-watchdog:${DOCKER_TAG} .

.PHONY: dockerpush
dockerpush:
	docker push mwennrich/partition-watchdog:${DOCKER_TAG}

.PHONY: clean
clean:
	rm -f bin/*

