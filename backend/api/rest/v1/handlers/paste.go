package handlers

import (
	"net/http"
	"pasteGo/backend/api/rest/v1/types"
	"pasteGo/backend/db"
	"pasteGo/backend/db/typesDB"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetPaste(c *gin.Context) {
	pastePsw := types.PastePassword{}
	pasteId := c.Param("id")
	if err := c.BindJSON(&pastePsw); err != nil {
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

	paste, exist, err := DBInstance.GetPasteRecordById(pasteId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	if !exist {
		c.IndentedJSON(http.StatusNotFound, types.APIResponse{
			Code:        types.ErrPasteNotFound,
			Explanation: types.ErrPasteNotFoundExp,
		})
		return
	}
	if paste.Lifetime > 0 && paste.Lifetime < time.Now().Unix() {
		c.IndentedJSON(http.StatusNotFound, types.APIResponse{
			Code:        types.ErrPasteNotFound,
			Explanation: types.ErrPasteNotFoundExp,
		})
		deletePaste(paste.Id)
		return
	}

	userDB, exists, err := DBInstance.GetUserRecordById(paste.UserId)
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

	//Если вставка непубличная
	if !typesDB.IntToBool(paste.Public) {
		accessToken, err := c.Cookie(types.CookieAccessToken)
		response := types.APIResponse{
			Code:        types.ErrNotPublicPaste,
			Explanation: types.ErrNotPublicPasteExp,
		}
		if err != nil {
			DumpCookies(c)
			c.IndentedJSON(http.StatusUnauthorized, response)
			return
		}

		claims, err := ParseClaims(accessToken)
		if err != nil {
			DumpCookies(c)
			c.IndentedJSON(http.StatusUnauthorized, response)
			return
		}
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			DumpCookies(c)
			c.IndentedJSON(http.StatusUnauthorized, response)
			return
		}
	}

	//Первый запрос: на вставке имеется пароль, но пользователь не знает об этом
	if paste.Password != "" && pastePsw.Password == "" {
		c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
			Code:        types.ErrPasswordPaste,
			Explanation: types.ErrPasswordPasteExp,
		})
		return
	} else if paste.Password != "" && pastePsw.Password != "" { //Второй запрос: указан пароль
		if paste.Password != ShaHashing(pastePsw.Password) {
			c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
				Code:        types.ErrWrongPasswordPaste,
				Explanation: types.ErrWrongPasswordPasteExp,
			})
			return
		}
	}

	respPaste := types.Paste{
		Id:          "",
		Author:      userDB.Username,
		Created:     paste.Created,
		Updated:     paste.Updated,
		ExpTime:     0,
		Lifetime:    "",
		Text:        paste.Text,
		Password:    "",
		HasPassword: false,
		Public:      typesDB.IntToBool(paste.Public),
	}

	c.IndentedJSON(http.StatusOK, types.APIResponse{
		Code:        types.OperationSuccess,
		Explanation: types.OperationSuccessExp,
		Message:     respPaste,
	})
}

func GetPasteList(c *gin.Context) {
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

	pasteList, err := DBInstance.GetPasteRecordsByUserId(userDB.Id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	if pasteList == nil {
		c.IndentedJSON(http.StatusNotFound, types.APIResponse{
			Code:        types.ErrPasteNotFound,
			Explanation: types.ErrPasteNotFoundExp,
		})
		return
	}
	finalPasteList := make([]types.Paste, 0, len(*pasteList))
	var PastePassword bool
	for i := range *pasteList {
		if (*pasteList)[i].Lifetime > 0 && (*pasteList)[i].Lifetime < time.Now().Unix() {
			deletePaste((*pasteList)[i].Id)
			continue
		}
		if (*pasteList)[i].Password != "" {
			PastePassword = true
		} else {
			PastePassword = false
		}
		finalPasteList = append(finalPasteList, types.Paste{
			Id:          (*pasteList)[i].Id,
			Author:      claims.Subject,
			Created:     (*pasteList)[i].Created,
			Updated:     (*pasteList)[i].Updated,
			ExpTime:     (*pasteList)[i].Lifetime,
			Text:        (*pasteList)[i].Text,
			Password:    "",
			HasPassword: PastePassword,
			Public:      typesDB.IntToBool((*pasteList)[i].Public),
		})
	}

	c.IndentedJSON(http.StatusCreated, types.APIResponse{
		Code:        types.OperationSuccess,
		Explanation: types.OperationSuccessExp,
		Message: types.PasteList{
			Pastes: finalPasteList,
		},
	})
}

func CreatePaste(c *gin.Context) {
	paste := types.Paste{}
	if err := c.BindJSON(&paste); err != nil {
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

	timeNow := time.Now()
	var expires int64 = -1
	switch paste.Lifetime {
	case "minute":
		expires = timeNow.Add(time.Minute).Unix()
	case "hour":
		expires = timeNow.Add(time.Hour).Unix()
	case "day":
		expires = timeNow.Add(time.Hour * 24).Unix()
	case "week":
		expires = timeNow.Add(time.Hour * 24 * 7).Unix()
	case "month":
		expires = timeNow.Add(time.Hour * 24 * 30).Unix()
	case "year":
		expires = timeNow.Add(time.Hour * 24 * 365).Unix()
	}

	if paste.HasPassword && paste.Password != "" {
		paste.Password = ShaHashing(paste.Password)
	} else {
		paste.Password = ""
	}

	pasteRecord := typesDB.PasteRecord{
		Id:       uuid.New().String(),
		UserId:   userDB.Id,
		Text:     paste.Text,
		Created:  timeNow.Unix(),
		Updated:  -1,
		Lifetime: expires,
		Password: paste.Password,
		Public:   typesDB.BoolToInt(paste.Public),
	}

	created, err := DBInstance.AddPasteRecord(&pasteRecord)
	if err != nil || !created {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, types.APIResponse{
		Code:        types.OperationSuccess,
		Explanation: types.OperationSuccessExp,
		Message: types.Paste{
			Id:          pasteRecord.Id,
			Author:      claims.Subject,
			Created:     pasteRecord.Created,
			Updated:     pasteRecord.Updated,
			ExpTime:     pasteRecord.Lifetime,
			Lifetime:    paste.Lifetime,
			Text:        pasteRecord.Text,
			Password:    "",
			HasPassword: paste.HasPassword,
			Public:      paste.Public,
		},
	})
}

func UpdatePaste(c *gin.Context) {
	paste := types.Paste{}
	pasteId := c.Param("id")
	if err := c.BindJSON(&paste); err != nil {
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

	oldPasteRecord, exists, err := DBInstance.GetPasteRecordById(pasteId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}
	if !exists {
		c.IndentedJSON(http.StatusNotFound, types.APIResponse{
			Code:        types.ErrPasteNotFound,
			Explanation: types.ErrPasteNotFoundExp,
		})
		return
	}
	if oldPasteRecord.Lifetime > 0 && oldPasteRecord.Lifetime < time.Now().Unix() {
		c.IndentedJSON(http.StatusNotFound, types.APIResponse{
			Code:        types.ErrPasteNotFound,
			Explanation: types.ErrPasteNotFoundExp,
		})
		deletePaste(oldPasteRecord.Id)
		return
	}

	timeNow := time.Now()
	var expires int64 = -1
	switch paste.Lifetime {
	case "minute":
		expires = timeNow.Add(time.Minute).Unix()
	case "hour":
		expires = timeNow.Add(time.Hour).Unix()
	case "day":
		expires = timeNow.Add(time.Hour * 24).Unix()
	case "week":
		expires = timeNow.Add(time.Hour * 24 * 7).Unix()
	case "month":
		expires = timeNow.Add(time.Hour * 24 * 30).Unix()
	case "year":
		expires = timeNow.Add(time.Hour * 24 * 365).Unix()
	}

	if paste.HasPassword && paste.Password != "" {
		paste.Password = ShaHashing(paste.Password)
	} else if paste.HasPassword && paste.Password == "" {
		paste.Password = oldPasteRecord.Password
		if oldPasteRecord.Password == "" {
			paste.HasPassword = false
		} else {
			paste.HasPassword = true
		}
	} else {
		paste.Password = ""
		paste.HasPassword = false
	}

	newPasteRecord := typesDB.PasteRecord{
		Id:       oldPasteRecord.Id,
		UserId:   userDB.Id,
		Text:     paste.Text,
		Created:  oldPasteRecord.Created,
		Updated:  timeNow.Unix(),
		Lifetime: expires,
		Password: paste.Password,
		Public:   typesDB.BoolToInt(paste.Public),
	}

	err = DBInstance.EditPasteRecord(&newPasteRecord)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, types.APIResponse{
			Code:        types.ErrServer,
			Explanation: types.ErrServerExp,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, types.APIResponse{
		Code:        types.OperationSuccess,
		Explanation: types.OperationSuccessExp,
		Message: types.Paste{
			Id:          newPasteRecord.Id,
			Author:      claims.Subject,
			Created:     newPasteRecord.Created,
			Updated:     newPasteRecord.Updated,
			ExpTime:     newPasteRecord.Lifetime,
			Lifetime:    paste.Lifetime,
			Text:        newPasteRecord.Text,
			Password:    "",
			HasPassword: paste.HasPassword,
			Public:      paste.Public,
		},
	})
}

func DeletePaste(c *gin.Context) {
	pasteId := c.Param("id")

	if err := deletePaste(pasteId); err != nil {
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

func deletePaste(id string) error {
	DBInstance, err := db.GetDBInstance()
	if err != nil {
		return err
	}
	return DBInstance.DeleteRecord(id, typesDB.PastesTable)
}
