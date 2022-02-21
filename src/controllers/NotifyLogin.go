package controllers

import (
	"fmt"
	"github.com/WangYiwei-oss/jdframe/src/configs"
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/gin-gonic/gin"
)

type NotifyLogin struct {
	db *configs.GormAdapter
}

func NewNotifyLogin() *NotifyLogin {
	return &NotifyLogin{}
}

type LoginPayLoad struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (n *NotifyLogin) Login(ctx *gin.Context) int {
	loginPayLoad := LoginPayLoad{}
	if err := ctx.BindJSON(&loginPayLoad); err != nil {
		return 422
	}
	fmt.Println(loginPayLoad)
	return 1
}

func (n *NotifyLogin) Build(jdft *jdft.Jdft) {
	jdft.Handle("POST", "notify_login", n.Login)
}
