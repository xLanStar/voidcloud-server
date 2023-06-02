package storage

import (
	fastio "github.com/xLanStar/go-fast-io"
)

type folderWriter struct {
	*fastio.FileWriter
}

var FolderWriter folderWriter

func (folderWriter) New() folderWriter {
	folderWriter := folderWriter{&fastio.FileWriter{}}
	folderWriter.Init()
	return folderWriter
}

func (folderWriter *folderWriter) WriteFolder(folder *Folder) {
	folderWriter.WriteString(folder.Id)
	folderWriter.WriteString(folder.Name)

	share := (folder.SharePermission.Index << 4) + folder.PublicPermission.Index

	folderWriter.WriteUint8(uint8(share))
	folderWriter.WriteUint8(uint8(len(folder.ShareUsers)))

	for _, userId := range folder.ShareUsers {
		folderWriter.WriteString(userId)
	}

	folderWriter.WriteUint32(uint32(len(folder.SubFolders)))

	for _, subFolder := range folder.SubFolders {
		folderWriter.WriteFolder(subFolder)
	}
}
