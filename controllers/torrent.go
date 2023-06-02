package controllers

import (
	"fmt"
	"log"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gin-gonic/gin"
)

func DownloadTorrent(c *gin.Context) {
	client, _ := torrent.NewClient(nil)
	defer client.Close()
	info, err := metainfo.Load(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	t, _ := client.AddTorrent(info)
	<-t.GotInfo()
	fmt.Println(t.NumPieces())
	fmt.Println(t.Stats().PiecesComplete)
	t.DownloadAll()
	var a int
	for true {
		fmt.Scanln(&a)
		fmt.Printf("%d/%d\n", t.NumPieces(), t.Stats().PiecesComplete)
	}

	client.WaitAll()
	log.Print("ermahgerd, torrent downloaded")
}

// func DownloadTorrent(c *gin.Context) {

// 	// fileHeader, _ := c.FormFile("file")
// 	// // data, _ := ioutil.ReadAll(c.Request.Body)
// 	// file, _ := fileHeader.Open()
// 	// data, _ := ioutil.ReadAll(file)
// 	// fmt.Println("download", string(data))

// 	// // Create a session
// 	// ses, _ := torrent.NewSession(torrent.DefaultConfig)

// 	// // Add magnet link
// 	// tor, err := ses.AddTorrent(c.Request.Body, &torrent.AddTorrentOptions{
// 	// 	StopAfterDownload: true,
// 	// })
// 	// torrent.

// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	return
// 	// }

// 	// // Watch the progress
// 	// for range time.Tick(time.Second) {
// 	// 	s := tor.Stats()
// 	// 	log.Printf("Status: %s, Downloaded: %d, Peers: %d", s.Status.String(), s.Bytes.Completed, s.Peers.Total)
// 	// }
// }
