package service

import (
	"fmt"

	"mediasec/entities"
	"mediasec/repository"
)

type ContentVerificationService struct {
	sources      *repository.MediaSourceRepository
	publications *repository.PublicationRepository
	incidents    *repository.IncidentRepository
}

func NewContentVerificationService(
	s *repository.MediaSourceRepository,
	p *repository.PublicationRepository,
	i *repository.IncidentRepository,
) *ContentVerificationService {
	return &ContentVerificationService{sources: s, publications: p, incidents: i}
}

func (svc *ContentVerificationService) CheckPublication(pubID int) (bool, string) {
	pub, ok := svc.publications.FindByID(pubID)
	if !ok {
		return false, "публикация не найдена"
	}

	if pub.Status != "pending" {
		return false, fmt.Sprintf("публикация уже обработана: статус — %s", pub.Status)
	}

	src, ok := svc.sources.FindByID(pub.SourceID)
	if !ok {
		svc.publications.UpdateStatus(pubID, "rejected")
		svc.incidents.Create("disinformation", "Источник не найден в базе", "critical", pubID)
		return false, "источник не найден. Публикация отклонена. Создан инцидент [critical]"
	}

	if src.IsBlocked {
		svc.publications.UpdateStatus(pubID, "rejected")
		svc.incidents.Create("disinformation", fmt.Sprintf("Публикация из заблокированного источника: %s", src.FullName), "high", pubID)
		return false, "источник заблокирован. Публикация отклонена. Создан инцидент [high]"
	}

	if src.TrustLevel == "low" {
		svc.publications.UpdateStatus(pubID, "rejected")
		svc.incidents.Create("disinformation", fmt.Sprintf("Источник с низким уровнем доверия: %s", src.FullName), "medium", pubID)
		return false, "источник ненадёжен. Публикация отклонена. Создан инцидент [medium]"
	}

	svc.publications.UpdateStatus(pubID, "approved")
	return true, "публикация проверена и одобрена"
}

type IncidentService struct {
	repo *repository.IncidentRepository
}

func NewIncidentService(r *repository.IncidentRepository) *IncidentService {
	return &IncidentService{repo: r}
}

func (svc *IncidentService) CreateIncident(category, description, severity string, publicationID int) entities.Incident {
	return svc.repo.Create(category, description, severity, publicationID)
}

func (svc *IncidentService) GetAllIncidents() []entities.Incident {
	return svc.repo.GetAll()
}

func (svc *IncidentService) UpdateSeverity(id int, severity string) bool {
	return svc.repo.UpdateSeverity(id, severity)
}
