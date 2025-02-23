package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func makeResponse(ctx context.Context, b *bot.Bot, chatID int64, message string) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		makeResponse(ctx, b, update.Message.Chat.ID,
			fmt.Sprintf("Sorry, unknown command: *%s*", update.Message.Text))
	} else if update.EditedMessage != nil {
		makeResponse(ctx, b, update.EditedMessage.Chat.ID,
			"editing unsupported :(")
	}
}
