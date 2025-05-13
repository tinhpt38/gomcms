// Script to fix indentation in the file
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	filePath := "/Users/tinhp/Workspace/ITC/PROJECTS/gomcms/server/service/checkins/attendance_check_in.go"

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Fix indentation issues
	fixedContent := bytes.ReplaceAll(content,
		[]byte("resultList := checkConditionsParallel(gc.agp, gc.conditions, req, ip)				// Xử lý kết quả"),
		[]byte("resultList := checkConditionsParallel(gc.agp, gc.conditions, req, ip) // Xử lý kết quả"))

	fixedContent = bytes.ReplaceAll(fixedContent,
		[]byte("					}					if result.Pass {"),
		[]byte("					}\n\n					if result.Pass {"))

	if err := ioutil.WriteFile(filePath, fixedContent, 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Println("File fixed successfully!")
}
