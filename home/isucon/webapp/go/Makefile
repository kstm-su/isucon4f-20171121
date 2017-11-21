.PHONY: deploy
deploy: build
	ssh root@app2 supervisorctl stop isucon_go
	ssh root@app3 supervisorctl stop isucon_go
	scp golang-webapp root@app2:/home/isucon/webapp/go/
	scp golang-webapp root@app3:/home/isucon/webapp/go/
	ssh root@app2 supervisorctl start isucon_go
	ssh root@app3 supervisorctl start isucon_go
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
