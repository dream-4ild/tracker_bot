package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"log"
	"time"
)

func notify(ctx context.Context, b *bot.Bot) {
	now := time.Now().Truncate(interval * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			rows, err := server.db.QueryContext(context.Background(),
				`SELECT * FROM tasks WHERE (status=$1 OR status=$2) AND deadline IS NOT NULL AND deadline - (NOW() + interval '3 hours') < interval '24 hour' AND NOW() + interval '3 hours' < deadline;`,
				newTask,
				backlogTask,
			)

			if err != nil {
				log.Println(err)
				goto suspend
			}

			for rows.Next() {
				var task Task
				err = rows.Scan(
					&task.id,
					&task.userId,
					&task.project,
					&task.task,
					&task.status,
					&task.deadline,
					&task.updatedAt,
				)
				go makeResponse(ctx,
					b,
					task.userId,
					fmt.Sprintf(
						"Deadline soon!\nTask from %s:\n%s;\ndeadline: %v",
						task.project,
						task.task,
						task.deadline.Time.Format(layout),
					),
				)
			}

		}
	suspend:
		now = now.Add(interval * time.Minute)

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Until(now)):
			continue
		}

	}
}
