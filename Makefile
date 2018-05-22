test:
	CORPUS_PATH=../corpus/ go test ./...

coverage:
	CORPUS_PATH=../corpus/ go test -coverprofile=coverage.out ./quran && go tool cover -html=coverage.out

run:
	MONGODB_URL=mongodb://localhost/quranize go run main.go
