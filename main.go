package main

import (
	"log"
	"os"
	"os/signal"
	"voidcloud-server/internal/router"
	"voidcloud-server/internal/storage"
	"voidcloud-server/internal/user"
	"voidcloud-server/server"
	"voidcloud-server/services/auth"

	_ "github.com/joho/godotenv/autoload"
)

var (
	webRouter     *router.WebRouter
	storageRouter *router.StorageRouter
)

func init() {
	// 設定為正式版
	// gin.SetMode(gin.ReleaseMode)

	auth.Init()
	storage.InitStorage()
	user.InitUser()

	// 建立路由
	webRouter = router.NewWebRouter()
	storageRouter = router.NewStorageRouter(os.Getenv("STORAGE_FOLDER"), os.Getenv("STORAGE_PREFIX"))
}

func main() {
	defer storage.SaveStorage()
	defer user.SaveUser()

	// DEBUG:
	// user.GetUserByAccount("ofei0305").RootFolder = storage.CreateRootFolder("歐菲系統廚具")
	// user.CreateUser("歐菲系統廚具", "ofei0305", "ofei0305@gmail.com", "48EmfQkusyMyC3Ur")
	// user.CreateUser("Lanstar", "Lanstar", "danny95624268@gmail.com", "aa95624268")

	// 啟用 Web HTTP 服務
	// 啟用 File Webdav 服務
	if os.Getenv("HTTP_ENABLE") == "true" {
		go server.ListenAndServe(":"+os.Getenv("STORAGE_PORT"), storageRouter)
		go server.ListenAndServe(":"+os.Getenv("WEB_PORT"), webRouter)
	}
	if os.Getenv("HTTPS_ENABLE") == "true" {
		go server.ListenAndServeTLS(":"+os.Getenv("STORAGE_PORT_TLS"), "./certs/cert.crt", "./certs/key.pem", storageRouter)
		go server.ListenAndServeTLS(":"+os.Getenv("WEB_PORT_TLS"), "./certs/cert.crt", "./certs/key.pem", webRouter)
	}

	// WEBDAV
	// http.ListenAndServe(":8080", &webdav.Handler{
	// 	// Prefix:     "",

	// 	FileSystem: webdav.Dir("."),
	// 	LockSystem: webdav.NewMemLS(),
	// })

	//
	var quit chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("關閉伺服器中...")
}
