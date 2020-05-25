package v1

import (
	myjwt "ecode/middleware/jwt"
	"ecode/models"
	"ecode/utils"
	"ecode/utils/md5"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Login -
func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		utils.HandelError(c, utils.StatusBadMessage.Fail.Login)
		return
	}
	password = md5.Md5(password)
	user, err := models.Login(email, password)
	if err != nil {
		utils.HandelError(c, utils.StatusBadMessage.Fail.Login)
		return
	}
	generateToken(c, user)
	return
}

// Logout -
func Logout(c *gin.Context) {

}

// 生成令牌
func generateToken(c *gin.Context, user models.User) {
	j := myjwt.NewJWT()

	claims := myjwt.CustomClaims{
		ID:    strconv.Itoa(user.ID),
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),      // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600*24*7), // 过期时间 一星期
			Issuer:    "ecode",                              //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功",
		"data":  user,
		"token": token,
	})
	return
}
