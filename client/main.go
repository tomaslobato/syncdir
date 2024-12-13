package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type FileItem struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	IsDir   bool   `json:"isDir"`
	Content string `json:"content"`
}

func main() {
	godotenv.Load(".env")

	url := os.Getenv("URL")
	dir := os.Getenv("FOLDER")

	files, err := RecursiveGetFiles("", dir)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(files)

	filesData, err := json.Marshal(files)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(filesData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Status code: %d\n", res.StatusCode)
	}
}

func RecursiveGetFiles(currentDir string, ROOT string) ([]FileItem, error) {
	var files []FileItem
	fls, err := os.ReadDir(ROOT + "/" + currentDir)
	if err != nil {
		return nil, err
	}

	for _, f := range fls {
		var file FileItem
		file.Name = f.Name()
		file.IsDir = f.IsDir()
		file.Id = filepath.Join(currentDir, file.Name)
		if file.IsDir {
			files = append(files, file)
			subFiles, err := RecursiveGetFiles(file.Id, ROOT)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, file)
		}
	}
	return files, nil
}
