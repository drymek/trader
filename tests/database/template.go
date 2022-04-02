package database

import (
	"fmt"
	"io"
	"os"
)

func CreateFromTemplate(src string, dst string) error {
	source, err := os.Open(fmt.Sprintf("../../../../../%s", src))
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(fmt.Sprintf("../../../../../%s", dst))
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)

	return err
}
