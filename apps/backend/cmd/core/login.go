package main

import (
	"encoding/json"
	"fmt"
	xdb "maestro/internal/db"
	xhttp "maestro/internal/http"
	"net/http"

	"github.com/ayaviri/goutils/timer"
)

//  ____   ____ _   _ _____ __  __    _    ____
// / ___| / ___| | | | ____|  \/  |  / \  / ___|
// \___ \| |   | |_| |  _| | |\/| | / _ \ \___ \
//  ___) | |___|  _  | |___| |  | |/ ___ \ ___) |
// |____/ \____|_| |_|_____|_|  |_/_/   \_\____/
//

type UserLoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponseBody struct {
	BearerToken string `json:"bearer_token"`
}

//  _   _    _    _   _ ____  _     _____ ____  ____
// | | | |  / \  | \ | |  _ \| |   | ____|  _ \/ ___|
// | |_| | / _ \ |  \| | | | | |   |  _| | |_) \___ \
// |  _  |/ ___ \| |\  | |_| | |___| |___|  _ < ___) |
// |_| |_/_/   \_\_| \_|____/|_____|_____|_| \_\____/
//

func LoginResourceHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody UserLoginRequestBody

	timer.WithTimer("reading & unmarshaling request JSON body", func() {
		err = xhttp.ReadUnmarshalRequestBody(request, &requestBody)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not read body of request: %v\n", err),
			http.StatusInternalServerError,
		)
		return
	}

	var bearerToken string

	timer.WithTimer("verifying user credentials", func() {
		bearerToken, err = xdb.AuthenticateAndGenerateToken(
			db,
			requestBody.Username,
			requestBody.Password,
		)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not validate user credentials: %v\n", err),
			http.StatusInternalServerError,
		)
		return
	}

	var responseJson []byte
	responseJson, err = json.Marshal(UserLoginResponseBody{BearerToken: bearerToken})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not marshal response body to JSON: %v\n", err),
			http.StatusInternalServerError,
		)
		return
	}

	writer.Write(responseJson)
	// http.SetCookie(
	// 	writer,
	// 	&http.Cookie{
	// 		Name:     "bearer_token",
	// 		Value:    bearerToken,
	// 		MaxAge:   0,
	// 		Secure:   true,
	// 		HttpOnly: true,
	// 		SameSite: http.SameSiteStrictMode,
	// 	},
	// )
}
