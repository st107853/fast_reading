package controllers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/logger"
	"github.com/st107853/fast_reading/services"
)

type UserStatus struct {
	IsLoggedIn bool
}

type ErrorData struct {
	UserStatus
	services.DomainError
}

var errorPage = template.Must(template.New("error_page.html").ParseFiles("./static/error_page.html", "./static/template.html"))

func NewUserStatus(ctx *gin.Context) UserStatus {
	userId, exists := ctx.Get("UserId")
	if !exists {
		return UserStatus{IsLoggedIn: false}
	}
	uID, ok := userId.(uint)
	return UserStatus{IsLoggedIn: ok && uID != 0}
}

func RenderError(ctx *gin.Context, err error) {
	// Build structured log fields shared across all error logs
	de, ok := services.IsDomainError(err)
	if !ok {
		de = services.ErrInternal(err)
	}

	log := logger.FromContext(ctx)

	switch {
	case de.Code >= 500:
		log.Error("internal server error",
			slog.Int("status", de.Code),
			slog.Any("error", de.InternalErr),
		)
	case de.Code == http.StatusForbidden ||
		de.Code == http.StatusUnauthorized:
		log.Info("access denied",
			slog.Int("status", de.Code),
		)
	default:
		log.Debug("client error",
			slog.Int("status", de.Code),
		)
	}

	data := ErrorData{
		UserStatus:  NewUserStatus(ctx),
		DomainError: *de,
	}

	if err := errorPage.Execute(ctx.Writer, data); err != nil {
		log.Error("failed to render error page",
			slog.Any("error", err),
		)
		ctx.String(de.Code, de.UserMsg)
	}
}
