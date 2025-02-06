package api

import (
	myjwt "SecKill/middleware/jwt"
	"SecKill/model"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

import "SecKill/data"

const kindKey = "kind"

// user log in
func LoginAuth(ctx *gin.Context) {
	var postUser model.LoginUser
	// use the ctx JOSN data into model.LoginUser data structure,
	// if fail, return error
	if err := ctx.BindJSON(&postUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{kindKey: "", ErrMsgKey: "Parse JSON format fail."})
		return
	} else {
		// find the user
		queryUser := model.User{Username: postUser.Username}
		// get the first *DB that suits the queryUser
		err := data.Db.Where(&queryUser).First(&queryUser).Error
		if err != nil && gorm.IsRecordNotFoundError(err) {
			ctx.JSON(http.StatusUnauthorized, gin.H{kindKey: "", ErrMsgKey: "Username or Password is null."})
			return
		}

		// if the username is OK, then check the password
		if queryUser.Password != model.GetMD5(postUser.Password) {
			// Kind - customer / seller
			ctx.JSON(http.StatusUnauthorized, gin.H{kindKey: queryUser.Kind, ErrMsgKey: "Username or Password is null."})
		}

		// generate the token
		generateToken(ctx, queryUser)
	}
}

func generateToken(ctx *gin.Context, user model.User) {
	j := myjwt.NewJWT()
	claims := myjwt.CustomClaims{
		Username: user.Username,
		Password: user.Password,
		Kind:     user.Kind,
		StandardClaims: jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    myjwt.Issuer,                    //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			kindKey:   user.Kind,
			ErrMsgKey: err,
		})
		return
	}

	//log.Println(token)
	ctx.Header("Authorization", token)
	ctx.JSON(http.StatusOK, gin.H{
		kindKey:   user.Kind,
		ErrMsgKey: "",
	})
	return

}

// the user log out,
// 用户登出，TODO：修改退出的token，主要是删除掉redis里的session
func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("user")
	if err := session.Save(); err != nil {
		//log.Warningf(ctx, "Error when save deleted session. %v", err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{ErrMsgKey: "log out."})
	return
}
