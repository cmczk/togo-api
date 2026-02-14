run-local:
	DATABASE_URL="postgres://togo_user:togo_password@localhost:5431/togo_db?sslmode=disable" PORT=3000 go run ./cmd/togo-api/main.go
