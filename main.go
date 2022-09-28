package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/esirk/dokku_agent/models"
)

func main() {
	http.HandleFunc("/apps", getApps)
	http.HandleFunc("/apps/create", createApp)
	http.HandleFunc("/apps/destroy", destroyApp)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func getApps(w http.ResponseWriter, r *http.Request) {
	// Read command from request
	out, err := runCommand("dokku", "apps:list")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	apps := strings.Split(strings.TrimSpace(out), "\n")[1:]
	// Write response as JSON
	resp := make(map[string][]models.DokkuApp)
	for _, app := range apps {
		// Get app report
		out, err := runCommand("dokku", "apps:report", app)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		// Parse app report
		lines := strings.Split(strings.TrimSpace(out), "\n")
		appReport := make(map[string]string)
		for _, line := range lines {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				appReport[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
		timestamp, err := strconv.ParseInt(appReport["Created at"], 0, 64)
		
		if err != nil {
			fmt.Println("Error: ", err)
		}
		createdAt := time.Unix((timestamp), 0)
		// Create app object
		resp["apps"] = append(resp["apps"], models.DokkuApp{
			Name:      app,
			GitUrl:    "https://github.com/"+app,
			GitBranch: "master",
			CreatedAt: createdAt.Format(time.RFC1123),
			Status:    "running",
			Details: appReport,
		})
	}

	jData, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jData)
}

func createApp(w http.ResponseWriter, r *http.Request) {
	// Read app name from request
	appName := r.URL.Query().Get("name")
	out, err := runCommand("dokku", "apps:create", appName)
	if err != nil {
		fmt.Println("Error: ", out)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func destroyApp(w http.ResponseWriter, r *http.Request) {
	// Read app name from request
	appName := r.URL.Query().Get("name")
	out, err := runCommand("dokku", "apps:destroy", appName)
	if err != nil {
		fmt.Println("Error: ", out)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	error := bytes.NewBuffer(nil)
	cmd.Stderr = error

	stdout, err := cmd.Output()
	if err != nil {
		return error.String(), err
	}
	return string(stdout), nil
}
