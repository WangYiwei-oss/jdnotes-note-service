package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/WangYiwei-oss/jdframe/src/configs"
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/services"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"path/filepath"
	"strings"
)

type NoteCtl struct {
	DB *configs.GormAdapter    `inject:"-"`
	ES *configs.EsAdapter      `inject:"-"`
	C  *services.CommonService `inject:"-"`
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
		Proto:    services.MD,
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
	ret := n.C.ConvertNotesToObj(user.Notes)
	return 1, ret
}

func (n *NoteCtl) GetNote(ctx *gin.Context) (int, string) {
	uuid := ctx.Query("uuid")
	if uuid == "" {
		return -300, ""
	}
	rep, err := n.ES.Get().Index("notes").Id(uuid).Do(ctx)
	if err != nil {
		return -301, err.Error()
	}
	if rep.Found {
		ret, _ := rep.Source.MarshalJSON()
		return 1, string(ret)
	} else {
		return -301, "查无此记录"
	}
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
