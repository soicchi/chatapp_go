path = ./...

container_up:
	docker compose --profile main up

go_get:
	docker compose run --rm app go get ${pkg}

go_fmt:
	docker compose run --rm app go fmt ${path}

go_tidy:
	docker compose run --rm app go mod tidy

go_vet:
	docker compose run --rm app go vet ${path}

go_test:
	docker compose up -d test-db
	docker compose run --rm app go test -v -cover ${path}
	docker compose down test-db
