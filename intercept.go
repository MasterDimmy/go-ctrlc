package ctrlc

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

/*
	Provides usefull CTRL+C interface, used to intercept ctrl+c signal
	to stop program outside or by internal command

	Usage:
	func main() {
		var ctrlc CtrlC
		defer ctrlc.DeferThisToWaitCtrlC(true)

		....
		go some_logic()
		....

		ctrl.InterceptKill(true, func() {
			fmt.Println("software was stopped via Ctrl+C")
		})
	}
*/

type CtrlC struct {
	m                       sync.Mutex
	force_stop_whole_system chan bool
}

func (c *CtrlC) init() {
	c.m.Lock()
	if c.force_stop_whole_system == nil {
		c.force_stop_whole_system = make(chan bool)
	}
	c.m.Unlock()
}

func (c *CtrlC) DeferThisToWaitCtrlC(stdout bool) {
	c.init()
	if stdout {
		fmt.Println("Waiting CTRL+C or ForceStopProgram")
	}
	for _ = range c.force_stop_whole_system {
	}
}

//stop program now
func (c *CtrlC) ForceStopProgram() {
	c.init()
	close(c.force_stop_whole_system)
}

//use  defer DeferThisToWaitCtrlC if you called InterceptKill !
//when called to stop via ctrl+c or ForceStop it call stop()
func (c *CtrlC) InterceptKill(stdout bool, stop func()) {
	c.init()

	s := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(s, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until we receive our signal.
	select {
	case v := <-s:
		if stdout {
			fmt.Printf("SYSTEM STOP: intercepted signal: %s\n", v.String())
		}
	case <-c.force_stop_whole_system:
		if stdout {
			fmt.Println("SYSTEM STOP: called ForceStopProgram()")
		}
	}

	if stop != nil {
		stop()
	}
	func() {
		defer func() {
			recover()
		}()
		close(c.force_stop_whole_system)
	}()
}