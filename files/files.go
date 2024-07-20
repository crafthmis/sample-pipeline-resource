package files

import "os"

func WriteFile(filepath string, data string) (filePath string, err error) {
	filePath = filepath
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return
	}
	return
}

func DeleteFile(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	return nil
}
