package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	for {
		time.Sleep(time.Second) // Пауза для снижения нагрузки на CPU
	}
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	// Инициализация HTTP сервера
	go func() {
		http.HandleFunc("/", executeHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Запуск службы
	if err := runService(); err != nil {
		log.Fatal(err)
	}
}

func runService() error {
	svcConfig := &service.Config{
		Name:        "Монитор запуска",
		DisplayName: "Монитор запуска",
		Description: "Служба для отслеживания команд в командной строке и их выполнение",
	}

	prg := &program{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return fmt.Errorf("ошибка создания службы: %v", err)
	}

	if err := svc.Run(); err != nil {
		return fmt.Errorf("ошибка запуска службы: %v", err)
	}
	return nil
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем команду из URL
	command := r.URL.Path[len("/"):]
	switch command {
	case "":
		http.Error(w, "Не указана команда", http.StatusBadRequest)
		return
	case "calculator":
		fmt.Fprintf(w, "Информация о калькуляторе\n")
	case "systeminfo":
		fmt.Fprintf(w, "Информация о системе\n")
	default:
		// Для неизвестных команд можно добавить обработку или вернуть ошибку
		http.Error(w, "Неизвестная команда", http.StatusBadRequest)
		return
	}

	// // Выполнение команды в командной строке
	// cmd := exec.Command("cmd", "/c", command)
	// output, err := cmd.Output()
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Ошибка выполнения команды: %s", err), http.StatusInternalServerError)
	// 	return
	// }

	// // Возвращаем результат выполнения команды
	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// w.Write(output)
}

// func onExit() {
// 	// Выполняется при выходе из приложения
// 	systray.Quit()
// }

// func onReady() {
// 	iconPath := "startup-monitor.ico"
// 	iconBytes, err := os.ReadFile(iconPath)
// 	if err != nil {
// 		fmt.Println("Ошибка чтения файла иконки:", err)
// 		return
// 	}

// 	systray.SetIcon(iconBytes)

// 	mOpenCmd := systray.AddMenuItem("Показать командную строку", "Открыть командную строку")
// 	systray.AddSeparator()
// 	mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")

// 	// Функция для обработки нажатий на элементы меню
// 	go func() {
// 		for {
// 			select {
// 			case <-mQuit.ClickedCh:
// 				systray.Quit()
// 				return
// 			case <-mOpenCmd.ClickedCh:
// 				openCmd.OpenCmd()
// 			}
// 		}
// 	}()

// 	// Устанавливаем всплывающее сообщение при наведении на иконку
// 	systray.SetTooltip("health-checker")
// }

// func onExit() {
// 	// Выполняется при выходе из приложения
// }

// func handleCommand(input string) {
// 	parts := strings.Fields(input)
// 	if len(parts) == 0 {
// 		return
// 	}

// 	switch parts[0] {
// 	case "info":
// 		if len(parts) < 2 {
// 			fmt.Println("Введите: info <application>")
// 			return
// 		}
// 		application := parts[1]
// 		cmdCommand(application)
// 	case "exit":
// 		os.Exit(1)
// 	default:
// 		fmt.Println("Команда не существует:", input)
// 	}
// }

// func cmdCommand(application string) {
// 	cmd := exec.Command("cmd", "-c", "ps aux | grep "+application)
// 	output, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println("Ошибка:", err)
// 		return
// 	}
// 	fmt.Println("Информация о приложении:")
// 	fmt.Println(string(output))
// }
