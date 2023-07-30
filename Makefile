up:
	docker-compose -f ./deploy/docker-compose.yml up --build

down:
	docker-compose -f ./deploy/docker-compose.yml down

lint:
	golangci-lint run