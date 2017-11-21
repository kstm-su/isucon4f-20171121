.PHONY: deploy
deploy: build
	supervisorctl restart isucon_go

.PHONY: build
build: golang-webapp
	./build.sh
