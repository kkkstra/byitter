package jwt

import (
	"byitter/src/controller/response"
	"byitter/src/model"
	"byitter/src/util/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		//fmt.Println(bearerToken)
		bearerTokenSplit := strings.Split(bearerToken, " ")
		if len(bearerTokenSplit) < 2 {
			response.Error(c, http.StatusBadRequest, "Empty token! ")
			return
		}
		token := jwt.TokenStr(bearerTokenSplit[1])
		if token == "" {
			response.Error(c, http.StatusBadRequest, "Empty token! ")
			return
		}
		claims, err := token.ParseToken(&jwt.UserClaims{})
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid token! ", err.Error())
			return
		} else if time.Now().Unix() > claims.Claims.(*jwt.UserClaims).ExpiresAt {
			response.Error(c, http.StatusBadRequest, "Expired token! ")
			return
		}

		// 检查token是否因更改密码而失效
		id := claims.Claims.(*jwt.UserClaims).Id
		m := model.GetModel(&model.UserModel{})
		u, err := m.(*model.UserModel).FindUserById(id)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "User does not exist. ", err.Error())
			return
		}
		if u.UpdatePasswordAt.Unix() > claims.Claims.(*jwt.UserClaims).IssuedAt {
			response.Error(c, http.StatusBadRequest, "Invalid toknen. ")
			return
		}

		userData := map[string]string{
			"id":       claims.Claims.(*jwt.UserClaims).Id,
			"username": claims.Claims.(*jwt.UserClaims).Username,
			"role":     claims.Claims.(*jwt.UserClaims).Role.Str(),
		}
		c.Set("userdata", userData)
		c.Next()
	}
}
