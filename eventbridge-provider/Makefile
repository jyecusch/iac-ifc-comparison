NAMESPACE := nitric
NAME := awseventbridge
VERSION := 0.0.1

# Don't modify below this line
FQN := ${NAMESPACE}/${NAME}

RUNTIME_BINARY := bin/${NAMESPACE}-${NAME}
RUNTIME_EMBED := cmd/deploy/runtime-extension

binaries: deploybin

runtimebin:
	@echo Building ${FQN} runtime server
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${RUNTIME_BINARY} -ldflags="-s -w -extldflags=-static" ./cmd/runtime

# move runtime binary so it can be embedded into the deployment engine binary
predeploybin: runtimebin
	@cp ${RUNTIME_BINARY} ${RUNTIME_EMBED}

deploybin: predeploybin
	@echo Building ${FQN} deployment server
	@CGO_ENABLED=0 go build -o bin/deploy-extension -ldflags="-s -w -extldflags=-static" -ldflags="-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore" ./cmd/deploy
	@rm ${RUNTIME_EMBED}

.PHONY: install
install: deploybin
	@echo installing ${FQN} to ${HOME}/.nitric/providers/${NAMESPACE}/${NAME}-${VERSION}
	@mkdir -p ${HOME}/.nitric/providers/${NAMESPACE}/
	@rm -f ${HOME}/.nitric/providers/${NAMESPACE}/${NAME}-${VERSION}
	@cp bin/deploy-extension ${HOME}/.nitric/providers/${NAMESPACE}/${NAME}-${VERSION}
	@echo "Installed successfully\n\nUpdate your stack file to use this provider.\n provider: ${FQN}@${VERSION}"

.PHONY: clean
clean:
	@rm -rf bin