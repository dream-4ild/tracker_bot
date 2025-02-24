package main

import (
	"context"
	"github.com/go-telegram/bot"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(os.Getenv("TRACKER_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/hello", bot.MatchTypeExact, helloHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/new_task", bot.MatchTypePrefix, newTaskHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/update_task", bot.MatchTypePrefix, updateTaskHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypePrefix, listHandler)

	go notify(ctx, b)

	b.Start(ctx)
}
