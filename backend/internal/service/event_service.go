package service

import (
	"time"

	"github.com/ryohighbridge/learn-github-copilot/backend/internal/domain"
)

type EventService struct {
	repo EventRepositoryInterface
}

type EventRepositoryInterface interface {
	GetAll() ([]domain.Event, error)
	GetByID(id int) (*domain.Event, error)
	GetByDateRange(start, end time.Time) ([]domain.Event, error)
	Create(event *domain.Event) error
	Update(event *domain.Event) error
	Delete(id int) error
}

func NewEventService(repo EventRepositoryInterface) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetAllEvents() ([]domain.Event, error) {
	return s.repo.GetAll()
}

func (s *EventService) GetEventByID(id int) (*domain.Event, error) {
	return s.repo.GetByID(id)
}

func (s *EventService) GetEventsByDateRange(start, end time.Time) ([]domain.Event, error) {
	return s.repo.GetByDateRange(start, end)
}

func (s *EventService) CreateEvent(event *domain.Event) error {
	// バリデーション
	if event.Title == "" {
		return domain.ErrInvalidInput
	}
	if event.EndDate.Before(event.StartDate) {
		return domain.ErrInvalidInput
	}

	return s.repo.Create(event)
}

func (s *EventService) UpdateEvent(event *domain.Event) error {
	// バリデーション
	if event.Title == "" {
		return domain.ErrInvalidInput
	}
	if event.EndDate.Before(event.StartDate) {
		return domain.ErrInvalidInput
	}

	existing, err := s.repo.GetByID(event.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrNotFound
	}

	return s.repo.Update(event)
}

func (s *EventService) DeleteEvent(id int) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrNotFound
	}

	return s.repo.Delete(id)
}
