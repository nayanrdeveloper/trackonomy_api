package category

import (
	"net/http"
	"strconv"
	"trackonomy/internal/dto"
	"trackonomy/internal/logger"
	"trackonomy/internal/response"
	"trackonomy/internal/utils"
	"trackonomy/internal/validators"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryController struct {
	service Service
}

func NewCategoryController(service Service) *CategoryController {
	return &CategoryController{service: service}
}

// CreateGlobalCategory creates a category that is for all users (is_global = true).
func (cc *CategoryController) CreateGlobalCategory(c *gin.Context) {
	// No auth check => userID = 0
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid category data", err.Error())
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		validationErrs := utils.ParseValidationErrors(err)
		response.BadRequest(c, "Validation error", validationErrs)
		return
	}

	cat := &Category{
		Name:     req.Name,
		IsGlobal: true, // Mark it global
		UserID:   0,    // userID=0 or no user
	}

	if err := cc.service.CreateCategory(cat); err != nil {
		response.InternalServerError(c, "Could not create global category", err.Error())
		return
	}
	response.Created(c, "Global category created successfully", cat)
}

// CreateCategory creates a new category.
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind CreateCategory JSON", zap.Error(err))
		response.BadRequest(c, "Invalid category data", err.Error())
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		// Convert validation errors to a map for a structured response
		validationErrs := utils.ParseValidationErrors(err)
		response.BadRequest(c, "Validation error", validationErrs)
		return
	}

	cat := &Category{
		Name:   req.Name,
		UserID: userID, // If categories belong to a user
	}

	if err := cc.service.CreateCategory(cat); err != nil {
		logger.Error("Failed to create category", zap.Error(err), zap.Uint("userID", userID))
		response.InternalServerError(c, "Could not create category", err.Error())
		return
	}
	response.Created(c, "Category created successfully", cat)
}

// GetAllGlobalCategories fetches only global categories (no user).
func (cc *CategoryController) GetAllGlobalCategories(c *gin.Context) {
	// userID=0 => repository returns is_global = true categories
	cats, err := cc.service.GetAllCategories(0)
	if err != nil {
		response.InternalServerError(c, "Could not retrieve global categories", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Global categories retrieved successfully", cats)
}

// GetAllCategories lists all categories for the user.
func (cc *CategoryController) GetAllCategories(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	categories, err := cc.service.GetAllCategories(userID)
	if err != nil {
		logger.Error("Failed to retrieve categories", zap.Error(err), zap.Uint("userID", userID))
		response.InternalServerError(c, "Could not retrieve categories", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Categories retrieved successfully", categories)
}

// GetCategoryByID retrieves a single category by ID.
func (cc *CategoryController) GetCategoryByID(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid category ID parameter", zap.String("id_param", idStr))
		response.BadRequest(c, "Invalid category ID", nil)
		return
	}

	category, err := cc.service.GetCategoryByID(uint(id), userID)
	if err != nil {
		logger.Error("Failed to retrieve category", zap.Error(err), zap.Int("categoryID", id))
		response.InternalServerError(c, "Failed to retrieve category", err.Error())
		return
	}
	if category == nil {
		response.NotFound(c, "Category not found", nil)
		return
	}
	response.Success(c, http.StatusOK, "Category retrieved successfully", category)
}

// UpdateCategory updates an existing category.
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid category ID parameter", zap.String("id_param", idStr))
		response.BadRequest(c, "Invalid category ID", nil)
		return
	}

	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind UpdateCategory JSON", zap.Error(err))
		response.BadRequest(c, "Invalid category data", err.Error())
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		validationErrs := utils.ParseValidationErrors(err)
		response.BadRequest(c, "Validation error", validationErrs)
		return
	}

	cat := &Category{
		ID:     uint(id),
		Name:   req.Name,
		UserID: userID, // keep user ownership
	}

	if err := cc.service.UpdateCategory(cat); err != nil {
		logger.Error("Failed to update category", zap.Error(err), zap.Uint("categoryID", cat.ID))
		response.InternalServerError(c, "Could not update category", err.Error())
		return
	}
	response.Updated(c, "Category updated successfully", cat)
}

// DeleteCategory removes a category by ID.
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid category ID parameter", zap.String("id_param", idStr))
		response.BadRequest(c, "Invalid category ID", nil)
		return
	}

	if err := cc.service.DeleteCategory(uint(id), userID); err != nil {
		logger.Error("Failed to delete category", zap.Error(err), zap.Int("categoryID", id))
		response.InternalServerError(c, "Could not delete category", err.Error())
		return
	}
	response.Deleted(c, "Category deleted successfully")
}
