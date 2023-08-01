up:
	docker-compose -f ./deploy/docker-compose.yml up --build

down:
	docker-compose -f ./deploy/docker-compose.yml down

lint:
	golangci-lint run

go-mod:
	go mod download && go mod tidy


test:
	go test ./...

pre-push: go-mod test lint
