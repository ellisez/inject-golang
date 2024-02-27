package utils

import (
	"errors"
	"fmt"
	"os"
)

func FileExists(filename string) (bool, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if stat.IsDir() {
		return false, errors.New(fmt.Sprintf("%s is not File!", filename))
	}
	return true, nil
}

func CreateFileIfNotExists(filename string) error {
	exists, err := FileExists(filename)
	if err != nil {
		return err
	}
	if !exists {
		_, err = os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func DirectoryExists(dirname string) (bool, error) {
	stat, err := os.Stat(dirname)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if !stat.IsDir() {
		return false, errors.New(fmt.Sprintf("%s is not Directory!", dirname))
	}
	return true, nil
}

func CreateDirectoryIfNotExists(dirname string) error {
	exists, err := DirectoryExists(dirname)
	if err != nil {
		return err
	}
	if !exists {
		err = os.Mkdir(dirname, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
