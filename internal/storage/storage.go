package storage

import (
	"fmt"
	"os"
	"voidcloud-server/internal/util"
)

var (
	folderIdMap map[string]*Folder = make(map[string]*Folder)
	rootFolders []*Folder          = make([]*Folder, 0)

	STORAGE_FOLDER string
	STORAGE_DATA   string
)

func InitStorage() {
	fmt.Println("[StorageManager] 初始化")
	STORAGE_FOLDER = os.Getenv("STORAGE_FOLDER")
	STORAGE_DATA = os.Getenv("STORAGE_DATA")

	if _, err := os.Stat(STORAGE_DATA); err != nil {
		return
	}

	folderReader := FolderReader.New()
	folderReader.OpenFile(STORAGE_DATA, os.O_RDONLY, 0666)
	for folderReader.Available() {
		folder, err := folderReader.ReadFolder(STORAGE_FOLDER)

		if err != nil {
			continue
		}

		rootFolders = append(rootFolders, folder)

		fmt.Printf("[StorageManager] 讀取資料夾 %s\n", folder)
	}
	folderReader.Close()
}

func SaveStorage() {
	fmt.Println("[StorageManager] 保存")
	folderWriter := FolderWriter.New()
	err := folderWriter.OpenFile(STORAGE_DATA, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, folder := range rootFolders {
		folderWriter.WriteFolder(folder)
		fmt.Printf("[StorageManager] 保存資料夾 %s\n", folder)
	}
	folderWriter.Close()
}

func generateFolderId() string {
	id := util.NewUUID()
	_, found := folderIdMap[id]
	for found {
		id = util.NewUUID()
		_, found = folderIdMap[id]
	}
	return id
}

func GetFolder(id string) *Folder {
	return folderIdMap[id]
}

func CreateRootFolder(name string) *Folder {
	folder := newFolder(name)

	rootFolders = append(rootFolders, folder)

	err := os.Mkdir(STORAGE_FOLDER+"/"+folder.Name, 0666)

	if err != nil {
		fmt.Printf("[StorageManager] 資料夾 Id:%s Name:%s 建立資料夾失敗\n", folder.Id, folder.Name)
		return folder
	}

	return folder
}
