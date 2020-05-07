package task

type Repository interface {
	GetCurrent() Task
	Save(t Task)
}
