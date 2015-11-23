package utils

import "sync"
import "os"
import "syscall"

type FileMutex struct {
	mutex *sync.Mutex
	file *os.File
}

func MakeFileMutex(filename string) *FileMutex {
	file, err := os.OpenFile(filename, os.O_RDONLY | os.O_CREATE, 0777)
	if err != nil {
		return &FileMutex{file: nil}
	}
	mutex := &sync.Mutex{}
	return &FileMutex{file: file, mutex: mutex}
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

