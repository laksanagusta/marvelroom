package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"sandbox/internal/usecase/grc"
)

// GRCHandler handles HTTP requests for GRC dashboard
type GRCHandler struct {
	getOverviewUseCase   *grc.GetOverviewUseCase
	listUnitsUseCase     *grc.ListUnitsUseCase
	getUnitDetailUseCase *grc.GetUnitDetailUseCase
	compareUnitsUseCase  *grc.CompareUnitsUseCase
	getCategoriesUseCase *grc.GetCategoriesUseCase
}

// NewGRCHandler creates a new GRC handler
func NewGRCHandler(
	getOverviewUseCase *grc.GetOverviewUseCase,
	listUnitsUseCase *grc.ListUnitsUseCase,
	getUnitDetailUseCase *grc.GetUnitDetailUseCase,
	compareUnitsUseCase *grc.CompareUnitsUseCase,
	getCategoriesUseCase *grc.GetCategoriesUseCase,
) *GRCHandler {
	return &GRCHandler{
		getOverviewUseCase:   getOverviewUseCase,
		listUnitsUseCase:     listUnitsUseCase,
		getUnitDetailUseCase: getUnitDetailUseCase,
		compareUnitsUseCase:  compareUnitsUseCase,
		getCategoriesUseCase: getCategoriesUseCase,
	}
}

// GetOverview returns dashboard overview with statistics (live from spreadsheet)
func (h *GRCHandler) GetOverview(c *fiber.Ctx) error {
	response, err := h.getOverviewUseCase.Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve overview data",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"source":  "Google Spreadsheet (live)",
	})
}

// ListUnits returns all units with optional filtering and sorting
func (h *GRCHandler) ListUnits(c *fiber.Ctx) error {
	category := c.Query("category", "")
	sortBy := c.Query("sort_by", "id")
	ascending := c.Query("ascending", "false") == "true"

	response, err := h.listUnitsUseCase.Execute(category, sortBy, ascending)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve units",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetUnitDetail returns detailed information for a specific unit
func (h *GRCHandler) GetUnitDetail(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid unit ID",
		})
	}

	response, err := h.getUnitDetailUseCase.Execute(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve unit detail",
			"details": err.Error(),
		})
	}

	if response == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Unit not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// CompareUnits returns comparison data for multiple units
func (h *GRCHandler) CompareUnits(c *fiber.Ctx) error {
	idsStr := c.Query("ids", "")
	if idsStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "ids parameter is required",
		})
	}

	idStrs := strings.Split(idsStr, ",")
	var ids []int
	for _, s := range idStrs {
		id, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "No valid unit IDs provided",
		})
	}

	response, err := h.compareUnitsUseCase.Execute(ids)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to compare units",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetCategories returns category breakdown
func (h *GRCHandler) GetCategories(c *fiber.Ctx) error {
	response, err := h.getCategoriesUseCase.Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve categories",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}
