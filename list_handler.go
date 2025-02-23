package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"regexp"
	"sort"
	"strings"
	"time"
)

type ByProject []Task

func (a ByProject) Len() int           { return len(a) }
func (a ByProject) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProject) Less(i, j int) bool { return a[i].project < a[j].project }

func listHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	re := regexp.MustCompile(`(?s)^/list(\s(\w+))?$`)
	matches := re.FindStringSubmatch(update.Message.Text)

	if matches == nil || len(matches) != 3 {
		defaultHandler(ctx, b, update)
		return
	}

	userId := update.Message.From.ID

	var projectName string
	if matches[2] != "" {
		projectName = matches[2]
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := func() (*sql.Rows, error) {
		if projectName != "" {
			return server.db.QueryContext(ctx,
				`SELECT * FROM tasks WHERE user_id = $1 AND project=$2`,
				userId,
				projectName,
			)
		} else {
			return server.db.QueryContext(ctx,
				`SELECT * FROM tasks WHERE user_id = $1`,
				userId,
			)
		}
	}()

	if err != nil {
		makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("smth wrong: %v", err))
		return
	}
	defer rows.Close()

	var activeTasks ByProject
	var closedTasks ByProject
	var backlogTasks ByProject

	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.id,
			&task.userId,
			&task.project,
			&task.task,
			&task.status,
			&task.deadline,
			&task.updatedAt,
		)

		if err != nil {
			makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("scan failed: %v", err))
			return
		}

		switch task.status {
		case newTask:
			activeTasks = append(activeTasks, task)
		case closeTask:
			if time.Since(task.updatedAt).Hours() <= 24*7 {
				closedTasks = append(closedTasks, task)
			}
		case backlogTask:
			backlogTasks = append(backlogTasks, task)
		default:
			makeResponse(ctx, b, update.Message.Chat.ID, fmt.Sprintf("unknown status: %v", task.status))
			return
		}
	}

	sort.Sort(activeTasks)
	sort.Sort(backlogTasks)
	sort.Sort(closedTasks)

	var response strings.Builder

	if len(activeTasks) > 0 {
		response.WriteString("Active tasks:\n")
	}

	var recProj string

	for _, task := range activeTasks {
		if recProj != task.project {
			response.WriteString(fmt.Sprintf(" %s\n", task.project))
			recProj = task.project
		}
		response.WriteString(fmt.Sprintf("  %v: %s", task.id, task.task))
		if task.deadline.Valid {
			response.WriteString(fmt.Sprintf(" deadline: %v", task.deadline.Time.Format(layout)))
		}
		response.WriteString("\n")
	}

	if len(backlogTasks) > 0 {
		response.WriteString("\n\nBacklog tasks:\n")
	}

	recProj = ""

	for _, task := range backlogTasks {
		if recProj != task.project {
			response.WriteString(fmt.Sprintf(" %s\n", task.project))
			recProj = task.project
		}
		response.WriteString(fmt.Sprintf("  %v: %s", task.id, task.task))
		if task.deadline.Valid {
			response.WriteString(fmt.Sprintf(" deadline: %v", task.deadline.Time.Format(layout)))
		}
		response.WriteString("\n")
	}

	if len(closedTasks) > 0 {
		response.WriteString("\n\nClosed tasks:\n")
	}

	recProj = ""

	for _, task := range closedTasks {
		if recProj != task.project {
			response.WriteString(fmt.Sprintf(" %s\n", task.project))
			recProj = task.project
		}
		response.WriteString(fmt.Sprintf("  %v: %s\n", task.id, task.task))
	}

	makeResponse(ctx, b, update.Message.Chat.ID, response.String())
}
