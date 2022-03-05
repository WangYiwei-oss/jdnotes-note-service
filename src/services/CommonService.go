package services

import (
	"github.com/WangYiwei-oss/jdframe/src/configs"
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
	"strings"
)

type CommonService struct {
	DB *configs.GormAdapter `inject:"-"`
}

func NewCommonService() *CommonService {
	return &CommonService{}
}

func (c *CommonService) GetUser(username string) *models.User {
	user := &models.User{
		User: jdft.User{
			UserName: username,
		},
	}
	c.DB.First(user)
	return user
}

func (c *CommonService) GetUserWithNote(username string) *models.User {
	user := &models.User{
		User: jdft.User{
			UserName: username,
		},
	}
	c.DB.Preload("Notes").First(user)
	return user
}

func (c *CommonService) ConvertNotesToObj(notes []models.Note) map[string]interface{} {
	ret := make(map[string]interface{})
	if notes == nil || len(notes) == 0 {
		return ret
	}
	currentMap := ret
	for _, note := range notes {
		relativePath := strings.TrimPrefix(note.NotePath, note.RootPath+"/")
		classes := strings.Split(relativePath, "/")
		currentMap = ret
		for i := 0; i < len(classes); i++ {
			if _, ok := currentMap[classes[i]]; ok {
				currentMap = currentMap[classes[i]].(map[string]interface{})
			} else {
				currentMap[classes[i]] = make(map[string]interface{})
				currentMap = currentMap[classes[i]].(map[string]interface{})
			}
		}
		currentMap["m_title_"+note.Title] = map[string]interface{}{
			"title": note.Title,
			"uuid":  note.UUID,
		}
	}
	return ret
}
