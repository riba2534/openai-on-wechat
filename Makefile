run:
	go build -o openai-on-wechat
	./openai-on-wechat

build:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -o openai-on-wechat
	zip openai-on-wechat.zip openai-on-wechat config.json.example prompt.txt.example