package utils

import (
	"fmt"
	"time"
)

func SimulateLogFile(filePath, fileType string) {
	// Create a new file util.
	fileUtil := NewFileUtil(filePath, fileType)
	// Create a new file.
	err := fileUtil.CreateFile()
	if err != nil {
		panic(err)
	}
	// Write some log lines to the file.
	for i := 0; i < 1000; i++ {
		fileUtil.WriteToFile(fmt.Sprintf("Log line %d \n", i))
		time.Sleep(500 * time.Millisecond)
	}
	fileUtil.CloseFile()
}
