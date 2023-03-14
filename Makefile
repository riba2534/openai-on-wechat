run:
	go build -o openai-on-wechat
	./openai-on-wechat

build:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -o openai-on-wechat
	zip openai-on-wechat.zip openai-on-wechat config.json.example prompt.txt.example
	GOOS=windows GOARCH=amd64 go build -o openai-on-wechat-windows.exe
	zip openai-on-wechat-windows.zip openai-on-wechat-windows.exe config.json.example prompt.txt.example
	GOOS=darwin GOARCH=amd64 go build -o openai-on-wechat-darwin-amd64
	zip openai-on-wechat-darwin-amd64.zip openai-on-wechat-darwin-amd64 config.json.example prompt.txt.example
	GOOS=darwin GOARCH=arm64 go build -o openai-on-wechat-darwin-arm64
	zip openai-on-wechat-darwin-arm64.zip openai-on-wechat-darwin-arm64 config.json.example prompt.txt.example
	rm -rf openai-on-wechat openai-on-wechat-windows.exe openai-on-wechat-darwin-amd64 openai-on-wechat-darwin-arm64