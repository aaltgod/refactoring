package handler

import (
	"net/http"

	"github.com/aaltgod/refactoring/internal/user"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	UserStorage user.Repository
	Logger      *logrus.Logger
}

func NewHandler(storage user.Repository, logger *logrus.Logger) Handler {
	return Handler{
		UserStorage: storage,
		Logger:      logger,
	}
}

func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) error {

	users, err := h.UserStorage.GetAll()
	if err != nil {
		h.Logger.Warnf("GetAll returns error: %s", err.Error())
		return err
	}

	render.JSON(w, r, users.List)
	return nil
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {

	id := r.Context().Value("id").(string)

	user, err := h.UserStorage.Get(id)
	if err != nil {
		h.Logger.Warnf("Get returns error: %s", err.Error())
		return err
	}

	render.JSON(w, r, user)
	return nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {

	request := user.CreateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		h.Logger.Warnf("Bind returns error: %s", err.Error())
		return err
	}

	id, err := h.UserStorage.Insert(request)
	if err != nil {
		h.Logger.Warnf("Insert returns error: %s", err.Error())
		return err
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})

	return nil
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {

	request := user.UpdateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		h.Logger.Warnf("Bind returns error: %s", err.Error())
		return err
	}

	id := r.Context().Value("id").(string)

	if err := h.UserStorage.Update(id, request); err != nil {
		h.Logger.Warnf("Update returns error: %s", err.Error())
		return err
	}

	render.Status(r, http.StatusNoContent)
	return nil
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {

	id := r.Context().Value("id").(string)

	if err := h.UserStorage.Delete(id); err != nil {
		h.Logger.Warnf("Delete returns error: %s", err.Error())
		return err
	}

	render.Status(r, http.StatusNoContent)
	return nil
}
