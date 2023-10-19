VERSION ?= 0.0.2
OUTPUT_IMAGE_HTTP ?= gcr.io/solo-test-236622/sample-tap-server-http:${VERSION}
OUTPUT_IMAGE_GRPC ?= gcr.io/solo-test-236622/sample-tap-server-grpc:${VERSION}

all: sample-tap-server-http-docker sample-tap-server-grpc-docker

sample-tap-server-http-docker:
	docker build -t ${OUTPUT_IMAGE_HTTP} --build-arg=GO_MAIN_FILE=tap-server-http/sample-tap-server-http/main.go .

sample-tap-server-grpc-docker:
	docker build -t ${OUTPUT_IMAGE_GRPC} --build-arg GO_MAIN_FILE=tap-server-grpc/sample-tap-server-grpc/main.go .

run-sample-tap-server-http:
	cd tap-server-http/sample-tap-server-http/ && go mod tidy && go run main.go -p 9001

run-sample-tap-server-grpc:
	cd tap-server-grpc/sample-tap-server-grpc/ && go mod tidy && go run main.go -p 9001

run-sample-tap-server-http-docker: sample-tap-server-http-docker
	docker run --rm -it -p 9001:8080 gcr.io/solo-test-236622/sample-tap-server-http:0.0.2

run-sample-tap-server-grpc-docker: sample-tap-server-grpc-docker
	docker run --rm -it -p 9001:8080 gcr.io/solo-test-236622/sample-tap-server-grpc:0.0.2

push-docker-image:
	docker push ${OUTPUT_IMAGE}

print-%:
	@echo $($*)

