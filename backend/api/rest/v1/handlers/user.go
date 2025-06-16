package handlers

import (
	"net/http"
	"pasteGo/backend/api/rest/v1/types"
	"pasteGo/backend/db"
	"pasteGo/backend/db/typesDB"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func UpdateUser(c *gin.Context) {
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

	userDB, exists, err := DBInstance.GetUserRecordByUsername(claims.Subject)
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

	_, exists, err = DBInstance.GetUserRecordByUsername(user.Username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	if exists {
		c.IndentedJSON(http.StatusConflict, types.APIResponse{
			Code:        types.ErrExistUser,
			Explanation: types.ErrExistUserExp,
		})
		return
	}

	newUser := typesDB.UserRecord{
		Id:       userDB.Id,
		Username: userDB.Username,
		Password: userDB.Password,
	}
	if user.Password == "" && user.Username == "" {
		c.IndentedJSON(http.StatusConflict, types.APIResponse{
			Code:        types.ErrUserEmptyCredentials,
			Explanation: types.ErrUserEmptyCredentialsExp,
		})
		return
	}
	if user.Password != "" {
		newUser.Password = ShaHashing(user.Password)
	}
	if user.Username != "" {
		newUser.Username = user.Username
	}
	if user.Username == userDB.Username || ShaHashing(user.Password) == userDB.Password {
		c.IndentedJSON(http.StatusConflict, types.APIResponse{
			Code:        types.ErrUserSameCredentials,
			Explanation: types.ErrUserSameCredentialsExp,
		})
		return
	}

	err = DBInstance.EditUserRecord(&newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
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

func DeleteUser(c *gin.Context) {
	DumpCookies(c)
	DBInstance, err := db.GetDBInstance()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}

	rawClaims, exists := c.Get("userClaims")
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
			Code:        types.ErrGetCookies,
			Explanation: types.ErrGetCookiesExp,
		})
		return
	}
	claims, ok := rawClaims.(*jwt.RegisteredClaims)
	if !ok {
		c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
			Code:        types.ErrJWTProcessing,
			Explanation: types.ErrJWTProcessingExp,
		})
		return
	}

	userDB, exists, err := DBInstance.GetUserRecordByUsername(claims.Subject)
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

	if err := DBInstance.DeleteRecord(userDB.Id, typesDB.UsersTable); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, types.APIResponse{
		Code:        types.OperationSuccess,
		Explanation: types.OperationSuccessExp,
	})
}
