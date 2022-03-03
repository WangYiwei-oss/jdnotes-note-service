package controllers

import (
	"fmt"
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/services"
	"github.com/gin-gonic/gin"
)

type NoteNotify struct {
	NoteService *services.NotifyProcessor `inject:"-"`
}

func NewNoteNotify() *NoteNotify {
	return &NoteNotify{}
}

func (n *NoteNotify) ProcessMessage(ctx *gin.Context) int {
	message := &models.NotifyMessage{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	fmt.Println(message)
	return 1
}

func (n *NoteNotify) DeleteFile(ctx *gin.Context) int {
	message := &models.NotifyMessage{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	n.NoteService.DelFile(message)
	return 1
}

func (n *NoteNotify) MoveDir(ctx *gin.Context) int {
	message := &models.NotifyMessage{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	n.NoteService.UpdateDir(message)
	return 1
}

func (n *NoteNotify) Build(jdft *jdft.Jdft) {
	jdft.Handle("DELETE", "notify/note", n.DeleteFile)
	jdft.Handle("PUT", "notify/note", n.MoveDir)
	jdft.Handle("POST", "note/notify", n.ProcessMessage)
}
