package repository

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/aaltgod/refactoring/internal/apperror"
	"github.com/aaltgod/refactoring/internal/user"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	store  string
	Logger *logrus.Logger
}

func NewRepository(store string, logger *logrus.Logger) *Repository {
	return &Repository{
		store:  store,
		Logger: logger,
	}
}

func (r *Repository) Get(id string) (user.User, error) {

	s := user.UserStore{}
	u := user.User{}

	f, err := ioutil.ReadFile(r.store)
	if err != nil {
		r.Logger.Warnln("ReadFile returns error")
		return u, apperror.ErrInternalServer
	}
	if err = json.Unmarshal(f, &s); err != nil {
		r.Logger.Warnln("Unmarshal returns error")
		return u, apperror.ErrInternalServer
	}

	u, ok := s.List[id]
	if !ok {
		r.Logger.Infof("User doesn't exist with id: %s", id)
		return u, apperror.ErrUserNotFound
	}

	return u, nil
}

func (r *Repository) Insert(u user.CreateUserRequest) (string, error) {

	f, err := ioutil.ReadFile(r.store)
	if err != nil {
		r.Logger.Warnln("ReadFile returns error")
		return "", apperror.ErrInternalServer
	}

	s := user.UserStore{}
	if err := json.Unmarshal(f, &s); err != nil {
		r.Logger.Warnln("Unmarshal returns error")
		return "", apperror.ErrInternalServer
	}

	s.Increment++
	createUser := user.User{
		CreatedAt:   time.Now(),
		DisplayName: u.DisplayName,
		Email:       u.Email,
	}

	id := strconv.Itoa(s.Increment)
	s.List[id] = createUser

	b, err := json.Marshal(&s)
	if err != nil {
		r.Logger.Warnln("Marshal returns error")
		return "", apperror.ErrInternalServer
	}

	if err := ioutil.WriteFile(r.store, b, fs.ModePerm); err != nil {
		r.Logger.Warnln("WriteFile returns error")
		return "", apperror.ErrInternalServer
	}

	return id, nil
}

func (r *Repository) Update(id string, u user.UpdateUserRequest) error {

	f, err := ioutil.ReadFile(r.store)
	if err != nil {
		r.Logger.Warnln("ReadFile returns error")
		return apperror.ErrInternalServer
	}

	s := user.UserStore{}
	if err := json.Unmarshal(f, &s); err != nil {
		r.Logger.Warnln("Unmarshal returns error")
		return apperror.ErrInternalServer
	}

	if _, ok := s.List[id]; !ok {
		r.Logger.Infof("User doesn't exist with id: %s", id)
		return apperror.ErrUserNotFound
	}

	updateUser := s.List[id]
	updateUser.DisplayName = u.DisplayName
	s.List[id] = updateUser

	b, err := json.Marshal(&s)
	if err != nil {
		r.Logger.Warnln("Marshal returns error")
		return apperror.ErrInternalServer
	}

	if err := ioutil.WriteFile(r.store, b, fs.ModePerm); err != nil {
		r.Logger.Warnln("WriteFile returns error")
		return apperror.ErrInternalServer
	}

	return nil
}

func (r *Repository) Delete(id string) error {

	f, err := ioutil.ReadFile(r.store)
	if err != nil {
		r.Logger.Warnln("ReadFile returns error")
		return apperror.ErrInternalServer
	}

	s := user.UserStore{}
	if err := json.Unmarshal(f, &s); err != nil {
		r.Logger.Warnln("Unmarshal returns error")
		return apperror.ErrInternalServer
	}

	if _, ok := s.List[id]; !ok {
		r.Logger.Infof("User doesn't exist with id: %s", id)
		return apperror.ErrUserNotFound
	}

	delete(s.List, id)

	b, err := json.Marshal(&s)
	if err != nil {
		r.Logger.Warnln("Marshal returns error")
		return apperror.ErrInternalServer
	}

	if err := ioutil.WriteFile(r.store, b, fs.ModePerm); err != nil {
		r.Logger.Warnln("WriteFile returns error")
		return apperror.ErrInternalServer
	}

	return nil
}

func (r *Repository) GetAll() (user.UserStore, error) {

	s := user.UserStore{}

	f, err := ioutil.ReadFile(r.store)
	if err != nil {
		r.Logger.Warnln("ReadFile returns error")
		return s, apperror.ErrInternalServer
	}
	if err = json.Unmarshal(f, &s); err != nil {
		r.Logger.Warnln("Unmarshal returns error")
		return s, apperror.ErrInternalServer
	}

	return s, nil
}
