package main

import (
	"notify"
	"notify/service/telegram"
)

func main() {
	notifier := notify.New()
	telegramService, _ := telegram.New("your_telegram_api_token")

	// Passing a telegram service. Ignoring error for demo simplicity
	telegramService.AddReceivers("-1234567890")

	notifier.UseServices(telegramService)

	_ = notifier.Send(
		"message subject1",
		"the actual message.hello world.")
}
