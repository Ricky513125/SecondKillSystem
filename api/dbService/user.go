package dbService

import (
	"SecKill/data"
	"SecKill/model"
)

// use userName to get the model.User
func GetUser(userName string) (model.User, error) {
	user := model.User{}
	operation := data.Db.Where("username = ?", userName).First(&user)
	return user, operation.Error
}
