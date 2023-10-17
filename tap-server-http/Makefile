VERSION = 0.0.1
# OUTPUT_IMAGE = quay.io/solo-io/sample-tap-server:${VERSION}
OUTPUT_IMAGE = gcr.io/solo-test-236622/sample-tap-server-http:${VERSION}

all: sample-tap-server-docker load-kind-image

sample-tap-server-docker:
	docker build -t ${OUTPUT_IMAGE} .

load-kind-image: sample-tap-server-docker
	kind load docker-image ${OUTPUT_IMAGE}

push-docker-image:
	docker push ${OUTPUT_IMAGE}
