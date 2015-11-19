package utils

import "io/ioutil"
import "os"

func WriteFile(filepath string, b []byte) error {
	return ioutil.WriteFile(filepath, b, 0660)
}

func GetFile(filepath string) (*os.File, error) {
	if _, err := os.Stat(filepath); err == nil {
		return os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	}
	return os.Create(filepath)
}

func DeleteFile(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return true
	}
	err = os.Remove(filepath)
	return err == nil
}
