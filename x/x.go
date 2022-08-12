package x

var (
	Hold  = false
	ExitC = make(chan bool)
)

func Close() {
	select {
	case <-ExitC:
	default:
		close(ExitC)
	}
}
