package file_handler

import (
	"io"
	"os"
)

func Read(file_name string) ([]byte, error) {
	//Open file
	file, err := os.Open(file_name)

	if err != nil {
		if err.Error() == "open "+file_name+": The system cannot find the file specified." {
			return []byte{}, nil
		}
		return nil, err
	}

	defer file.Close()

	//Read contents of file
	data, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func Save(file_name string, data []byte) error {
	//Open file
	file, err := os.OpenFile(file_name, os.O_CREATE, 0600)

	if err != nil {
		return err
	}

	defer file.Close()

	//Clear contents of file
	file.Truncate(0)

	//Write data to file
	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
