package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.Run()
	return nil
}

func (p *program) Run() {
	// Запускаем HTTP сервер
	go func() {
		http.HandleFunc("/", executeHandler)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// Бесконечный цикл для предотвращения завершения работы службы
	for {
		time.Sleep(time.Second)
	}
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	// Проверяем, запущена ли программа как служба
	exists, err := isServiceInstalled("Startup-monitor")
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		// Если программа не запущена как служба, устанавливаем её
		if err := installService(); err != nil {
			log.Fatal(err)
		}

		// Проверяем текущую операционную систему
		switch os := runtime.GOOS; os {
		case "windows":
			// В Windows используем команду sc
			cmd := exec.Command("sc", "start", "Startup-monitor")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		case "linux":
			// В Linux используем соответствующую команду для запуска службы
			cmd := exec.Command("systemctl", "start", "startup-monitor.service")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}
		return
	}

	// Запуск службы
	if err := runService(); err != nil {
		log.Fatal(err)
	}
}

func isServiceInstalled(serviceName string) (bool, error) {
	svcConfig := &service.Config{
		Name: serviceName,
	}

	prg := &program{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return false, fmt.Errorf("ошибка создания службы: %v", err)
	}

	status, err := svc.Status()
	if err != nil {
		if status == service.StatusUnknown {
			return false, nil // Служба не установлена
		}
		return false, fmt.Errorf("ошибка проверки статуса службы: %v", err)
	}

	return true, nil // Служба установлена
}

func runService() error {
	svcConfig := &service.Config{
		Name:        "Startup-monitor",
		DisplayName: "Монитор запуска",
		Description: "Служба для отслеживания статистики какого-либо приложения через командную строку",
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

func installService() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	svcConfig := &service.Config{
		Name:        "Startup-monitor",
		DisplayName: "Монитор запуска",
		Description: "Служба для отслеживания статистики какого-либо приложения через командную строку",
		Executable:  exePath,
	}

	prg := &program{}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return fmt.Errorf("ошибка создания службы: %v", err)
	}

	if err := svc.Install(); err != nil {
		return fmt.Errorf("ошибка установки службы: %v", err)
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
	case "help":
		fmt.Fprintf(w, "Команды:\n 1) info - информация о приложении;\n 2) calculator - проверочная команда;\n 3) systeminfo - проверочная команда.\n")
	case "info":
		fmt.Fprintf(w, "Данное приложение позволяет собирать и получать статистику разных процессов и приложений\n")
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
