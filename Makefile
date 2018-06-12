build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/twitter-crawler twitter-crawler/main.go