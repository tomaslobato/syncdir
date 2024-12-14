package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type SyncRequest []FileItem

type FileItem struct {
	Id      string `json:"id"`
	IsDir   bool   `json:"isDir"`
	Content string `json:"content"`
}

func main() {
	godotenv.Load("../.env")

	http.HandleFunc("POST /sync", handleSync)

	fmt.Println("running port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func handleSync(w http.ResponseWriter, r *http.Request) {
	folder := os.Getenv("FOLDER")

	var files SyncRequest
	err := json.NewDecoder(r.Body).Decode(&files)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	for _, f := range files {
		path := filepath.Join(folder, f.Id)
		if f.IsDir {
			err := os.Mkdir(path, 0755)
			if err != nil {
				if os.IsNotExist(err) {
					os.Mkdir(folder, 0755)
					continue
				}
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("created dir at %s\n", path)
			continue
		}
		err := os.WriteFile(path, []byte(f.Content), 0644)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("created file at %s\n", path)
	}

	cmd := exec.Command("docker", "exec", "-u", "33", "-it", "nextcloud", "php", "occ", "files:scan", "--all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error executing command: ", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
