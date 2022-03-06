package controllers

import (
	"fmt"
	"github.com/WangYiwei-oss/jdframe/src/jdft"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/services"
	"github.com/gin-gonic/gin"
)

type NoteNotify struct {
	NoteService *services.NotifyProcessor `inject:"-"`
}

func NewNoteNotify() *NoteNotify {
	return &NoteNotify{}
}

func (n *NoteNotify) ProcessMessage(ctx *gin.Context) int {
	message := &models.NotifyMessage{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	fmt.Println(message)
	return 1
}

func (n *NoteNotify) DeleteFile(ctx *gin.Context) int {
	message := &models.NotifyMessage{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	n.NoteService.DelFile(ctx, message)
	return 1
}

func (n *NoteNotify) MoveDir(ctx *gin.Context) int {
	message := &models.NotifyMessage{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	n.NoteService.UpdateDir(message)
	return 1
}

func (n *NoteNotify) AddFile(ctx *gin.Context) int {
	message := &models.NotifyMessageWithContent{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	n.NoteService.AddFile(ctx, message)
	return 1
}

// UpdateFile 更新文件内容
func (n *NoteNotify) UpdateFile(ctx *gin.Context) int {
	message := &models.NotifyMessageWithContent{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	n.NoteService.UpdateFile(ctx, message)
	return 1
}

// MoveFile 移动文件
func (n *NoteNotify) MoveFile(ctx *gin.Context) int {
	message := &models.NotifyMessage{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	err = n.NoteService.MoveFile(message)
	if err != nil {
		return -400
	}
	return 1
}

func (n *NoteNotify) ChangeRootDir(ctx *gin.Context) int {
	message := &models.ChangeUserRootModel{}
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	n.NoteService.ChangeRoot(message)
	return 1
}

func (n *NoteNotify) DeleteAllUserNote(ctx *gin.Context) int {
	username := &models.UserNameModel{}
	err := ctx.ShouldBindJSON(username)
	if err != nil {
		fmt.Println(err)
		return -400
	}
	n.NoteService.DeleteAllUserNote(ctx, username.UserName)
	return 1
}

func (n *NoteNotify) Build(jdft *jdft.Jdft) {
	jdft.Handle("DELETE", "notify/file", n.DeleteFile) //删除文件或文件夹
	jdft.Handle("PUT", "notify/dir", n.MoveDir)        //修改文件夹名
	jdft.Handle("POST", "notify/file", n.AddFile)      //增加文件
	jdft.Handle("PATCH", "notify/file", n.MoveFile)    //移动文件
	jdft.Handle("PUT", "notify/file", n.UpdateFile)    //修改文件内容
	jdft.Handle("POST", "notify/test", n.ProcessMessage)
	jdft.Handle("PATCH", "notify/user_root", n.ChangeRootDir)       //用户改变根目录
	jdft.Handle("DELETE", "notify/user_notes", n.DeleteAllUserNote) //删除用户全部笔记
}
