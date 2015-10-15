package readc

// #cgo LDFLAGS: -Llibomron -lomron -lm
// #include "readc.h"
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"strconv"
	"strings"
	"time"

	"github.com/thomasf/lg"
	"github.com/thomasf/bpchart/pkg/omron"
)


func Open() error {
	ret := C.m_open()
	str := C.GoString(ret)
	if str != "" {
		return errors.New(str)
	}
	return nil
}

func Close() error {
	ret := C.m_close()
	str := C.GoString(ret)
	if str != "" {
		return errors.New(str)
	}
	return nil
}

func Read(bank int) ([]omron.Entry, error) {
	counted := C.m_count(C.int(bank))
	var entries []omron.Entry
	for i := 0; i < int(counted)+1; i++ {
		ret := C.m_read(C.int(bank), C.int(i))
		str := C.GoString(ret)
		C.free(unsafe.Pointer(ret))
		if str == "" {
			fmt.Print(".")
		} else {
			fmt.Print("*")
			fields := strings.Split(str, ",")
			t, err := time.ParseInLocation("2006-01-02 15:04:05", fields[0], time.Local)
			if err != nil {
				lg.Fatal(err)
			}
			sys, err := strconv.Atoi(fields[1])
			if err != nil {
				lg.Fatal(err)
			}
			dia, err := strconv.Atoi(fields[2])
			if err != nil {
				lg.Fatal(err)
			}
			pulse, err := strconv.Atoi(fields[3])
			if err != nil {
				lg.Fatal(err)
			}
			entry := omron.Entry{
				Time:  t,
				Sys:   sys,
				Dia:   dia,
				Pulse: pulse,
				Bank:  bank,
			}
			entries = append(entries, entry)
		}
	}
	fmt.Print("\n")

	return entries, nil
}
