package controllers

import (
	"github.com/gin-gonic/gin"
)

type UserStatus struct {
	IsLoggedIn bool
}

func NewUserStatus(ctx *gin.Context) UserStatus {
	userId, exists := ctx.Get("UserId")
	if !exists {
		return UserStatus{IsLoggedIn: false}
	}
	uID, ok := userId.(uint)
	return UserStatus{IsLoggedIn: ok && uID != 0}
}
