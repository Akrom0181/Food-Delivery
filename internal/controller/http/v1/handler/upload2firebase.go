package handler

import (
	"github.com/Akrom0181/Food-Delivery/config"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/Akrom0181/Food-Delivery/pkg/firebase"
	"github.com/gin-gonic/gin"
)

// UploadFiles godoc
// @ID upload_multiple_files
// @Router /firebase [post]
// @Summary Upload Multiple Files
// @Description Upload Multiple Files
// @Security BearerAuth
// @Tags Upload File
// @Accept multipart/form-data
// @Produce json
// @Param file formData []file true "File to upload"
// @Success 200 {object} entity.MultipleFileUploadResponse "Success Request"
// @Failure 400 {object} entity.ErrorResponse "Bad Request"
// @Failure 500 {object} entity.ErrorResponse "Server error"
func (h *Handler) UploadFiles(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid file upload request", 400)
		return
	}

	resp, err := firebase.UploadFiles(form)
	if h.HandleDbError(ctx, err, "Error uploading files") {
		return
	}

	ctx.JSON(200, resp)
}

// DeleteFile godoc
// @ID delete_file
// @Router /firebase/{id} [delete]
// @Summary Delete File
// @Description Delete File
// @Security BearerAuth
// @Tags Upload File
// @Accept json
// @Produce json
// @Param id query string true "ID of the file to delete"
// @Success 204 {string} string "Success Request"
// @Failure 400 {object} entity.ErrorResponse "Bad Request"
// @Failure 500 {object} entity.ErrorResponse "Server error"
func (h *Handler) DeleteFile(ctx *gin.Context) {
	fileID := ctx.Query("id")
	if fileID == "" {
		h.ReturnError(ctx, config.ErrorBadRequest, "Missing file ID", 400)
		return
	}

	err := firebase.DeleteFile(fileID)
	if h.HandleDbError(ctx, err, "Error deleting file") {
		return
	}

	ctx.JSON(204, entity.SuccessResponse{Message: "File deleted successfully"})
}
