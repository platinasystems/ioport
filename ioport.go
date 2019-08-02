// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package ioport

import (
	"errors"
	"os"
	"syscall"
)

const (
	sys_iopl   = 172 //amd64
	sys_ioperm = 173 //amd64
)

// Outb writes a byte to the specified ioport.
func Outb(addr uint16, data byte) (err error) {
	if err = setIoperm(addr); err != nil {
		return err
	}

	f, err := os.OpenFile("/dev/port", os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	if _, err = f.Seek(int64(addr), 0); err != nil {
		return
	}
	if _, err = f.Write([]byte{data}); err != nil {
		return
	}
	f.Sync()

	if err = clrIoperm(addr); err != nil {
		return err
	}
	return nil
}

// Inb reads a byte from the specified ioport
func Inb(addr uint16) (data byte, err error) {
	if err = setIoperm(addr); err != nil {
		return
	}

	f, err := os.Open("/dev/port")
	if err != nil {
		return
	}

	defer f.Close()
	if _, err = f.Seek(int64(addr), 0); err != nil {
		return
	}
	b := []byte{0}
	if n, err := f.Read(b); err != nil || n == 0 {
		if err == nil {
			err = errors.New("Read zero bytes")
		}
		return 0, err
	}
	data = b[0]

	err = clrIoperm(addr)
	return
}

func setIoperm(addr uint16) (err error) {
	level := 3
	if _, _, errno := syscall.Syscall(sys_iopl,
		uintptr(level), 0, 0); errno != 0 {
		return err
	}
	num := 1
	on := 1
	if _, _, errno := syscall.Syscall(sys_ioperm, uintptr(addr),
		uintptr(num), uintptr(on)); errno != 0 {
		return err
	}

	return nil
}

func clrIoperm(addr uint16) (err error) {
	num := 1
	on := 0
	if _, _, errno := syscall.Syscall(sys_ioperm, uintptr(addr),
		uintptr(num), uintptr(on)); errno != 0 {
		return err
	}
	level := 0
	if _, _, errno := syscall.Syscall(sys_iopl,
		uintptr(level), 0, 0); errno != 0 {
		return err

	}
	return nil
}
