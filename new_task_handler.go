package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"regexp"
	"time"
)

func newTaskHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	re := regexp.MustCompile(`(?s)^/new_task\s(\p{L}+)(\s(\d{2}\.\d{2}\.\d{4}))?\s(.+)$`)
	matches := re.FindStringSubmatch(update.Message.Text)
	if matches == nil || len(matches) != 5 {
		defaultHandler(ctx, b, update)
		return
	}

	var err error

	userId := update.Message.From.ID

	var project any = sql.Null[string]{}
	if matches[1] != backlogTask {
		project = matches[1]
	}

	var deadline any = sql.Null[time.Time]{}

	if matches[3] != "" {

		deadline, err = time.Parse(layout, matches[3])
		if err != nil {
			makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("smth wrong: %v", err))
			return
		}
	}

	task := matches[4]

	taskType := newTask
	if _, err := project.(string); !err {
		taskType = backlogTask
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = server.db.Exec(
		`INSERT INTO tasks (user_id, project, task, status, deadline) values ($1, $2, $3, $4, $5);`,
		userId, project, task, taskType, deadline,
	)
	if err != nil {
		makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("smth wrong: %v", err))
		return
	}

	makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Successfuly create!"))
}
