package utils

import "io/ioutil"
import "os"

import "github.com/BurntSushi/toml"

func WriteFile(filepath string, b []byte) error {
	return ioutil.WriteFile(filepath, b, 0660)
}

func GetFile(filepath string) (*os.File, error) {
	return os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
}

func SafeReadTomlFile(filename string, v interface{}) error {
	fileLock := MakeFileMutex(filename)
	fileLock.Lock()
	defer fileLock.Unlock()
	_, err := toml.DecodeFile(filename, v)

	return err
}

func SafeWriteTomlFile(v interface{}, filename string) error {
	fileLock := MakeFileMutex(filename)
	fileLock.Lock()
	defer fileLock.Unlock()
	f, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0777)
	defer f.Close()
	if err != nil {
		return err
	}	
	encoder := toml.NewEncoder(f)
	return encoder.Encode(v)
}

func DeleteFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	err = os.Remove(filepath)
	return err
}
