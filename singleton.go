package run

var singleton *Group = New()

func Always(runnables ...Runnable) {
	singleton.Always(runnables...)
}

func Add(when bool, runnables ...Runnable) {
	singleton.Add(when, runnables...)
}

func Run() error {
	return singleton.Run()
}

func Alive() bool {
	return singleton.Alive()
}
