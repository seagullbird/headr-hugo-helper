APP?=hugo-helper
GOARCH?=amd64
GOOS?=linux
COMMIT?=$(shell git rev-parse --short HEAD)
IMAGE_NAME?=hugo-helper

clean:
	rm -f ${APP}

build: clean
	GOOS=${GOOS} GOARCH=${GOARCH} go build \
	-o ${APP}

container: build
	docker build -t ${IMAGE_NAME}:${COMMIT} .
