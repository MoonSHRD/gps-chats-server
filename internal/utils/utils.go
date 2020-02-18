package utils

import "time"

func SetInterval(someFunc func(...interface{}), milliseconds int, async bool, args ...interface{}) chan bool {

	// How often to fire the passed in function
	// in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	callFuncionWithArgs := func() {
		switch len(args) {
		case 0:
			{
				someFunc()
			}
		case 1:
			{
				someFunc(args[0])
			}
		case 2:
			{
				someFunc(args[0], args[1])
			}
		case 3:
			{
				someFunc(args[0], args[1], args[3])
			}
			// add more cases if you need this
		}
	}

	// Put the selection in a go routine
	// so that the for loop is none blocking
	go func() {
		for {

			select {
			case <-ticker.C:
				if async {
					// This won't block
					go callFuncionWithArgs()
				} else {
					// This will block
					callFuncionWithArgs()
				}
			case <-clear:
				ticker.Stop()
				return
			}

		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear

}
