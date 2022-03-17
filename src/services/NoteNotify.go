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
	log.Println("触发创建文件", message.SrcPath, message.DestPath, message.IsDirectory)
	user := n.C.GetUserWithNote(message.Username)
	index := strings.LastIndex(message.SrcPath, "/")
	rootPath := GetUserRootPath(user)
	notePath := ReplaceRootPath(user, message.SrcPath[:index+1])
	if rootPath == notePath {
		return
	}
	note := models.Note{
		Title:    DeleteExt(message.SrcPath[index+1:]),
		RootPath: rootPath,
		NotePath: notePath,
		UUID:     uuid.NewV1().String(),
		Proto:    message.Proto,
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
	log.Println("触发改变文件", message.SrcPath, message.DestPath, message.IsDirectory)
	user := n.C.GetUser(message.Username)
	index := strings.LastIndex(message.SrcPath, "/")
	title := DeleteExt(message.SrcPath[index+1:])
	notePath := ReplaceRootPath(user, message.SrcPath[:index])
	note := &models.Note{}
	n.DB.Model(&models.Note{}).
		Where("user_id = ?", user.ID).
		Where("title = ?", title).
		Where("note_path = ?", notePath).First(note)
	if note.ID == 0 { //复制一个文件进来时也会触发modify，所以如果找不到就添加
		n.AddFile(ctx, message)
		return
	}
	esNote := &models.EsNote{
		Text: message.Content,
	}
	esJson, _ := json.Marshal(esNote)
	n.ES.Index().Index("notes").Id(note.UUID).BodyString(string(esJson)).Do(ctx)
}

func (n *NotifyProcessor) MoveFile(message *models.NotifyMessage) error {
	log.Println("触发移动文件", message.SrcPath, message.DestPath, message.IsDirectory)
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

func (n *NotifyProcessor) DelFile(ctx *gin.Context, message *models.NotifyMessage) {
	log.Println("触发删除文件", message.SrcPath, message.DestPath, message.IsDirectory)
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
			n.DB.Unscoped().Delete(&models.Note{}, note.ID)
			n.ES.Delete().Id(note.UUID).Do(ctx) //es删除
		}
	}
}

func (n *NotifyProcessor) AddDir() {
	fmt.Println("增加了文件夹，但我不管")
}

func (n *NotifyProcessor) UpdateDir(message *models.NotifyMessage) {
	log.Println("触发改变文件夹", message.SrcPath, message.DestPath, message.IsDirectory)
	user := n.C.GetUser(message.Username)
	//prefix, suffix := "/", ""
	srcPath := path.Join("/", "data", user.UserName, strings.TrimPrefix(message.SrcPath, user.RootPath))
	dstPath := path.Join("/", "data", user.UserName, strings.TrimPrefix(message.DestPath, user.RootPath))
	//srcSlice := strings.Split(srcPath, "/")
	//dstSlice := strings.Split(dstPath, "/")
	//fmt.Println(srcSlice, dstSlice)
	//i := 0
	//for ; i < len(srcSlice); i++ {
	//	if srcSlice[i] == dstSlice[i] {
	//		prefix = path.Join(prefix, srcSlice[i])
	//	} else {
	//		break
	//	}
	//}
	//for j := i + 1; j < len(srcSlice); j++ {
	//	suffix += dstSlice[j]
	//}
	notes := make([]models.Note, 0)
	//fmt.Println(prefix,suffix)
	//n.DB.Raw(fmt.Sprintf("select * from notes where user_id=1 and note_path LIKE CONCAT ('%s%%','%%%s');", prefix, suffix)).Find(&notes)
	n.DB.Raw(fmt.Sprintf("select * from notes where user_id=1 and note_path='%s';", srcPath)).Find(&notes)
	for _, note := range notes {
		//if len(strings.Split(note.NotePath, "/")) == len(srcSlice) {
		//	note.NotePath = path.Join(prefix, dstSlice[i], suffix)
		//	n.DB.Save(&note)
		//}
		note.NotePath = dstPath
		n.DB.Save(&note)
	}
}

func (n *NotifyProcessor) DeleteAllUserNote(ctx *gin.Context, username string) {
	log.Println("触发删除用户所有文件")
	user := n.C.GetUserWithNote(username)
	for _, note := range user.Notes {
		if note.NoteType == 0 {
			fmt.Println("删除笔记", note.ID)
			n.DB.Unscoped().Delete(&note)
			n.ES.Delete().Id(note.UUID).Do(ctx)
		}
	}
}

func (n *NotifyProcessor) ChangeRoot(model *models.ChangeUserRootModel) {
	fmt.Println("====================更换源", model)
	n.DB.Model(&models.User{}).Where("user_name = ?", model.Username).Update("root_path", model.NewRoot)
}
