package cpuset

import (
	handy "github.com/itsabgr/go-handy"
	unix "golang.org/x/sys/unix"
	"runtime"
)

func CPUGet(tid int) []int {
	if tid < 0 {
		tid = unix.Gettid()
	}
	set := &unix.CPUSet{}
	set.Zero()
	err := unix.SchedGetaffinity(tid, set)
	handy.Throw(err)
	result := make([]int, 0, runtime.NumCPU())
	for i := range handy.N(uint(runtime.NumCPU())) {
		if set.IsSet(i) {
			result = append(result, i)
		}
	}
	return result
}
func CPUSet(tid int, ns ...int) error {
	if tid < 0 {
		tid = unix.Gettid()
	}
	set := &unix.CPUSet{}
	set.Zero()
	for _, n := range ns {
		set.Set(n)
	}
	return unix.SchedSetaffinity(tid, set)
}
