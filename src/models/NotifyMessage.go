package models

type NotifyMessage struct {
	Username    string `json:"username"`
	Action      string `json:"action"`
	IsDirectory bool   `json:"is_directory"`
	SrcPath     string `json:"src_path"`
	DestPath    string `json:"dest_path"`
	Proto       string `json:"proto"`
}

type NotifyMessageWithContent struct {
	NotifyMessage
	Content string `json:"content"`
}

type ChangeUserRootModel struct {
	Username string `json:"username"`
	NewRoot  string `json:"new_root"`
}
