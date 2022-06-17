package user

import (
	"errors"
	"tdp-cloud/core/dborm"
	"tdp-cloud/core/utils"

	"golang.org/x/crypto/bcrypt"
)

// 注册账号

type RegisterParam struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(param *RegisterParam) error {

	hash, _ := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	result := dborm.Db.Create(&dborm.User{Username: param.Username, Password: string(hash)})

	return result.Error

}

// 登录账号

type LoginParam struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResult struct {
	Keyid int    `json:"keyid"`
	Token string `json:"token"`
}

func Login(param *LoginParam) (LoginResult, error) {

	var res LoginResult

	var user dborm.User
	var secret dborm.Secret

	// 验证账号

	dborm.Db.First(&user, "username = ?", param.Username)

	if user.Id == 0 {
		return res, errors.New("账号错误")
	}

	// 验证密码

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password))

	if err != nil {
		return res, errors.New("密码错误")
	}

	// 创建令牌

	token := utils.RandString(32)
	dborm.Db.Create(&dborm.Session{UserId: user.Id, Token: token})

	// 获取密钥

	dborm.Db.First(&secret, "user_id = ?", user.Id)

	res.Keyid = secret.Id
	res.Token = token

	return res, nil

}

// 修改资料

type ModifyParam struct {
	UserId      int    `json:"userId"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	Description string `json:"description"`
}

func Modify(param *ModifyParam) error {

	var user dborm.User

	// 验证账号

	dborm.Db.First(&user, "user_id = ?", param.UserId)

	if user.Id == 0 {
		return errors.New("账号错误")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.OldPassword))

	if err != nil {
		return errors.New("密码错误")
	}

	// 更新资料

	if param.NewPassword != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(param.NewPassword), bcrypt.DefaultCost)
		user.Password = string(password)
	}

	user.Description = param.Description

	dborm.Db.Save(&user)

	return nil

}
