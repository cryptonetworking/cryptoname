package handy

import "runtime"

func Loop() {
	for {
		runtime.Gosched()
	}
}
