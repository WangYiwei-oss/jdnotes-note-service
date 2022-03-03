package services

import (
	"fmt"
	"github.com/WangYiwei-oss/jdframe/src/configs"
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
	"path"
	"strings"
)

type NotifyProcessor struct {
	DB *configs.GormAdapter `inject:"-"`
}

func NewNotifyProcessor() *NotifyProcessor {
	return &NotifyProcessor{}
}

func (n *NotifyProcessor) AddFile() {

}

func (n *NotifyProcessor) UpdateFile() {

}

func (n *NotifyProcessor) DelFile(message *models.NotifyMessage) {
	user := &models.User{
		User: jdft.User{
			UserName: message.Username,
		},
	}
	n.DB.First(user)
	notes := make([]models.Note, 0)
	path := path.Join("/", "data", user.UserName, strings.TrimPrefix(message.SrcPath, user.RootPath))
	fmt.Println(path)
	n.DB.Raw("select * from notes where user_id = ? and note_path LIKE '"+path+"%'", user.ID).Find(&notes)
	fmt.Println(notes)
}

func (n *NotifyProcessor) AddDir() {
	fmt.Println("增加了文件夹，但我不管")
}

func (n *NotifyProcessor) UpdateDir(message *models.NotifyMessage) {
	user := &models.User{
		User: jdft.User{
			UserName: message.Username,
		},
	}
	n.DB.First(user)
	prefix, suffix := "", ""
	srcPath := path.Join("/", "data", user.UserName, strings.TrimPrefix(message.SrcPath, user.RootPath))
	dstPath := path.Join("/", "data", user.UserName, strings.TrimPrefix(message.DestPath, user.RootPath))
	srcRune := []rune(srcPath)
	dstRune := []rune(dstPath)
	i := 0
	for ; i < len(srcRune); i++ {
		if srcRune[i] == dstRune[i] {
			prefix += string(srcRune[i])
		} else {
			break
		}
	}
	j := len(srcRune) - 1
	k := len(dstRune) - 1
	for k >= 0 && j >= 0 {
		if srcRune[j] == dstRune[k] {
			suffix += string(srcRune[j])
		} else {
			break
		}
		j--
		k--
	}

	prefix = prefix[:strings.LastIndex(prefix, "/")+1]
	suffix = reverseString(suffix)
	suffix = suffix[strings.Index(suffix, "/"):]
	fmt.Println(prefix, suffix)
	//notes:=make([]models.Note,0)

}
