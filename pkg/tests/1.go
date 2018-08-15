package main

import (
	"os"
	"syscall"
	"unsafe"
	"fmt"
	"time"
)

func Open(regionName string, flags int, perm os.FileMode) (*os.File, error) {
	name, err := syscall.BytePtrFromString(regionName)
	if err != nil {
		return nil, err
	}
	fd, _, errno := syscall.Syscall(syscall.SYS_SHM_OPEN,
		uintptr(unsafe.Pointer(name)),
		uintptr(flags), uintptr(perm),
	)
	if errno != 0 {
		return nil, errno
	}
	return os.NewFile(fd, regionName), nil
}

type stream struct {
}

func newstream() interface{} {
	return new(stream)
}

func resetstream(i interface{}) {
	var null stream
	s := i.(*stream)
	*s = null
}

func Unlink(regionName string) error {
	name, err := syscall.BytePtrFromString(regionName)
	if err != nil {
		return err
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_SHM_UNLINK,
		uintptr(unsafe.Pointer(name)), 0, 0,
	); errno != 0 {
		return errno
	}
	return nil
}


func main() {
	file, err := Open("my_region", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	if err := syscall.Ftruncate(
		int(file.Fd()), 50,
	); err != nil {
	   fmt.Println(err)
	}

	b, err := syscall.Mmap(int(file.Fd()), 0, 32, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	// syscall.Ftruncate if new, etc
	_, err = file.Write([]byte("hello"))
	if err != nil {
		fmt.Println(err)
	}
	b := make([]byte, 10)
	_, err = file.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	time.Sleep(10 * time.Second)
	defer file.Close()
	defer Unlink(file.Name())


}
