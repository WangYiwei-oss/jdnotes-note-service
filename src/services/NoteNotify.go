package services

import (
	"encoding/json"
	"fmt"
	"github.com/WangYiwei-oss/jdframe/src/configs"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"path"
	"strings"
)

type NotifyProcessor struct {
	DB *configs.GormAdapter `inject:"-"`
	C  *CommonService       `inject:"-"`
	ES *configs.EsAdapter   `inject:"-"`
}

func NewNotifyProcessor() *NotifyProcessor {
	return &NotifyProcessor{}
}

func (n *NotifyProcessor) AddFile(ctx *gin.Context, message *models.NotifyMessageWithContent) {
	user := n.C.GetUserWithNote(message.Username)
	index := strings.LastIndex(message.SrcPath, "/")
	note := models.Note{
		Title:    DeleteExt(message.SrcPath[index+1:]),
		RootPath: GetUserRootPath(user),
		NotePath: ReplaceRootPath(user, message.SrcPath[:index+1]),
		UUID:     uuid.NewV1().String(),
	}
	user.Notes = append(user.Notes, note)
	n.DB.Save(user)
	esNote := &models.EsNote{
		Text: message.Content,
	}
	esJson, _ := json.Marshal(esNote)
	n.ES.Index().Index("notes").Id(note.UUID).BodyString(string(esJson)).Do(ctx)
}

func (n *NotifyProcessor) UpdateFile(ctx *gin.Context, message *models.NotifyMessageWithContent) {
	user := n.C.GetUser(message.Username)
	index := strings.LastIndex(message.SrcPath, "/")
	title := DeleteExt(message.SrcPath[index+1:])
	notePath := ReplaceRootPath(user, message.SrcPath[:index])
	note := &models.Note{}
	n.DB.Model(&models.Note{}).
		Where("user_id = ?", user.ID).
		Where("title = ?", title).
		Where("note_path = ?", notePath).First(note)
	esNote := &models.EsNote{
		Text: message.Content,
	}
	esJson, _ := json.Marshal(esNote)
	n.ES.Index().Index("notes").Id(note.UUID).BodyString(string(esJson)).Do(ctx)
}

func (n *NotifyProcessor) MoveFile(message *models.NotifyMessage) error {
	user := n.C.GetUser(message.Username)
	index := strings.LastIndex(message.SrcPath, "/")
	title := DeleteExt(message.SrcPath[index+1:])
	newTitle := DeleteExt(message.DestPath[strings.LastIndex(message.DestPath, "/")+1:])
	notePath := ReplaceRootPath(user, message.SrcPath[:index])
	err := n.DB.Model(&models.Note{}).
		Where("user_id = ?", user.ID).
		Where("title = ?", title).
		Where("note_path = ?", notePath).
		Update("title", newTitle).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (n *NotifyProcessor) DelFile(message *models.NotifyMessage) {
	user := n.C.GetUser(message.Username)
	notes := make([]models.Note, 0)
	path := path.Join("/", "data", user.UserName, strings.TrimPrefix(message.SrcPath, user.RootPath))
	fmt.Println(path)
	n.DB.Raw("select * from notes where user_id = ? and note_path LIKE '"+path+"%';", user.ID).Find(&notes)
	if len(notes) == 0 {
		index := strings.LastIndex(message.SrcPath, "/")
		title := DeleteExt(message.SrcPath[index+1:])
		notePath := ReplaceRootPath(user, message.SrcPath[:index])
		fmt.Println("============", user.ID, title, notePath)
		n.DB.Exec(fmt.Sprintf("DELETE FROM notes where user_id=%d and title='%s' and note_path='%s'", user.ID, title, notePath))
	} else {
		for _, note := range notes {
			n.DB.Delete(&models.Note{}, note.ID)
		}
	}
}

func (n *NotifyProcessor) AddDir() {
	fmt.Println("增加了文件夹，但我不管")
}

func (n *NotifyProcessor) UpdateDir(message *models.NotifyMessage) {
	user := n.C.GetUser(message.Username)
	prefix, suffix := "/", ""
	srcPath := path.Join("/", "data", user.UserName, strings.TrimPrefix(message.SrcPath, user.RootPath))
	dstPath := path.Join("/", "data", user.UserName, strings.TrimPrefix(message.DestPath, user.RootPath))
	srcSlice := strings.Split(srcPath, "/")
	dstSlice := strings.Split(dstPath, "/")
	fmt.Println(srcSlice, dstSlice)
	i := 0
	for ; i < len(srcSlice); i++ {
		if srcSlice[i] == dstSlice[i] {
			prefix = path.Join(prefix, srcSlice[i])
		} else {
			break
		}
	}
	for j := i + 1; j < len(dstSlice); j++ {
		suffix += dstSlice[j]
	}
	notes := make([]models.Note, 0)
	n.DB.Raw(fmt.Sprintf("select * from notes where user_id=1 and note_path LIKE CONCAT ('%s%%','%%%s');", prefix, suffix)).Find(&notes)
	for _, note := range notes {
		if len(strings.Split(note.NotePath, "/")) == len(srcSlice) {
			note.NotePath = path.Join(prefix, dstSlice[i], suffix)
			n.DB.Save(&note)
		}
	}
}
