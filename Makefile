TEST?=$$(go list ./... | grep -v 'vendor')
# set the following values when using local-filesystem as a registry
#HOSTNAME=registry.terraform.io
#NAMESPACE=terraform-providers

HOSTNAME=terraform-registry.acme.com
NAMESPACE=acme
NAME=provisioner
VERSION=1.0.0
OS_ARCH=linux_amd64
PROTOCOL=x5
BUILD_DIR=binary_output
BINARY=terraform-provider-${NAME}
BINARY_OUT=${BUILD_DIR}/${BINARY}

default: install

build:
	go build -o ${BINARY_OUT}

release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY_OUT}_v${VERSION}_${PROTOCOL}
	zip -j ./bin/${BINARY_OUT}_${VERSION}_${OS_ARCH}.zip ./bin/${BINARY_OUT}_v${VERSION}_${PROTOCOL}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/v${VERSION}/${OS_ARCH}
	mkdir -p ~/.terraform.d/plugins/${OS_ARCH}
	cp ${BINARY_OUT} ~/.terraform.d/plugins/${OS_ARCH}/${BINARY}_v${VERSION}
	cp ${BINARY_OUT} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/v${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   