package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"altoai_mvp/internal/middleware"
	"altoai_mvp/internal/models"
	"altoai_mvp/internal/repository"
	"altoai_mvp/internal/services"
	"altoai_mvp/pkg/errors"
	"altoai_mvp/pkg/response"
)

type UserHandler struct {
	svc services.UserService
}

func NewUserHandler(svc services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to list users")
		return
	}
	response.OK(c, users)
}

func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	u, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			response.Error(c, http.StatusNotFound, "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get user")
		return
	}
	response.OK(c, u)
}

func (h *UserHandler) Create(c *gin.Context) {
	var dto models.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, errs.FromBinding(err))
		return
	}
	u, err := h.svc.Create(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create user")
		return
	}
	response.Created(c, u)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var dto models.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, errs.FromBinding(err))
		return
	}
	u, err := h.svc.Update(c.Request.Context(), id, dto)
	if err != nil {
		if err == repository.ErrNotFound {
			response.Error(c, http.StatusNotFound, "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to update user")
		return
	}
	response.OK(c, u)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if err == repository.ErrNotFound {
			response.Error(c, http.StatusNotFound, "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to delete user")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	claims := c.MustGet("user").(*middleware.MyClaims)
	// Get user by email
	user, err := h.svc.GetByEmail(c.Request.Context(), claims.Email)
	if err != nil {
		if err == repository.ErrNotFound {
			response.Error(c, http.StatusNotFound, "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get user")
		return
	}
	
	var dto models.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ValidationError(c, errs.FromBinding(err))
		return
	}
	
	u, err := h.svc.Update(c.Request.Context(), user.ID, dto)
	if err != nil {
		if err == repository.ErrNotFound {
			response.Error(c, http.StatusNotFound, "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to update profile")
		return
	}
	response.OK(c, u)
}
