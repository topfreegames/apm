package utils

import "sync"
import "os"
import "syscall"
import "fmt"

type FileMutex struct {
	mutex sync.Mutex
	file *os.File
}

func MakeFileMutex(filename string) *FileMutex {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	return &FileMutex{file: file}
}

func (fMutex *FileMutex) Lock() {
	fMutex.mutex.Lock()
	if fMutex.file != nil {
		if err := syscall.Flock(int(fMutex.file.Fd()), syscall.LOCK_EX); err != nil {
			panic(err)
		}
	}
}

func (fMutex *FileMutex) Unlock() {
	fMutex.mutex.Unlock()
	if fMutex.file != nil {
		if err := syscall.Flock(int(fMutex.file.Fd()), syscall.LOCK_UN); err != nil {
			panic(err)
		}
	}
}

