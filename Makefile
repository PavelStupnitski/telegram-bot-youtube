.PHONY:
.SILENT:

build:
	go build -o ./.bin/telegramBotYouTube cmd/bot/main.go

run: build
	./.bin/telegramBotYouTube