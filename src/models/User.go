package models

import "github.com/WangYiwei-oss/jdframe/src/jdft"

type User struct {
	jdft.User
	Notes []Note
}
