package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func updateTaskHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	re := regexp.MustCompile(`(?s)^/update_task (\d+)\s(\w+)\s(.+)$`)
	matches := re.FindStringSubmatch(update.Message.Text)

	if matches == nil || len(matches) != 4 {
		defaultHandler(ctx, b, update)
		return
	}

	userId := update.Message.From.ID

	taskId, err := strconv.Atoi(matches[1])
	if err != nil {
		makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Invalid task id: %s", matches[1]))
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := server.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		makeResponse(ctx, b, update.Message.Chat.ID, "Internal Server Error")
		return
	}

	var task Task
	err = tx.QueryRowContext(ctx,
		`SELECT * FROM tasks where id=$1 and user_id=$2;`,
		taskId, userId,
	).Scan(&task.id, &task.userId, &task.project, &task.task, &task.status, &task.deadline, &task.updatedAt)

	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("обманывать плохо! у тебя нет такой задачи"))
		} else {
			makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("smth wrong: %v", err))
		}
		return
	}

	collumn := matches[2]
	if _, exist := mutableCollumns[collumn]; !exist {
		makeResponse(ctx, b, update.Message.Chat.ID, "bad collumn")
		return
	}

	newValue := matches[3]
	switch collumn {
	case project:
		task.project = newValue
	case description:
		task.task = newValue
	case status:
		task.status = newValue
	case deadline:
		timeDeadline, err := time.Parse(layout, newValue)
		if err != nil {
			makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("smth wrong: %v", err))
			return
		}

		task.deadline = sql.NullTime{Time: timeDeadline, Valid: true}

	}

	_, err = tx.ExecContext(ctx,
		`UPDATE tasks SET project=$1, task=$2, status=$3, deadline=$4 WHERE id=$5`,
		task.project,
		task.task,
		task.status,
		task.deadline,
		task.id)

	if err != nil {
		tx.Rollback()
		makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("smth wrong: %v", err))
		return
	}

	tx.Commit()

	makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Successfuly update!"))
}
