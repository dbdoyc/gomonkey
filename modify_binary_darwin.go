package gomonkey

import (
	"syscall"
	"time"
)

func modifyBinary(target uintptr, bytes []byte) {
	function := entryAddress(target, len(bytes))
	err := mprotectCrossPage(target, len(bytes), syscall.PROT_READ|syscall.PROT_WRITE)
	if err != nil {
		panic(err)
	}
	copy(function, bytes)
	err = mprotectCrossPage(target, len(bytes), syscall.PROT_READ|syscall.PROT_EXEC)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond)
}

func mprotectCrossPage(addr uintptr, length int, prot int) error {
	pageSize := syscall.Getpagesize()
	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := entryAddress(p, pageSize)
		if err := syscall.Mprotect(page, prot); err != nil {
			return err
		}
	}
	return nil
}
