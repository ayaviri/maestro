package main

import (
	"errors"
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

type UserRegistrationRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

//  _   _    _    _   _ ____  _     _____ ____  ____
// | | | |  / \  | \ | |  _ \| |   | ____|  _ \/ ___|
// | |_| | / _ \ |  \| | | | | |   |  _| | |_) \___ \
// |  _  |/ ___ \| |\  | |_| | |___| |___|  _ < ___) |
// |_| |_/_/   \_\_| \_|____/|_____|_____|_| \_\____/
//

func RegistrationResourceHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody UserRegistrationRequestBody

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

	timer.WithTimer("creating user and user cart", func() {
		var usernameAvailable bool
		usernameAvailable, err = xdb.IsUsernameAvailable(db, requestBody.Username)

		if err != nil {
			return
		}

		if !usernameAvailable {
			err = errors.New("Username already in use")
			return
		}

		_, err = xdb.CreateUser(
			db,
			requestBody.Username,
			requestBody.Password,
			requestBody.Email,
		)
	})

	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Could not create user: %v\n", err),
			http.StatusInternalServerError,
		)
		return
	}
}
