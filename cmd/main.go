package main

import (
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/controllers"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
)

func migration() {
	jdft.Gorm.AutoMigrate(&models.Note{})
	jdft.Gorm.AutoMigrate(&models.User{})
}

func main() {
	migration()
	//common.GetFdNotify().Mount("D:\\test").Start()
	jdft.NewJdft().DefaultBean().
		Mount("v1", controllers.NewNoteCtl()).
		Mount("v1", controllers.NewNoteView(), controllers.NewNotifyLogin()).
		Launch()
}
