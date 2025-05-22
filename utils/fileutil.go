package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type FileUtil struct {
	// filePath is the path of the file.
	filePath string
	// fileType is the type of the file.
	fileType string
	// file is the file pointer.
	file *os.File
}

// NewFileUtil creates a new file util.
func NewFileUtil(filePath string, fileType string) *FileUtil {
	return &FileUtil{
		filePath: filePath,
		fileType: fileType,
	}
}

// CreateFile creates a new file with the given name and type.
// It returns an error if the file already exists or cannot be created.
func (f *FileUtil) CreateFile() error {

	// Check if the file already exists.
	if _, err := os.Stat(f.filePath); !os.IsNotExist(err) {
		os.Remove(f.filePath)
	}

	// Create the file.
	file, err := os.Create(f.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err.Error())
	}

	f.file = file
	return nil
}

// OpenFile opens a file and returns the file pointer.
// It returns an error if the file does not exist or cannot be opened.
func (f *FileUtil) OpenFile() error {
	// Check if the file exists.
	if _, err := os.Stat(f.filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", err.Error())
	}

	// Open the file.
	file, err := os.Open(f.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err.Error())
	}

	f.file = file
	return nil
}

// TailFile tails the file and sends each new log line to the provided channel.
func (f *FileUtil) TailFile(ctx context.Context, logChan chan<- string) {
	// Seek to the end of the file to only get new lines.
	_, err := f.file.Seek(0, io.SeekEnd)
	if err != nil {
		logChan <- fmt.Sprintf("failed to seek file: %s", err.Error())
		close(logChan)
		return
	}
	reader := bufio.NewReader(f.file)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(logChan)
			return
		case <-ticker.C:
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				logChan <- line
			}
		}
	}
}

func (f *FileUtil) CloseFile() error {
	err := f.file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %s", err.Error())
	}

	return nil
}

func (f *FileUtil) WriteToFile(data string) error {
	if _, err := f.file.WriteString(data); err != nil {
		return fmt.Errorf("failed to write to file: %s", err.Error())
	}

	return nil
}
