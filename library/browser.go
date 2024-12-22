package library

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		err = fmt.Errorf("不支援的作業系統")
	}
	if err != nil {
		log.Fatalf("無法啟動瀏覽器: %v", err)
	}
}
