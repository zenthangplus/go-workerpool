package workerpool

type Job interface {
	Id() string
	Exec()
}
