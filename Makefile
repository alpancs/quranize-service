test:
	go test -race ./...

run:
	MONGODB_URL=mongodb://localhost/quranize go run main.go
