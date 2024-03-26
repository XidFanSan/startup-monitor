package openCmd

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenCmd() {
	// Определяем команду для открытия командной строки в зависимости от ОС
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "cmd")
	case "darwin":
		cmd = exec.Command("open", "-a", "Terminal")
	case "linux":
		cmd = exec.Command("x-terminal-emulator")
	default:
		fmt.Println("Не поддерживаемая операционная система")
		return
	}

	// Запускаем службу
	if err := cmd.Start(); err != nil {
		fmt.Println("Ошибка при открытии командной строки:", err)
		return
	}
}
