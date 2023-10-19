VERSION ?= 0.0.2
OUTPUT_IMAGE ?= gcr.io/solo-test-236622/sample-tap-server-http:${VERSION}

all: sample-tap-server-docker

sample-tap-server-docker:
	docker build -t ${OUTPUT_IMAGE} .

push-docker-image:
	docker push ${OUTPUT_IMAGE}

print-%:
	@echo $($*)
