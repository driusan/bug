package bugs

import (
	"fmt"
	"os"
)

func (b *Bug) Read(p []byte) (int, error) {
	if b.descFile == nil {
		dir := b.GetDirectory()
		fp, err := os.OpenFile(string(dir)+"/Description", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		b.descFile = fp
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %s", err.Error())
			return 0, NoDescriptionError
		}
	}

	return b.descFile.Read(p)
}

func (b *Bug) Write(data []byte) (n int, err error) {
	if b.descFile == nil {
		dir := b.GetDirectory()
		os.MkdirAll(string(dir), 0755)
		fp, err := os.OpenFile(string(dir)+"/Description", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to bug: %s", err.Error())
			return 0, err
		}
		b.descFile = fp
	}
	return b.descFile.Write(data)
}

func (b *Bug) WriteAt(data []byte, off int64) (n int, err error) {
	if b.descFile == nil {
		dir := b.GetDirectory()
		os.MkdirAll(string(dir), 0755)
		fp, err := os.OpenFile(string(dir)+"/Description", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to bug: %s", err.Error())
			return 0, err
		}
		b.descFile = fp
	}
	return b.descFile.WriteAt(data, off)
}
func (b Bug) Close() error {
	if b.descFile != nil {
		err := b.descFile.Close()
		b.descFile = nil
		return err
	}
	return nil
}
