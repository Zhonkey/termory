ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
CMD ?= ./main.go
build:
	docker compose build
up: build
	docker compose up -d
	#docker compose exec app go mod tidy
sh: up
	docker compose exec app sh -c "/bin/bash"
npm: up
	docker compose exec vue sh -c "/bin/bash"
db: up
	docker compose exec db psql -U trainer -d trainer
down:
	docker compose down --remove-orphans
debug:
	docker compose exec app sh -c 'go build -gcflags "all=-N -l" -o main_debug.bin .'
	docker compose exec app sh -c 'dlv exec ./main_debug.bin --headless --listen=:2345 --api-version=2 --accept-multiclient'
