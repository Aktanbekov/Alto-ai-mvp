package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"altoai_mvp/internal/models"
)

var ErrNotFound = errors.New("not found")

type UserRepo interface {
	List() ([]models.User, error)
	Get(id string) (models.User, error)
	Create(email, name string) (models.User, error)
	Update(id string, email, name *string) (models.User, error)
	Delete(id string) error
	Close() error
}

type userMemoryRepo struct {
	mu    sync.RWMutex
	store map[string]models.User
}

func NewUserMemoryRepo() UserRepo {
	return &userMemoryRepo{store: map[string]models.User{}}
}

func (r *userMemoryRepo) List() ([]models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.User, 0, len(r.store))
	for _, u := range r.store {
		out = append(out, u)
	}
	return out, nil
}

func (r *userMemoryRepo) Get(id string) (models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.store[id]
	if !ok {
		return models.User{}, ErrNotFound
	}
	return u, nil
}

func (r *userMemoryRepo) Create(email, name string) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	u := models.User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.store[u.ID] = u
	return u, nil
}

func (r *userMemoryRepo) Update(id string, email, name *string) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.store[id]
	if !ok {
		return models.User{}, ErrNotFound
	}
	if email != nil {
		u.Email = *email
	}
	if name != nil {
		u.Name = *name
	}
	u.UpdatedAt = time.Now().UTC()
	r.store[id] = u
	return u, nil
}

func (r *userMemoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}

func (r *userMemoryRepo) Close() error {
	// Nothing to close for in-memory repository
	return nil
}
