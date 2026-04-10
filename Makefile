DB_URL=postgres://root:p@ssw0rd@localhost:5432/blogcms?sslmode=disable

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down

migrate-force:
	migrate -path db/migrations -database "$(DB_URL)" force $(version)

run:
	go run cmd/api/*.go
