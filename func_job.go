package workerpool

type FuncJob func()

func (f FuncJob) Id() string {
	return ""
}

func (f FuncJob) Exec() {
	f()
}
