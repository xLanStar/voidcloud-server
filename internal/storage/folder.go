package storage

import (
	"bytes"
	"fmt"
	"os"
)

type Folder struct {
	Id               string      `json:"id"`
	Name             string      `json:"name"`
	Files            []string    `json:"-"`
	SubFolders       []*Folder   `json:"-"`
	PublicPermission *Permission `json:"-"`
	SharePermission  *Permission `json:"-"`
	ShareUsers       []string    `json:"-"`
	ParentFolder     *Folder     `json:"-"`
}

func (folder *Folder) String() string {
	return fmt.Sprintf("資料夾{ ID=%v 名字=%v }", folder.Id, folder.Name)
}

func (folder *Folder) newSubFolder(name string) *Folder {
	newFolder := newFolder(name)
	newFolder.ParentFolder = folder
	folder.SubFolders = append(folder.SubFolders, newFolder)
	return newFolder
}

func (folder *Folder) GetAbsoluteFilePath() string {
	current := folder
	folders := []string{}

	for current != nil {
		folders = append(folders, current.Name)
		current = current.ParentFolder
	}

	var buffer bytes.Buffer
	for i := len(folders) - 1; i >= 0; i-- {
		buffer.WriteRune('/')
		buffer.WriteString(folders[i])
	}

	return buffer.String()
}

func (folder *Folder) CreateFolder(name ...string) (*Folder, error) {
	fmt.Print("CreateFolder:", name)
	var buffer bytes.Buffer
	buffer.WriteString(STORAGE_FOLDER)
	buffer.WriteString(folder.GetAbsoluteFilePath())

	for _, folderName := range name {
		fmt.Println("Folder:", folderName)
		buffer.WriteRune('/')
		buffer.WriteString(folderName)

		filepath := buffer.String()
		if _, err := os.Stat(filepath); err == nil {
			// folder exists
			fmt.Println("Folder exists")
			for _, subFolder := range folder.SubFolders {
				if subFolder.Name == folderName {
					folder = subFolder
					fmt.Println("use subFolder", subFolder)
					continue
				}
			}
		} else {
			// folder not exists
			fmt.Println("Folder not exists")
			err := os.Mkdir(filepath, 0644)
			if err != nil {
				fmt.Println("MkDir Error", err)
				return nil, err
			}
		}

		folder = folder.newSubFolder(folderName)
	}

	return folder, nil
}

func (folder *Folder) FindFolder(name ...string) (*Folder, error) {
	fmt.Println("FindFolder:", name)
	var buffer bytes.Buffer
	buffer.WriteString(folder.GetAbsoluteFilePath())

	for _, folderName := range name {
		fmt.Println("Folder:", folderName)
		buffer.WriteRune('/')
		buffer.WriteString(folderName)

		found := false
		for _, subFolder := range folder.SubFolders {
			fmt.Println("Find SubFolder:", subFolder.Name)
			if subFolder.Name == folderName {
				folder = subFolder
				found = true
				break
			}
		}

		if found {
			continue
		}

		filepath := buffer.String()
		if _, err := os.Stat(filepath); err != nil {
			folder = folder.newSubFolder(folderName)
			continue
		}

		return nil, fmt.Errorf("不存在此路徑")
	}
	return folder, nil
}

func newFolder(name string) *Folder {
	fmt.Println("newFolder:", name)
	folder := &Folder{
		Id:               generateFolderId(),
		Name:             name,
		Files:            make([]string, 0),
		SubFolders:       make([]*Folder, 0),
		ShareUsers:       make([]string, 0),
		PublicPermission: NO_PERMISSION,
		SharePermission:  NO_PERMISSION,
	}
	folderIdMap[folder.Id] = folder
	return folder
}
