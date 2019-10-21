.PHONY: build
build:
	env GOOS=linux GOARCH=amd64 go build -o loki_conn main.go 
	docker build -t loki_conn .

.PHONY: run
run:
	go run main.go -url="http://127.0.0.1:3100/api/prom/push" -interval=10s

.PHONY: push
push:
	docker tag loki_conn hellodudu86/loki_conn
	docker push hellodudu86/loki_conn


