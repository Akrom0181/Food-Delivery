package handler

import (
	"strconv"

	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/Akrom0181/Food-Delivery/pkg/firebase"
	"github.com/Akrom0181/Food-Delivery/pkg/hash"
	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Router /user [post]
// @Summary Create a new user
// @Description Create a new user
// @Security BearerAuth
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body entity.User true "User object"
// @Success 201 {object} entity.User
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateUser(ctx *gin.Context) {
	var (
		body entity.User
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	body.Password, err = hash.HashPassword(body.Password)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Error hashing password", 400)
		return
	}

	user, err := h.UseCase.UserRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating user") {
		return
	}

	ctx.JSON(201, user)
}

// GetUser godoc
// @Router /user/{id} [get]
// @Summary Get a user by ID
// @Description Get a user by ID
// @Security BearerAuth
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} entity.User
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetUser(ctx *gin.Context) {
	var (
		req entity.UserSingleRequest
	)

	req.ID = ctx.Param("id")

	user, err := h.UseCase.UserRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting user") {
		return
	}

	user.Password = ""

	ctx.JSON(200, user)
}

// GetUsers godoc
// @Router /user/list [get]
// @Summary Get a list of users
// @Description Get a list of users
// @Security BearerAuth
// @Tags user
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param search query string false "search"
// @Success 200 {object} entity.UserList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetUsers(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "full_name",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "username",
			Type:   "search",
			Value:  search,
		},
		entity.Filter{
			Column: "email",
			Type:   "search",
			Value:  search,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	users, err := h.UseCase.UserRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting users") {
		return
	}

	ctx.JSON(200, users)
}

// UpdateUser godoc
// @Router /user [put]
// @Summary Update a user
// @Description Update a user
// @Security BearerAuth
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body entity.User true "User object"
// @Success 200 {object} entity.User
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateUser(ctx *gin.Context) {
	var (
		body entity.User
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_type") == "user" {
		body.ID = ctx.GetHeader("sub")
	}

	if body.Password != "" {
		body.Password, err = hash.HashPassword(body.Password)
		if err != nil {
			h.ReturnError(ctx, config.ErrorBadRequest, "Error hashing password", 400)
			return
		}
	}

	user, err := h.UseCase.UserRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating user") {
		return
	}

	ctx.JSON(200, user)
}

// DeleteUser godoc
// @Router /user/{id} [delete]
// @Summary Delete a user
// @Description Delete a user
// @Security BearerAuth
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteUser(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	if ctx.GetHeader("user_type") == "user" {
		req.ID = ctx.GetHeader("sub")
	}

	err := h.UseCase.UserRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting user") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "User deleted successfully",
	})
}

// UploadProfilePic godoc
// @ID upload_profile_pic_file
// @Router /user/upload [post]
// @Summary Upload Multiple Files
// @Description Upload Multiple Files
// @Security BearerAuth
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param file formData []file true "File to upload"
// @Success 200 {object} entity.MultipleFileUploadResponse "Success Request"
// @Failure 400 {object} entity.ErrorResponse "Bad Request"
// @Failure 500 {object} entity.ErrorResponse "Server error"
func (h *Handler) UploadProfilePic(ctx *gin.Context) {
	var (
		id entity.UserSingleRequest
	)

	form, err := ctx.MultipartForm()
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid file upload request", 400)
		return
	}

	if ctx.GetHeader("user_type") == "user" || ctx.GetHeader("user_type") == "admin" {
		id.ID = ctx.GetHeader("sub")
	}

	user, err := h.UseCase.UserRepo.GetSingle(ctx, id)
	if h.HandleDbError(ctx, err, "Error getting user") {
		return
	}

	resp, err := firebase.UploadFiles(form)
	if h.HandleDbError(ctx, err, "Error uploading files") {
		return
	}

	_, err = h.UseCase.UserRepo.Update(ctx, entity.User{
		ID:         id.ID,
		FullName:   user.FullName,
		UserName:   user.UserName,
		Email:      user.Email,
		Gender:     user.Gender,
		ProfilePic: resp.Url[0].Url,
		Bio:        user.Bio,
		Status:     user.Status,
		UserType:   user.UserType,
		UserRole:   user.UserRole,
	})

	if h.HandleDbError(ctx, err, "Error updating user") {
		return
	}

	ctx.JSON(200, resp)
}
