package controllers

import (
	"net/http"

	"github.com/aikintech/scim/pkg/dto"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
)

func SignIn(c echo.Context) error {
	input := new(dto.SignInDTO)

	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "An error occurred while parsing your request",
		})
	}

	validator := validate.Struct(input)

	if validator.Validate() {
		return c.JSON(http.StatusOK, input)
	}

	return c.JSON(http.StatusUnprocessableEntity, utils.FormatValidationErrors(validator.Errors.All()))
}
