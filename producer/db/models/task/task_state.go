package task

type TaskState string

const (
	Pending    TaskState = "pending"
	InProgress TaskState = "in_progress"
	Completed  TaskState = "completed"
	Failed     TaskState = "failed"
)
