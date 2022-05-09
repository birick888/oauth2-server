BINARY=engine
test: 
	go test -v -cover -covermode=atomic ./...

engine:
	go build -o ${BINARY} app/*.go


unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker:
	docker build -t go_oauth2 .

run:
	docker-compose --env-file ./env/test.env up --build -d

remove:
	docker stop go_oauth2_api go_oauth2_mysql go_oauth2_redis
	docker rm go_oauth2_api go_oauth2_mysql go_oauth2_redis

stop:
	docker-compose --env-file ./env/test.env down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: clean install unittest build docker run stop vendor lint-prepare lint remove