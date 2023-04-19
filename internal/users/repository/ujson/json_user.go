package ujson

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"strconv"
	"sync"
	"time"

	"refactoring/internal/domain"
)

type UserStore struct {
	Increment int             `json:"increment"`
	List      domain.UserList `json:"list"`
}

var (
	ErrReadStore     = errors.New("error with read store")
	ErrWithParseUser = errors.New("error with parse user")
	ErrWriteStore    = errors.New("error with write store")
	ErrUserNotFound  = errors.New("error with found user")
)

type userRepository struct {
	store *os.File
	mu    *sync.RWMutex
}

func NewUserRepository(store *os.File) domain.UserRepository {
	return &userRepository{
		store: store,
		mu:    &sync.RWMutex{},
	}
}

func (u *userRepository) Create(ctx context.Context, displayName, email string) (domain.User, error) {
	data, err := os.ReadFile(u.store.Name())
	if err != nil {
		return domain.User{}, ErrReadStore
	}

	s := UserStore{}
	err = json.Unmarshal(data, &s)
	if err != nil {

		return domain.User{}, ErrWithParseUser
	}

	s.Increment++
	id := strconv.Itoa(s.Increment)
	user := domain.User{
		ID:          id,
		CreatedAt:   time.Now(),
		DisplayName: displayName,
		Email:       email,
	}

	u.mu.Lock()
	s.List[id] = user
	u.mu.Unlock()

	b, err := json.Marshal(&s)
	if err != nil {
		return domain.User{}, ErrWithParseUser
	}

	err = os.WriteFile(u.store.Name(), b, fs.ModePerm)
	if err != nil {
		return domain.User{}, ErrWriteStore
	}

	return user, nil
}

func (u *userRepository) GetByID(ctx context.Context, userID string) (domain.User, error) {
	data, err := os.ReadFile(u.store.Name())
	if err != nil {
		return domain.User{}, ErrReadStore
	}

	s := UserStore{}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return domain.User{}, ErrWithParseUser
	}

	u.mu.RLock()
	val, ok := s.List[userID]
	u.mu.RUnlock()

	if !ok {
		return domain.User{}, ErrUserNotFound
	}
	return val, nil
}

func (u *userRepository) GetAll(ctx context.Context) (domain.UserList, error) {
	data, err := os.ReadFile(u.store.Name())
	if err != nil {
		return nil, ErrReadStore
	}

	s := UserStore{}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, ErrWithParseUser
	}

	return s.List, nil
}

func (u *userRepository) Update(ctx context.Context, userID, displayName string) error {
	data, err := os.ReadFile(u.store.Name())
	if err != nil {
		return ErrReadStore
	}

	s := UserStore{}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return ErrWithParseUser
	}
	u.mu.RLock()
	val, ok := s.List[userID]
	u.mu.RUnlock()

	if !ok {
		return ErrUserNotFound
	}

	val.DisplayName = displayName

	u.mu.Lock()
	s.List[userID] = val
	u.mu.Unlock()

	b, err := json.Marshal(&s)
	if err != nil {
		return ErrWithParseUser
	}

	err = os.WriteFile(u.store.Name(), b, fs.ModePerm)
	if err != nil {
		return ErrWriteStore
	}

	return nil
}
func (u *userRepository) Delete(ctx context.Context, userID string) error {
	data, err := os.ReadFile(u.store.Name())
	if err != nil {
		return ErrReadStore
	}

	s := UserStore{}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return ErrWithParseUser
	}

	if _, ok := s.List[userID]; !ok {
		return ErrUserNotFound
	}

	u.mu.Lock()
	delete(s.List, userID)
	u.mu.Unlock()

	b, err := json.Marshal(&s)
	if err != nil {
		return ErrWithParseUser
	}

	err = os.WriteFile(u.store.Name(), b, fs.ModePerm)
	if err != nil {
		return ErrWriteStore
	}

	return nil
}
