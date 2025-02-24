package main

const (
	newTask     = "active"
	closeTask   = "close"
	backlogTask = "backlog"
)

const (
	layout = "02.01.2006"
)

const (
	project     = "project"
	description = "task"
	status      = "status"
	deadline    = "deadline"
)

var mutableCollumns = map[string]struct{}{project: {}, description: {}, status: {}, deadline: {}}

const interval = 60 * 5
