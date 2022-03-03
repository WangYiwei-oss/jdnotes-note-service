package main

import (
	"fmt"
	"github.com/WangYiwei-oss/jdnotes-note-service/src/models"
	"strings"
)

func convertNotesToObj(note *models.Note) map[string]interface{} {
	ret := make(map[string]interface{})
	relativePath := strings.TrimPrefix(note.NotePath, note.RootPath+"/")
	classes := strings.Split(relativePath, "/")
	currentMap := ret
	for i := 0; i < len(classes); i++ {
		if _, ok := currentMap[classes[i]]; ok {
			currentMap = currentMap[classes[i]].(map[string]interface{})
		} else {
			if i == len(classes)-1 {
				currentMap[classes[i]] = map[string]interface{}{
					"title": note.Title,
					"uuid":  note.UUID,
				}
			} else {
				currentMap[classes[i]] = make(map[string]interface{})
				currentMap = currentMap[classes[i]].(map[string]interface{})
			}
		}
	}
	return ret
}

func main() {
	a := "wangyiwei"
	fmt.Println(a[:strings.LastIndex(a, "e")])
	//a := models.Note{
	//	Title:    "1.数组",
	//	RootPath: "/data/wangyiwei",
	//	NotePath: "/data/wangyiwei/c++/c++基础/bilibili",
	//	UUID:     "111",
	//}
	//b := convertNotesToObj(&a)
	//fmt.Println(b)
	//watcher, err := fsnotify.NewWatcher() //1. 先new
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer watcher.Close() //记得释放
	//done := make(chan bool)
	//go func() {
	//	for {
	//		select {
	//		case event, ok := <-watcher.Events:
	//			if !ok {
	//				return
	//			}
	//			log.Println("event:", event)
	//			if event.Op&fsnotify.Write == fsnotify.Write {
	//				log.Println("modified file:", event.Name)
	//			}
	//			if event.Op&fsnotify.Write == fsnotify.Create {
	//				log.Println("create file:", event.Name)
	//			}
	//			if event.Op&fsnotify.Write == fsnotify.Remove {
	//				log.Println("remove file:", event.Name)
	//			}
	//		case err, ok := <-watcher.Errors:
	//			if !ok {
	//				return
	//			}
	//			log.Println("error:", err)
	//		}
	//	}
	//}()
	//
	//err = watcher.Add("D:\\test")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//<-done
}
