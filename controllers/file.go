package controllers

import (
	"fmt"
	"net/http"
	"voidcloud-server/internal/exception"
	"voidcloud-server/internal/storage"
	"voidcloud-server/internal/user"
	"voidcloud-server/services/permission"

	"github.com/gin-gonic/gin"
)

type FolderData struct {
	Permission *storage.Permission `json:"permission"`
	IsOwner    bool                `json:"is_owner"`

	// READ_ONLY
	Folder     *storage.Folder   `json:"folder"`
	Files      []string          `json:"files"`
	SubFolders []*storage.Folder `json:"subfolders"`
	Parents    []*storage.Folder `json:"parents"`

	// ALL
	PublicPermission *storage.Permission `json:"public_permission,omitempty"`
	SharePermission  *storage.Permission `json:"share_permission,omitempty"`
	ShareUsers       []string            `json:"share_users,omitempty"`
}

func GetFile(c *gin.Context) {
}

func GetFolder(c *gin.Context) {
	_folder, _ := c.Get("folder")
	folder := _folder.(*storage.Folder)

	_folderPermission, _ := c.Get("permission")
	folderPermission := _folderPermission.(*permission.FolderPermission)

	data := &FolderData{
		Permission: folderPermission.Permission,
		IsOwner:    folderPermission.IsOwner,
	}

	// if folderPermission.Permission >= storage.READ_ONLY {
	// 	data.Folder = folder
	// 	data.Files = folder.Files
	// 	data.SubFolders = folder.SubFolders

	// 	if _, requestParaent := c.GetQuery("parent"); requestParaent {
	// 		parents := make([]*storage.Folder, 0, folderPermission.TopDepth)
	// 		if folderPermission.TopDepth != 0 {
	// 			curFolder := folder
	// 			for curFolder != folderPermission.TopReadableFolder {
	// 				curFolder = curFolder.ParentFolder
	// 				parents = append(parents, curFolder)
	// 			}
	// 		}
	// 		data.Parents = parents
	// 	}
	// }

	if folderPermission.Permission == storage.ALL {
		data.PublicPermission = folder.PublicPermission
		data.ShareUsers = folder.ShareUsers
		data.SharePermission = folder.SharePermission
	}

	c.JSON(http.StatusOK, data)
}

func GetFileShareInfo(c *gin.Context) {

}

func ShareFile(c *gin.Context) {
	_user, _ := c.Get("user")

	var data struct {
		Folder           []string           `json:"folder"`
		Name             string             `json:"name"`
		IsDir            bool               `json:"isDir"`
		PublicPermission storage.Permission `json:"public"`
		SharePermission  storage.Permission `json:"share"`
	}

	c.BindJSON(&data)
	fmt.Println(data)

	if len(data.Folder) == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": exception.UNKNOWN,
		})
		return
	}

	var rootFolder *storage.Folder

	if data.Folder[0] == "drive" {
		rootFolder = _user.(*user.User).RootFolder
		data.Folder = data.Folder[1:]
	} else if true {
		// TODO: Share Folder
		fmt.Println("TODO: Share Folder")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": exception.UNKNOWN,
		})
		return
	}

	folder, err := rootFolder.FindFolder(data.Folder...)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": exception.UNKNOWN,
		})
	}

	if data.IsDir {
		folder, err = folder.FindFolder(data.Name)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": exception.UNKNOWN,
			})
			return
		}

		folder.PublicPermission = storage.READ_ONLY
	}

	c.JSON(http.StatusOK, gin.H{
		"folder": folder.Id,
	})
}

func DeleteFile(c *gin.Context) {
	var data struct {
		Folder []string `json:"folder"`
		Name   string   `json:"name"`
		IsDir  bool     `json:"isDir"`
	}

	c.BindJSON(&data)

	fmt.Println(data)
}
