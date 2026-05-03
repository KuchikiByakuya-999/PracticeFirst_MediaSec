package repository

import (
	"errors"

	"mediasec/entities"
)

type MediaSourceRepository struct {
	data   []entities.MediaSource
	nextID int
}

func NewMediaSourceRepository() *MediaSourceRepository {
	return &MediaSourceRepository{
		nextID: 3,
		data: []entities.MediaSource{
			{ID: 1, FullName: "РБК", URL: "rbc.ru", TrustLevel: "high", IsBlocked: false},
			{ID: 2, FullName: "Неизвестный блог", URL: "fakeblog.ru", TrustLevel: "low", IsBlocked: false},
		},
	}
}

func (r *MediaSourceRepository) GetAll() []entities.MediaSource {
	return r.data
}

func (r *MediaSourceRepository) FindByID(id int) (entities.MediaSource, bool) {
	for _, s := range r.data {
		if s.ID == id {
			return s, true
		}
	}
	return entities.MediaSource{}, false
}

func (r *MediaSourceRepository) Create(fullName, url, trustLevel string) entities.MediaSource {
	s := entities.MediaSource{
		ID:         r.nextID,
		FullName:   fullName,
		URL:        url,
		TrustLevel: trustLevel,
		IsBlocked:  false,
	}
	r.data = append(r.data, s)
	r.nextID++
	return s
}

func (r *MediaSourceRepository) Block(id int) bool {
	for i := range r.data {
		if r.data[i].ID == id {
			r.data[i].IsBlocked = true
			return true
		}
	}
	return false
}

type PublicationRepository struct {
	data   []entities.Publication
	nextID int
}

func NewPublicationRepository() *PublicationRepository {
	return &PublicationRepository{nextID: 1}
}

func (r *PublicationRepository) GetAll() []entities.Publication {
	return r.data
}

func (r *PublicationRepository) FindByID(id int) (entities.Publication, bool) {
	for _, p := range r.data {
		if p.ID == id {
			return p, true
		}
	}
	return entities.Publication{}, false
}

func (r *PublicationRepository) Create(title, url string, sourceID int) entities.Publication {
	p := entities.Publication{
		ID:       r.nextID,
		Title:    title,
		URL:      url,
		SourceID: sourceID,
		Status:   "pending",
	}
	r.data = append(r.data, p)
	r.nextID++
	return p
}

func (r *PublicationRepository) UpdateStatus(id int, status string) error {
	for i := range r.data {
		if r.data[i].ID == id {
			r.data[i].Status = status
			return nil
		}
	}
	return errors.New("публикация не найдена")
}

type IncidentRepository struct {
	data   []entities.Incident
	nextID int
}

func NewIncidentRepository() *IncidentRepository {
	return &IncidentRepository{nextID: 1}
}

func (r *IncidentRepository) GetAll() []entities.Incident {
	return r.data
}

func (r *IncidentRepository) FindByID(id int) (entities.Incident, bool) {
	for _, i := range r.data {
		if i.ID == id {
			return i, true
		}
	}
	return entities.Incident{}, false
}

func (r *IncidentRepository) Create(category, description, severity string, publicationID int) entities.Incident {
	inc := entities.Incident{
		ID:            r.nextID,
		Category:      category,
		Description:   description,
		Severity:      severity,
		PublicationID: publicationID,
	}
	r.data = append(r.data, inc)
	r.nextID++
	return inc
}

func (r *IncidentRepository) UpdateSeverity(id int, severity string) bool {
	for i := range r.data {
		if r.data[i].ID == id {
			r.data[i].Severity = severity
			return true
		}
	}
	return false
}
