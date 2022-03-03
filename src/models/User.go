package models

import "github.com/WangYiwei-oss/jdframe/src/jdft"

type User struct {
	jdft.User
	RootPath string `json:"root_path" gorm:"column:root_path"`
	Notes    []Note
}
