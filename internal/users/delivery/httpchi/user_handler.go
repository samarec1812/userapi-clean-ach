package httpchi

import (
	"errors"
	"log"
	"net/http"
	"time"

	"refactoring/internal/domain"
	uerrors "refactoring/internal/users/repository/ujson"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewUserHandler(r chi.Router, userUsecase domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: userUsecase,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(time.Now().String()))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", handler.searchUsers)
				r.Post("/", handler.createUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", handler.getUser)
					r.Patch("/", handler.updateUser)
					r.Delete("/", handler.deleteUser)
				})
			})
		})
	})
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {

	request := CreateUserRequest{}
	if err := render.Bind(r, &request); err != nil {
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Print(err)
		}
		return
	}

	user, err := h.UserUsecase.Create(r.Context(), request.DisplayName, request.Email)
	if err != nil {
		err = render.Render(w, r, ErrorResponse(err, http.StatusInternalServerError, ServerErrorStatusText))
		if err != nil {
			log.Print(err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateUserResponse{ID: user.ID})
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.UserUsecase.GetByID(r.Context(), id)
	if err != nil {
		err = render.Render(w, r, ErrorResponse(err, http.StatusInternalServerError, ServerErrorStatusText))
		if err != nil {
			log.Print(err)
		}
		return
	}

	render.JSON(w, r, user)
}

func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	request := UpdateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		err = render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			log.Println(err)
		}
		return
	}

	id := chi.URLParam(r, "id")
	err := h.UserUsecase.Update(r.Context(), id, request.DisplayName)
	if err != nil {
		if errors.Is(err, uerrors.ErrUserNotFound) {
			err = render.Render(w, r, ErrInvalidRequest(UserNotFound))
			if err != nil {
				log.Println(err)
			}
			return
		}

		err = render.Render(w, r, ErrorResponse(err, http.StatusInternalServerError, ServerErrorStatusText))
		if err != nil {
			log.Println(err)
		}
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	err := h.UserUsecase.Delete(r.Context(), id)

	if err != nil {
		if errors.Is(err, uerrors.ErrUserNotFound) {
			err = render.Render(w, r, ErrInvalidRequest(UserNotFound))
			if err != nil {
				log.Print(err)
			}
			return
		}

		err = render.Render(w, r, ErrorResponse(err, http.StatusInternalServerError, ServerErrorStatusText))
		if err != nil {
			log.Print(err)
		}
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *UserHandler) searchUsers(w http.ResponseWriter, r *http.Request) {
	list, err := h.UserUsecase.GetAll(r.Context())
	if err != nil {
		if errors.Is(err, uerrors.ErrUserNotFound) {
			err = render.Render(w, r, ErrInvalidRequest(UserNotFound))
			if err != nil {
				log.Print(err)
			}
			return
		}

		err = render.Render(w, r, ErrorResponse(err, http.StatusInternalServerError, ServerErrorStatusText))
		if err != nil {
			log.Print(err)
		}
		return
	}
	render.JSON(w, r, list)
}
