package docs

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func LaunchSwagger(port string) {
	host := os.Getenv("SWAGGER_HOST")
	if host == "" {
		host = os.Getenv("HOST")
	}
	if host == "" {
		host = "localhost"
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
		cmd = "xdg-open"
		args = []string{url}
	}

	exec.Command(cmd, args...).Start()
}
