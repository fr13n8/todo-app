build:
	$(MAKE) test && docker compose build todo-app

run:
	docker compose up todo-app

test:
	go test -v ./...

migrate_up:
	migrate -path ./schema -database 'postgres://postgres:1234@localhost:5436/postgres?sslmode=disable' up

migrate_down:
	migrate -path ./schema -database 'postgres://postgres:1234@localhost:5436/postgres?sslmode=disable' down