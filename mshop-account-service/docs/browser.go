package docs

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

func LaunchSwagger(port string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("SWAGGER_HOST")
	if host == "" {
		log.Fatal("Missing SWAGGER_HOST in .env")
	}

	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")

	url := fmt.Sprintf("http://%s:%s/swagger/index.html", host, port)
	fmt.Printf("Swagger UI available at: %s\n", url)

	go openBrowser(url)
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		cmd = "cmd.exe"
		args = []string{"/c", "start", url}
	}

	exec.Command(cmd, args...).Start()
}
