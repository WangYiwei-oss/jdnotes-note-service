package main

import (
	"fmt"
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/config"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/controllers"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
)

func migration() {
	jdft.Gorm.AutoMigrate(&models.Note{})
	err := jdft.Gorm.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("迁移user表错误", err)
	}
}

func main() {
	migration()
	//common.GetFdNotify().Mount("D:\\test").Start()
	jdft.NewJdft().
		DefaultBean().
		Beans(config.NewMServiceConfig()).
		Mount("v1", controllers.NewNoteCtl()).
		Mount("v1", controllers.NewNoteView(), controllers.NewNotifyLogin(), controllers.NewNoteNotify()).
		Launch()
}
