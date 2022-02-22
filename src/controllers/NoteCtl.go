package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/WangYiwei-oss/jdframe/src/configs"
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"path/filepath"
	"strings"
)

type NoteCtl struct {
	DB *configs.GormAdapter `inject:"-"`
	ES *configs.EsAdapter   `inject:"-"`
}

func NewNoteCtl() *NoteCtl {
	return &NoteCtl{}
}

func (n *NoteCtl) CreateNote(ctx *gin.Context) (int, string) {
	var notePost models.NotePost
	err := ctx.ShouldBindJSON(&notePost)
	if err != nil {
		return -300, err.Error()
	}
	fmt.Println(notePost)
	user := &models.User{
		User: jdft.User{
			UserName: notePost.User,
		},
	}
	n.DB.Preload("notes").First(user)
	note := models.Note{
		Title:    notePost.Title,
		RootPath: "/" + strings.Replace(filepath.Join("data", notePost.User), "\\", "/", -1),
		NotePath: "/" + strings.Replace(filepath.Join("data", notePost.User, notePost.FirstClass, notePost.SecondClass, notePost.ThirdClass), "\\", "/", -1),
		NoteType: notePost.NoteType,
		UUID:     uuid.NewV1().String(),
	}
	user.Notes = append(user.Notes, note)
	n.DB.Save(user)
	esNote := &models.EsNote{
		Text: notePost.Text,
	}
	esJson, _ := json.Marshal(esNote)
	n.ES.Index().Index("notes").Id(note.UUID).BodyString(string(esJson)).Do(ctx)
	return 1, ""
}

func (n *NoteCtl) DeleteNote(ctx *gin.Context) (int, string) {
	return 1, "删除笔记"
}

type NoteModel struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
}

func (n *NoteCtl) convertNotesToObj(notes []models.Note) map[string]interface{} {
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
		currentMap["m_title"] = map[string]interface{}{
			"title": note.Title,
			"uuid":  note.UUID,
		}
	}
	return ret
}

func (n *NoteCtl) GetUserNotes(ctx *gin.Context) (int, jdft.Json) {
	userName := ctx.Query("user")
	if userName == "" {
		return -300, ""
	}
	user := &models.User{
		User: jdft.User{
			UserName: userName,
		},
	}
	n.DB.Preload("Notes").First(user)
	ret := n.convertNotesToObj(user.Notes)
	return 1, ret
}

func (n *NoteCtl) GetNote(ctx *gin.Context) (int, string) {
	//userName,uuid := ctx.Query("user"),ctx.Query("uuid")

	return 1, "获取笔记内容了"
}

func (n *NoteCtl) UpdateNote(ctx *gin.Context) (int, string) {
	return 1, "修改笔记内容了"
}

func (n *NoteCtl) Build(jdft *jdft.Jdft) {
	jdft.Handle("POST", "note", n.CreateNote)
	jdft.Handle("DELETE", "note", n.DeleteNote)
	jdft.Handle("GET", "note", n.GetNote)
	jdft.Handle("PUT", "note", n.UpdateNote)
	jdft.Handle("GET", "notes", n.GetUserNotes)
}
