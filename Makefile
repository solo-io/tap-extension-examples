VERSION ?= 0.0.2
OUTPUT_IMAGE_HTTP ?= gcr.io/solo-test-236622/sample-tap-server-http:${VERSION}
OUTPUT_IMAGE_GRPC ?= gcr.io/solo-test-236622/sample-tap-server-grpc:${VERSION}
TAP_SERVER_PORT ?= 9001

all: sample-tap-server-http-docker sample-tap-server-grpc-docker

sample-tap-server-http-docker:
	docker build -t ${OUTPUT_IMAGE_HTTP} --build-arg=GO_MAIN_FILE=tap-server-http/cmd/main.go .

sample-tap-server-grpc-docker:
	docker build -t ${OUTPUT_IMAGE_GRPC} --build-arg GO_MAIN_FILE=tap-server-grpc/cmd/main.go .

run-sample-tap-server-http:
	cd tap-server-http/cmd/ && go mod tidy && go run main.go -p ${TAP_SERVER_PORT}

run-sample-tap-server-grpc:
	cd tap-server-grpc/cmd/ && go mod tidy && go run main.go -p ${TAP_SERVER_PORT}

run-sample-tap-server-http-docker: sample-tap-server-http-docker
	docker run --rm -it -p ${TAP_SERVER_PORT}:8080 ${OUTPUT_IMAGE_HTTP}

run-sample-tap-server-grpc-docker: sample-tap-server-grpc-docker
	docker run --rm -it -p ${TAP_SERVER_PORT}:8080 ${OUTPUT_IMAGE_GRPC}

push-docker-image-http: sample-tap-server-http-docker
	docker push ${OUTPUT_IMAGE_HTTP}

push-docker-image-grpc: sample-tap-server-grpc-docker
	docker push ${OUTPUT_IMAGE_GRPC}

push-docker-images: push-docker-image-http push-docker-image-grpc

print-%:
	@echo $($*)
