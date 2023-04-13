.PHONY: run
run:
	go run cmd/bot/main.go

.PHONY: build
build:
	go build -o ./bin/bot$(shell go env GOEXE) cmd/bot/main.go

.PHONY: docker
docker:
	docker build . -t hablof/logistic-package-api-bot

.PHONY: docker-run
docker-run:
	docker run -d --env-file ./.env \
		-v $(PWD)/config.yml:/root/config.yml \
		hablof/logistic-package-api-bot
