package permission

import (
	"fmt"
	"voidcloud-server/internal/storage"
	"voidcloud-server/internal/user"
)

type FolderPermission struct {
	TopReadableFolder *storage.Folder
	TopDepth          uint16
	IsOwner           bool
	Permission        *storage.Permission
	Folder            *storage.Folder
}

// 必須有
func _GetFolderPermission(folder *storage.Folder, user *user.User) *FolderPermission {

	// 檢查權限
	topReadableFolder := folder
	topDepth := uint16(0)
	permission := storage.NO_PERMISSION

	// 使用者相關設定
	isOwner := false
	if user != nil {
		curFolder := folder
		depth := uint16(0)
		for {
			for _, shareUser := range curFolder.ShareUsers {
				if shareUser == user.Account {
					if curFolder.SharePermission.Index > permission.Index {
						permission = curFolder.SharePermission
					}
					if curFolder.SharePermission.Index >= storage.READ_ONLY.Index {
						topReadableFolder = curFolder
						topDepth = depth
					}
					break
				}
			}
			if curFolder.ParentFolder == nil || permission == storage.ALL {
				break
			}
			depth++
			curFolder = curFolder.ParentFolder
		}

		if curFolder == user.RootFolder {
			// 擁有者
			permission = storage.ALL
			isOwner = true
		}
	}

	if permission != storage.ALL {
		// 公開設定
		curFolder := folder
		depth := uint16(0)
		for {
			if curFolder.PublicPermission.Index >= storage.READ_ONLY.Index && depth > topDepth {
				topReadableFolder = curFolder
				topDepth = depth
			}
			if curFolder.PublicPermission.Index > permission.Index {
				permission = curFolder.PublicPermission
				if permission == storage.ALL {
					break
				}
			}
			if curFolder.ParentFolder == nil {
				break
			}
			depth++
			curFolder = curFolder.ParentFolder
		}
	}

	// c.Set("permission", FolderPermission{
	// 	TopReadableFolder: topReadableFolder,
	// 	TopDepth:          topDepth,
	// 	IsOwner:           isOwner,
	// 	Permission:        permission,
	// 	Folder:            folder,
	// })
	return &FolderPermission{
		TopReadableFolder: topReadableFolder,
		TopDepth:          topDepth,
		IsOwner:           isOwner,
		Permission:        permission,
	}
}

// 必須有
func GetFolderPermission(folder *storage.Folder, user *user.User) *storage.Permission {

	// 檢查權限
	permission := storage.NO_PERMISSION
	fmt.Println(permission)

	// 使用者相關設定
	if user != nil {
		curFolder := folder
		for {
			// 公開
			if curFolder.PublicPermission.Index > permission.Index {
				permission = curFolder.PublicPermission
				fmt.Println(permission)
			}
			// 共享
			for _, shareUser := range curFolder.ShareUsers {
				if shareUser == user.Account {
					if curFolder.SharePermission.Index > permission.Index {
						permission = curFolder.SharePermission
						fmt.Println(permission)
					}
					break
				}
			}
			if curFolder.ParentFolder == nil || permission == storage.ALL {
				break
			}
			curFolder = curFolder.ParentFolder
		}

		if curFolder == user.RootFolder {
			// 擁有者
			permission = storage.ALL
		}
	}

	return permission
}
