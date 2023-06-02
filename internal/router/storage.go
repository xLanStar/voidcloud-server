package router

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"voidcloud-server/internal/storage"
	"voidcloud-server/internal/user"
	"voidcloud-server/services/auth"
	"voidcloud-server/services/permission"

	"golang.org/x/net/webdav"
)

type StorageRouter struct {
	webdav.Handler
}

func NewStorageRouter(storageRoot, WebdavPrefix string) *StorageRouter {
	router := &StorageRouter{webdav.Handler{
		Prefix:     WebdavPrefix,
		FileSystem: webdav.Dir(storageRoot),
		LockSystem: webdav.NewMemLS(),
	}}
	return router
}

func (router StorageRouter) Init() {
}

const BaseAuthorizationPrefix = "Basic "

// ServeHTTP conforms to the http.Handler interface.
func (router *StorageRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[", req.Method, "] URL:", req.URL.Path)

	// CORS
	cors(w, req)

	// OPTIONS
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		// fmt.Println("回應:", http.StatusOK)
		return
	}

	// 取得資料夾
	folder, reqFolder, reqPath := getFolder(w, req)
	if folder != nil {
		// 指定資料夾
		// fmt.Println("Case 1:", folder, reqFolder, reqPath)
		router.routeShareFolder(w, req, folder, reqPath)
	} else if reqFolder == "drive" {
		// 自己的資料夾
		// fmt.Println("Case 2:", folder, reqFolder, reqPath)
		router.routeUserFolder(w, req, reqPath)
	} else {
		// 未知路徑
		// fmt.Println("Case 3:", folder, reqFolder, reqPath)
		w.WriteHeader(http.StatusForbidden)
		// fmt.Println("回應:", http.StatusForbidden)
		return
	}
}

func (router *StorageRouter) routeUserFolder(w http.ResponseWriter, req *http.Request, reqPath string) {
	// 身分驗證
	var user *user.User = verifyAuthorization(w, req)

	// DEBUG
	// if user != nil {
	// 	fmt.Printf("[StorageRouter] 已登入: %s\n", user)
	// }

	if user == nil {
		log.Println("Unuthorized: Need Auth or Need Permission")
		responseUnauthorization(w)
		return
	}

	// MOVE
	if req.Method == "MOVE" {
		userFolderPath := user.RootFolder.GetAbsoluteFilePath()
		destination, _ := url.QueryUnescape(req.Header["Destination"][0])
		destPath := destination[strings.Index(destination, "/drive")+6:]
		os.Rename(storage.STORAGE_FOLDER+userFolderPath+reqPath, storage.STORAGE_FOLDER+userFolderPath+destPath)
		w.WriteHeader(http.StatusOK)
	}

	req.URL.Path = "/" + user.Account + reqPath
	router.Handler.ServeHTTP(w, req)
}

func (router *StorageRouter) routeShareFolder(w http.ResponseWriter, req *http.Request, folder *storage.Folder, reqPath string) {
	// fmt.Println("public permission:", folder.PublicPermission, "share permission:", folder.SharePermission)

	if folder.PublicPermission.Can[req.Method] {
		// 指定分享的資料夾
		req.URL.Path = folder.GetAbsoluteFilePath() + reqPath
		// log.Println("重導至指定的資料夾:", folder.GetAbsoluteFilePath(), "最終URL:", req.URL.Path)
		// 以 Public 權限存取
		router.Handler.ServeHTTP(w, req)
		return
	}

	// 身分驗證
	var user *user.User = verifyAuthorization(w, req)

	// DEBUG
	// if user != nil {
	// 	fmt.Printf("[StorageRouter] 已登入: %s\n", user)
	// }

	if !permission.GetFolderPermission(folder, user).Can[req.Method] {
		// log.Println("Unuthorized: Need Auth or Need Permission")
		responseUnauthorization(w)
		return
	}
	router.Handler.ServeHTTP(w, req)
}

func cors(w http.ResponseWriter, req *http.Request) {
	if origin, ok := req.Header["Origin"]; ok {
		// fmt.Println("Has origin", origin, origin[0])
		w.Header().Set("Access-Control-Allow-Origin", origin[0])
		w.Header().Set("Access-Control-Allow-Headers", "Folder, Authorization, depth") //, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With
		w.Header().Set("Access-Control-Allow-Methods", "PROPFIND")
	}
}

func getFolder(w http.ResponseWriter, req *http.Request) (*storage.Folder, string, string) {
	if len(req.URL.Path) == 0 {
		return nil, "", ""
	}
	indexFolder := 1 + strings.IndexByte(req.URL.Path[1:], '/')
	if indexFolder == 0 {
		if req.URL.Path[len(req.URL.Path)-1] != '/' {
			rootFolder := req.URL.Path[1:]
			return storage.GetFolder(rootFolder), rootFolder, ""
		}
		return nil, "", ""
	}
	rootFolder := req.URL.Path[1:indexFolder]
	return storage.GetFolder(rootFolder), rootFolder, req.URL.Path[indexFolder:]
}

func verifyToken(w http.ResponseWriter, req *http.Request, token string) *user.User {
	if token == "" {
		return nil
	}
	authClaim, err := auth.Verify(token)
	if err != nil {
		log.Println("[StorageRouter] 無效的 Token")
		responseUnauthorization(w)
		return nil
	}
	return user.GetUserByAccount(authClaim.Account)
}

func verifyAuthorization(w http.ResponseWriter, req *http.Request) *user.User {
	var _user *user.User
	if authorization, hasAuth := req.Header["Authorization"]; hasAuth {
		if strings.HasPrefix(authorization[0], BaseAuthorizationPrefix) {
			// Case 1: 以帳號、密碼登入
			account, password, authOK := req.BasicAuth()
			if !authOK {
				log.Println("[StorageRouter] 無效的 Authorization")
				responseUnauthorization(w)
				return nil
			}

			_user = user.GetUserByAccount(account)
			if _user == nil || _user.Password != password {
				log.Println("[StorageRouter] 身分驗證發生錯誤")
				responseUnauthorization(w)
				return nil
			}

		} else {
			// Case 2: Header - Token
			_user = verifyToken(w, req, authorization[0])
			if _user == nil {
				responseUnauthorization(w)
				return nil
			}
		}
	} else {
		// Case 3: Query - Token
		query := req.URL.Query()
		if query.Has("Authorization") {
			_user = verifyToken(w, req, query.Get("Authorization"))
			if _user == nil {
				responseUnauthorization(w)
				return nil
			}
		}
	}
	return _user
}

func responseUnauthorization(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Println("回應:", http.StatusUnauthorized)
}
