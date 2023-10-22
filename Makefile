path = ./...

go_get:
	docker compose run --rm app go get ${pkg}

go_fmt:
	docker compose run --rm app go fmt ${path}

go_tidy:
	docker compose run --rm app go mod tidy

go_vet:
	docker compose run --rm app go vet ${path}

go_test:
	docker compose run --rm app go test -v ${path}
