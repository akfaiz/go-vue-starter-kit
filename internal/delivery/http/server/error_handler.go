package server

import (
	"fmt"
	"net/http"

	"github.com/akfaiz/go-vue-starter-kit/internal/errdefs"
	"github.com/akfaiz/go-vue-starter-kit/internal/validator"
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
)

const ContentTypeProblemJSON = "application/problem+json"

func customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return // If the response is already committed, do nothing
	}

	instance := c.Path()
	requestID, ok := c.Response().Header()[echo.HeaderXRequestID]
	if ok && len(requestID) > 0 {
		instance = requestID[0]
	}

	// Check if the error is a custom application error
	var appError *errdefs.AppError
	if errors.As(err, &appError) {
		c.JSON(appError.Status, appError.WithInstance(instance))
		c.Response().Header().Set(echo.HeaderContentType, ContentTypeProblemJSON)
		return
	}

	// Check if the error is a validation error
	var validationErr *validator.ValidationError
	if errors.As(err, &validationErr) {
		appError := errdefs.ErrValidation().
			WithErrors(validationErr).
			WithCause(err).
			WithInstance(instance)
		c.JSON(appError.Status, appError)
		c.Response().Header().Set(echo.HeaderContentType, ContentTypeProblemJSON)
		return
	}

	code := http.StatusInternalServerError
	// Retrieve the custom status code if it's a *echo.HTTPError
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		code = httpErr.Code
		appError = errdefs.New(fmt.Sprintf("%v", httpErr.Message), "about:blank", code)
	} else {
		appError = errdefs.ErrInternalServer()
	}
	c.JSON(code, appError.WithInstance(instance))
	c.Response().Header().Set(echo.HeaderContentType, ContentTypeProblemJSON)
}
