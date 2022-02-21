package controllers

import (
	"github.com/WangYiwei-oss/jdframe/src/configs"
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/gin-gonic/gin"
)

type NoteView struct {
	db *configs.GormAdapter
}

func NewNoteView() *NoteView {
	return &NoteView{}
}

func (n *NoteView) CreateNote(ctx *gin.Context) (int, string) {
	return 1, "获取所有笔记列表了"
}

func (n *NoteView) Build(jdft *jdft.Jdft) {
	jdft.Handle("GET", "note_list", n.CreateNote)
}
