// +build linux

package system

import (
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type perfEventAttr struct {
	Type                       uint32
	Size                       uint32
	Config                     uint64
	UnionSamplePeriodFreq      uint64
	SampleType                 uint64
	ReadFormat                 uint64
	Flags                      uint64
	UnionWakeupEventsWaterMark uint32
	BPType                     uint32
	UnionBPaddrConfig1         uint64
	UnionBPLenConfig2          uint64
	BranchSampleType           uint64
	SampleRegsUser             uint64
	SampleStackUser            uint32
	ClockID                    uint32
	SampleRegsIntr             uint64
	AuxWatermark               uint32
	SampleMaxStack             uint16
	Reserved2                  uint16
}

const evSzV5 = 112

type perfEvReadData struct {
	Count uint64
}

func (c *Counter) read() (float64, error) {
	var v perfEvReadData
	err := binary.Read(c.fd, binary.LittleEndian, &v)
	return float64(v.Count), err
}

func newProcessEventCounter(evtype uint32, ev uint64) (*os.File, error) {
	attrs := &perfEventAttr{
		Type:   evtype,
		Config: ev,
		Flags:  perfFlagExcludeKernel | perfFlagExcludeHv | perfFlagExcludeIdle,
	}
	fd, err := perfProcessEventOpen(attrs)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func perfProcessEventOpen(a *perfEventAttr) (*os.File, error) {
	pid := int64(syscall.Getpid())

	return perfEventOpen(a, pid, -1, -1, 0)
}

func perfEventOpen(a *perfEventAttr, pid, cpu, groupFD int64, fdFlags uint64) (*os.File, error) {
	a.Size = evSzV5
	a.ReadFormat = 0

	fd, _, err := syscall.Syscall6(
		syscall.SYS_PERF_EVENT_OPEN,
		uintptr(unsafe.Pointer(a)),
		uintptr(pid),
		uintptr(cpu),
		uintptr(groupFD),
		uintptr(fdFlags),
		0,
	)
	if err != 0 {
		return nil, os.NewSyscallError("perf counter open", err)
	}

	name := fmt.Sprintf("<perf counter fd=%d>", fd)
	return os.NewFile(fd, name), nil
}
