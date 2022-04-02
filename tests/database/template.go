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

	x, err := io.Copy(destination, source)
	_ = x
	y, err := source.Stat()
	_ = y
	z, err := destination.Stat()
	_ = z

	return err
}
