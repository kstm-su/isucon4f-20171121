.PHONY: deploy
deploy: build
	supervisorctl restart isucon_go

.PHONY: fmt
fmt:
	go fmt *.go

.PHONY: build
build: fmt
	./build.sh

.PHONY: test
test:
	go test *.go
