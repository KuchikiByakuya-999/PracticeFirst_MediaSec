package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"mediasec/repository"
	"mediasec/service"
)

type ConsoleUI struct {
	sourceRepo  *repository.MediaSourceRepository
	pubRepo     *repository.PublicationRepository
	verifySvc   *service.ContentVerificationService
	incidentSvc *service.IncidentService
	reader      *bufio.Reader
}

func NewConsoleUI(
	sr *repository.MediaSourceRepository,
	pr *repository.PublicationRepository,
	vs *service.ContentVerificationService,
	is *service.IncidentService,
) *ConsoleUI {
	return &ConsoleUI{
		sourceRepo:  sr,
		pubRepo:     pr,
		verifySvc:   vs,
		incidentSvc: is,
		reader:      bufio.NewReader(os.Stdin),
	}
}

func (ui *ConsoleUI) Run() {
	choice := -1
	for choice != 0 {
		ui.printMenu()
		var err error
		choice, err = ui.readInt()
		if err != nil {
			fmt.Println("Некорректный ввод. Введите номер пункта меню.")
			continue
		}

		switch choice {
		case 1:
			ui.showSources()
		case 2:
			ui.createSource()
		case 3:
			ui.showPublications()
		case 4:
			ui.createPublication()
		case 5:
			ui.checkPublication()
		case 6:
			ui.showIncidents()
		case 7:
			ui.createIncidentManual()
		case 8:
			ui.updateIncidentSeverity()
		case 0:
			fmt.Println("Выход из программы.")
		default:
			fmt.Println("Некорректный пункт меню.")
		}
	}
}

func (ui *ConsoleUI) printMenu() {
	fmt.Println("1. Показать все источники")
	fmt.Println("2. Создать источник")
	fmt.Println("3. Показать все публикации")
	fmt.Println("4. Добавить публикацию")
	fmt.Println("5. Верифицировать публикацию")
	fmt.Println("6. Показать все инциденты")
	fmt.Println("7. Создать инцидент вручную")
	fmt.Println("8. Изменить серьёзность инцидента")
	fmt.Println("0. Выход")
	fmt.Print("Выберите действие: ")
}

func (ui *ConsoleUI) showSources() {
	sources := ui.sourceRepo.GetAll()
	if len(sources) == 0 {
		fmt.Println("Источники отсутствуют.")
		return
	}
	for _, s := range sources {
		blocked := "нет"
		if s.IsBlocked {
			blocked = "да"
		}
		fmt.Printf("ID: %d, название: %s, URL: %s, доверие: %s, заблокирован: %s\n",
			s.ID, s.FullName, s.URL, s.TrustLevel, blocked)
	}
}

func (ui *ConsoleUI) createSource() {
	fmt.Print("Введите название источника: ")
	name := ui.readLine()
	fmt.Print("Введите URL: ")
	url := ui.readLine()
	fmt.Print("Введите уровень доверия (high/medium/low): ")
	trust := ui.readLine()

	if name == "" || url == "" || trust == "" {
		fmt.Println("Название, URL и уровень доверия не должны быть пустыми.")
		return
	}

	s := ui.sourceRepo.Create(name, url, trust)
	fmt.Printf("Источник создан. ID: %d\n", s.ID)
}

func (ui *ConsoleUI) showPublications() {
	pubs := ui.pubRepo.GetAll()
	if len(pubs) == 0 {
		fmt.Println("Публикации отсутствуют.")
		return
	}
	for _, p := range pubs {
		fmt.Printf("ID: %d, заголовок: %s, URL: %s, источник ID: %d, статус: %s\n",
			p.ID, p.Title, p.URL, p.SourceID, p.Status)
	}
}

func (ui *ConsoleUI) createPublication() {
	fmt.Print("Введите заголовок: ")
	title := ui.readLine()
	fmt.Print("Введите URL материала: ")
	url := ui.readLine()
	fmt.Print("Введите ID источника: ")
	sourceID, err := ui.readInt()
	if err != nil || sourceID <= 0 {
		fmt.Println("Некорректный ID источника.")
		return
	}

	if title == "" || url == "" {
		fmt.Println("Заголовок и URL не должны быть пустыми.")
		return
	}

	if _, ok := ui.sourceRepo.FindByID(sourceID); !ok {
		fmt.Println("Источник с таким ID не найден.")
		return
	}

	p := ui.pubRepo.Create(title, url, sourceID)
	fmt.Printf("Публикация добавлена. ID: %d\n", p.ID)
}

func (ui *ConsoleUI) checkPublication() {
	fmt.Print("Введите ID публикации: ")
	id, err := ui.readInt()
	if err != nil || id <= 0 {
		fmt.Println("Некорректный ID публикации.")
		return
	}
	_, msg := ui.verifySvc.CheckPublication(id)
	fmt.Println(msg)
}

func (ui *ConsoleUI) showIncidents() {
	incidents := ui.incidentSvc.GetAllIncidents()
	if len(incidents) == 0 {
		fmt.Println("Инциденты отсутствуют.")
		return
	}
	for _, i := range incidents {
		pubID := "не указана"
		if i.PublicationID != 0 {
			pubID = strconv.Itoa(i.PublicationID)
		}
		fmt.Printf("ID: %d, категория: %s, серьёзность: %s, публикация ID: %s, описание: %s\n",
			i.ID, i.Category, i.Severity, pubID, i.Description)
	}
}

func (ui *ConsoleUI) createIncidentManual() {
	fmt.Print("Введите категорию (disinformation/copyright/data_leak): ")
	category := ui.readLine()
	fmt.Print("Введите описание: ")
	description := ui.readLine()
	fmt.Print("Введите серьёзность (low/medium/high/critical): ")
	severity := ui.readLine()
	fmt.Print("Введите ID публикации или 0, если не применимо: ")
	pubID, err := ui.readInt()
	if err != nil || pubID < 0 {
		fmt.Println("Некорректный ID публикации.")
		return
	}

	if category == "" || description == "" || severity == "" {
		fmt.Println("Категория, описание и серьёзность не должны быть пустыми.")
		return
	}

	if pubID != 0 {
		if _, ok := ui.pubRepo.FindByID(pubID); !ok {
			fmt.Println("Публикация с таким ID не найдена.")
			return
		}
	}

	inc := ui.incidentSvc.CreateIncident(category, description, severity, pubID)
	fmt.Printf("Инцидент создан. ID: %d\n", inc.ID)
}

func (ui *ConsoleUI) updateIncidentSeverity() {
	fmt.Print("Введите ID инцидента: ")
	id, err := ui.readInt()
	if err != nil || id <= 0 {
		fmt.Println("Некорректный ID инцидента.")
		return
	}
	fmt.Print("Введите новую серьёзность (low/medium/high/critical): ")
	severity := ui.readLine()
	if severity == "" {
		fmt.Println("Серьёзность не должна быть пустой.")
		return
	}
	if ui.incidentSvc.UpdateSeverity(id, severity) {
		fmt.Println("Серьёзность инцидента обновлена.")
	} else {
		fmt.Println("Инцидент не найден.")
	}
}

func (ui *ConsoleUI) readInt() (int, error) {
	line := ui.readLine()
	return strconv.Atoi(strings.TrimSpace(line))
}

func (ui *ConsoleUI) readLine() string {
	line, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(line)
}
