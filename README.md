# go-ctrlc
Waits and gracefully stops golang program if CTRL+C called

Provides usefull CTRL+C interface, used to intercept ctrl+c signal
to stop program outside or by internal command

Usage:
```
func main() {
	var ctrlc CtrlC
	defer ctrlc.DeferThisToWaitCtrlC(true)

	....

	go some_your_logic()
	....

	ctrl.InterceptKill(true, func() {
		fmt.Println("software was stopped via Ctrl+C")
	})
}
```

