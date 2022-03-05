package services

import (
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
	"path"
	"strings"
)

func reverseString(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}

	return string(runes)
}

// GetUserRootPath 根据用户名获取笔记的rootpath
func GetUserRootPath(user *models.User) string {
	return path.Join("/data", user.UserName)
}

// ReplaceRootPath 把本地空间的rootpath换成用户空间的rootpath
func ReplaceRootPath(user *models.User, srcPath string) string {
	return path.Join(GetUserRootPath(user), strings.TrimPrefix(srcPath, user.RootPath))
}

// DeleteExt 删除文件后缀
func DeleteExt(filename string) string {
	if index := strings.LastIndex(filename, "."); index != -1 {
		return filename[:index]
	}
	return filename
}
