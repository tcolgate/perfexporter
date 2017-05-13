// +build !linux

package system

func newProcessEventCounter(evtype uint32, ev uint64) (*perfCounter, error) {
	return nil, errrors.New("cannot count perf events on your platform")
}

func (c *perfCounter) read() (float64, error) {
	return 0, errrors.New("cannot count perf events on your platform")
}
