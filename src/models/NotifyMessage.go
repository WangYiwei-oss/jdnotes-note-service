package models

type NotifyMessage struct {
	Username    string `json:"username"`
	Action      string `json:"action"`
	IsDirectory bool   `json:"is_directory"`
	SrcPath     string `json:"src_path"`
	DestPath    string `json:"dest_path"`
}
