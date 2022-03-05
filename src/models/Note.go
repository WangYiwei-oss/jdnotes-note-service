package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Title    string `gorm:"column:title;type:varchar(20);NOT NULL" json:"title"`
	RootPath string `gorm:"column:root_path;type:varchar(50);NOT NULL" json:"root_path"`
	NotePath string `gorm:"column:note_path;type:varchar(50);NOT NULL" json:"note_path"`
	NoteType uint   `gorm:"column:note_type;type:uint;NOT NULL" json:"note_type"`
	UUID     string `gorm:"column:uuid;type:varchar(50);NOT NULL" json:"uuid"`
	UserID   uint   //外键索引
}

func (n *Note) TableName() string {
	return "notes"
}

type NoteList []*Note

type NotePost struct {
	User        string `json:"user" binding:"required"`
	Title       string `json:"title" binding:"required"`
	FirstClass  string `json:"first_class" binding:"required"`
	SecondClass string `json:"second_class"`
	ThirdClass  string `json:"third_class"`
	NoteType    uint   `json:"note_type" binding:"required"`
	Text        string `json:"text" binding:"required"`
}

type EsNote struct {
	Text string `json:"text"`
}
