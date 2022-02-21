package controllers

import (
	"github.com/WangYiwei-oss/jdframe/src/configs"
)

type NoteNotify struct {
	db *configs.GormAdapter
}

func NewNoteNotify() *NoteNotify {
	return &NoteNotify{}
}

type NotifyInfo struct {
}

//func (n *NoteNotify) Build(jdft *jdft.Jdft) {
//	jdft.Handle("GET", "note_list", n.CreateNote)
//}
