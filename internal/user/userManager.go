package user

import (
	"fmt"
	"os"
	"voidcloud-server/internal/storage"
	"voidcloud-server/internal/util"
)

var (
	userIdMap  map[string]*User = make(map[string]*User)
	accountMap map[string]*User = make(map[string]*User)
	USER_DATA  string
)

func InitUser() {
	fmt.Println("[UserManager] 初始化")
	USER_DATA = os.Getenv("USER_DATA")

	if _, err := os.Stat(USER_DATA); err != nil {
		os.Create(USER_DATA)
		return
	}

	userReader := UserReader.New()
	userReader.OpenFile(USER_DATA, os.O_RDONLY, 0666)
	for userReader.Available() {
		user := userReader.ReadUser()

		if user == nil {
			continue
		}

		userIdMap[user.Id] = user
		accountMap[user.Account] = user

		fmt.Printf("[UserManager] 讀取使用者 %s\n", user)
	}
	userReader.Close()
}

func SaveUser() {
	fmt.Println("[UserManager] 保存")
	userWriter := UserWriter.New()
	err := userWriter.OpenFile(USER_DATA, os.O_CREATE, 0666)

	if err != nil {
		return
	}

	for _, user := range userIdMap {
		userWriter.WriteUser(user)
		fmt.Printf("[UserManager] 保存使用者 %s\n", user)
	}
	userWriter.Close()
}

func GenerateUserId() string {
	id := util.NewUUID()
	_, found := userIdMap[id]
	for found {
		id = util.NewUUID()
		_, found = userIdMap[id]
	}
	return id
}

func CreateUser(name, account, email, password string) *User {
	user := &User{
		Id:       GenerateUserId(),
		Name:     name,
		Account:  account,
		Email:    email,
		Password: password,
	}

	// 建立使用者雲端空間
	user.RootFolder = storage.CreateRootFolder(user.Account)

	// 註冊 ID
	userIdMap[user.Id] = user
	accountMap[user.Account] = user

	fmt.Printf("[UserManager] 建立使用者 %s\n", user)

	return user
}

func GetUserById(id string) *User {
	return userIdMap[id]
}

func GetUserByAccount(account string) *User {
	return accountMap[account]
}
