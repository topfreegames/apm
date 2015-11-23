package utils

import "io/ioutil"
import "os"

import "github.com/BurntSushi/toml"

func WriteFile(filepath string, b []byte) error {
	return ioutil.WriteFile(filepath, b, 0660)
}

func GetFile(filepath string) (*os.File, error) {
	if _, err := os.Stat(filepath); err == nil {
		return os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	}
	return os.Create(filepath)
}

func SafeReadTomlFile(filename string, v interface{}) error {
	fileLock := MakeFileMutex(filename)
	fileLock.Lock()
	defer fileLock.Unlock()
	_, err := toml.Decode(filename, v)
	return err
}

func SafeWriteTomlFile(v interface{}, filename string) error {
	fileLock := MakeFileMutex(filename)
	fileLock.Lock()
	defer fileLock.Unlock()
	f, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		return err
	}	
	encoder := toml.NewEncoder(f)
	return encoder.Encode(v)
}

func DeleteFile(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return true
	}
	err = os.Remove(filepath)
	return err == nil
}
