package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"pasteGo/backend/api/rest/v1/types"
	"pasteGo/backend/db"
	"pasteGo/backend/db/typesDB"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	user := types.User{}
	if err := c.BindJSON(&user); err != nil {
		return
	}

	DBInstance, err := db.GetDBInstance()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	newUUID := uuid.New().String()
	userRecord := typesDB.UserRecord{
		Id:       newUUID,
		Username: user.Username,
		Password: ShaHashing(user.Password),
	}

	created, err := DBInstance.AddUserRecord(&userRecord)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	if !created {
		c.IndentedJSON(http.StatusConflict, types.APIResponse{
			Code:        types.ErrExistUser,
			Explanation: types.ErrExistUserExp,
		})
		return
	}

	newTokens, err := GenerateTokens(user.Username)
	if err != nil {
		c.IndentedJSON(http.StatusCreated, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	created, err = DBInstance.AddToken(&typesDB.TokenRecord{
		UserId:       newUUID,
		RefreshToken: newTokens.RefreshToken,
	})
	if err != nil || !created {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}

	setCookies(c, newTokens, user.Username)

	c.IndentedJSON(http.StatusCreated, types.APIResponse{
		Code:        types.OperationSuccess,
		Explanation: types.OperationSuccessExp,
	})
}

func Login(c *gin.Context) {
	user := types.User{}
	if err := c.BindJSON(&user); err != nil {
		return
	}

	DBInstance, err := db.GetDBInstance()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}

	userDB, exists, err := DBInstance.GetUserRecordByUsername(user.Username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	if !exists {
		c.IndentedJSON(http.StatusNotFound, types.APIResponse{
			Code:        types.ErrUserNotFound,
			Explanation: types.ErrUserNotFoundExp,
		})
		return
	}

	if userDB.Password != ShaHashing(user.Password) {
		c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
			Code:        types.ErrWrongCredentials,
			Explanation: types.ErrWrongCredentialsExp,
		})
		return
	}

	oldTokens, err := DBInstance.GetTokenByUserId(userDB.Id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	for i := range *oldTokens {
		DBInstance.DeleteToken((*oldTokens)[i].RefreshToken)
	}

	newTokens, err := GenerateTokens(user.Username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	created, err := DBInstance.AddToken(&typesDB.TokenRecord{
		UserId:       userDB.Id,
		RefreshToken: newTokens.RefreshToken,
	})
	if err != nil || !created {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}

	setCookies(c, newTokens, user.Username)

	c.IndentedJSON(http.StatusOK, types.APIResponse{
		Code:        types.OperationSuccess,
		Explanation: types.OperationSuccessExp,
	})
}

func Logout(c *gin.Context) {
	DumpCookies(c)
	DBInstance, err := db.GetDBInstance()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}

	oldRefreshToken, err := c.Cookie(types.CookieRefreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
			Code:        types.ErrGetCookies,
			Explanation: types.ErrGetCookiesExp,
		})
		return
	}
	exists, err := DBInstance.CheckIfExistToken(oldRefreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	if exists {
		DBInstance.DeleteToken(oldRefreshToken)
	}

	c.IndentedJSON(http.StatusOK, types.APIResponse{
		Code:        types.OperationSuccess,
		Explanation: types.OperationSuccessExp,
	})
}

func Refresh(c *gin.Context) {
	DBInstance, err := db.GetDBInstance()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		DumpCookies(c)
		return
	}

	rawClaims, exists := c.Get("userClaims")
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
			Code:        types.ErrGetCookies,
			Explanation: types.ErrGetCookiesExp,
		})
		DumpCookies(c)
		return
	}

	claims, ok := rawClaims.(*jwt.RegisteredClaims)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
			Code:        types.ErrJWTProcessing,
			Explanation: types.ErrJWTProcessingExp,
		})
		DumpCookies(c)
		return
	}

	oldRefreshToken, err := c.Cookie(types.CookieRefreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
			Code:        types.ErrGetCookies,
			Explanation: types.ErrGetCookiesExp,
		})
		DumpCookies(c)
		return
	}

	exists, err = DBInstance.CheckIfExistToken(oldRefreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	//Если такой refresh-токен есть в БД
	if exists {
		//Создание нового токена
		newTokens, err := GenerateTokens(claims.Subject)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
				Code:        types.ErrJWTProcessing,
				Explanation: types.ErrJWTProcessingExp,
			})
			return
		}
		err = DBInstance.ChangeToken(&typesDB.TokenRecord{
			UserId:       claims.Subject,
			RefreshToken: newTokens.RefreshToken,
		}, oldRefreshToken)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
				Code:        types.ErrServer,
				Explanation: types.ErrServerExp,
			})
			return
		}

		setCookies(c, newTokens, claims.Subject)

		c.IndentedJSON(http.StatusOK, types.APIResponse{
			Code:        types.OperationSuccess,
			Explanation: types.OperationSuccessExp,
		})

	} else { //Если такого refresh-токена нет в БД
		DumpCookies(c)
		c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
			Code:        types.ErrJWTNotFound,
			Explanation: types.ErrJWTNotFoundExp,
		})
		return
	}
}

func setCookies(c *gin.Context, tokens types.Tokens, username string) {
	hour := 60 * 60 * 1
	week := 60 * 60 * 24 * 7
	var expTime int64 = time.Now().Add(time.Hour * 1).Unix()
	c.SetCookie(types.CookieAccessToken, tokens.AccessToken, hour, "/", "localhost", false, true)   // 1 hour
	c.SetCookie(types.CookieRefreshToken, tokens.RefreshToken, week, "/", "localhost", false, true) // week
	c.SetCookie(types.CookieUsername, username, hour, "/", "localhost", false, false)
	c.SetCookie(types.CookieExp, strconv.FormatInt(expTime, 10), week, "/", "localhost", false, false)
}

func DumpCookies(c *gin.Context) {
	c.SetCookie(types.CookieAccessToken, "", -1, "/", "localhost", false, true)
	c.SetCookie(types.CookieRefreshToken, "", -1, "/", "localhost", false, true)
	c.SetCookie(types.CookieUsername, "", -1, "/", "localhost", false, false)
	c.SetCookie(types.CookieExp, "", -1, "/", "localhost", false, false)
}

func GenerateTokens(username string) (types.Tokens, error) {
	claimsAccess := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), //Срок жизни
		IssuedAt:  jwt.NewNumericDate(time.Now()),                    //Время создания
		Subject:   username,
	}
	accessTokenString, err := generatejwt(claimsAccess)
	if err != nil {
		return types.Tokens{}, err
	}

	claimsRefresh := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   username,
	}
	refreshTokenString, err := generatejwt(claimsRefresh)
	if err != nil {
		return types.Tokens{}, err
	}

	return types.Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func generatejwt(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(types.SecretKey)
}

func ShaHashing(input string) string {
	plainText := []byte(input)
	sha256Hash := sha256.Sum256(plainText)
	return hex.EncodeToString(sha256Hash[:])
}

func ParseClaims(tokenString string) (*jwt.RegisteredClaims, error) {
	// Парсинг токена с проверкой подписи
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return types.SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Проверка валидности токена
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// Приведение claims к типу jwt.RegisteredClaims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("claims are not of type RegisteredClaims")
	}

	return claims, nil
}
