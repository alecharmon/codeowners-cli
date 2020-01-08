package core

import "os"

func OpenFile(path string) (bool, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
