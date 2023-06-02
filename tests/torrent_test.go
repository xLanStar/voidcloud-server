package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestDownloadTorrent(t *testing.T) {
	file, _ := os.OpenFile("./test.torrent", os.O_RDONLY, 0644)

	data, _ := ioutil.ReadAll(file)

	fmt.Println(string(data))

	responseBody := bytes.NewBuffer(data)

	http.Post("http://voidcloud.net/api/torrent/download", "text/plain", responseBody)
}
