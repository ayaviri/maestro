package main

import (
	"encoding/json"
	"maestro/internal"
	xdb "maestro/internal/db"
	xhttp "maestro/internal/http"
	xyoutube "maestro/internal/youtube"
	"net/http"
	"net/url"
)

//  ____   ___  _   _ _____ _____ ____
// |  _ \ / _ \| | | |_   _| ____|  _ \
// | |_) | | | | | | | | | |  _| | |_) |
// |  _ <| |_| | |_| | | | | |___|  _ <
// |_| \_\\___/ \___/  |_| |_____|_| \_\
//

func cartResourceHandler(writer http.ResponseWriter, request *http.Request) {
	var user xdb.User

	internal.WithTimer("getting user from auth bearer token", func() {
		var bearerToken string
		bearerToken, _ = xhttp.GetAuthBearerToken(request)
		user, err = xdb.GetUserFromToken(db, bearerToken)
	})

	if err != nil {
		http.Error(
			writer,
			"Could not get user from bearer token",
			http.StatusInternalServerError,
		)
		return
	}

	var cartId int64

	internal.WithTimer("getting user cart", func() {
		cartId, err = xdb.GetUserCartId(db, user.Id)
	})

	if err != nil {
		http.Error(
			writer,
			"Could not get/create cart for user",
			http.StatusInternalServerError,
		)
		return
	}

	switch request.Method {
	case http.MethodGet:
		getCartHandler(writer, request, user, cartId)
	case http.MethodPost:
		addToCartHandler(writer, request, user, cartId)
	case http.MethodDelete:
		removeFromCartHandler(writer, request, user, cartId)
	}
}

//  ____   ____ _   _ _____ __  __    _    ____
// / ___| / ___| | | | ____|  \/  |  / \  / ___|
// \___ \| |   | |_| |  _| | |\/| | / _ \ \___ \
//  ___) | |___|  _  | |___| |  | |/ ___ \ ___) |
// |____/ \____|_| |_|_____|_|  |_/_/   \_\____/
//

type GetCartResponseBody struct {
	CartItems []xyoutube.Video `json:"cart_items"`
}

//  _   _    _    _   _ ____  _     _____ ____  ____
// | | | |  / \  | \ | |  _ \| |   | ____|  _ \/ ___|
// | |_| | / _ \ |  \| | | | | |   |  _| | |_) \___ \
// |  _  |/ ___ \| |\  | |_| | |___| |___|  _ < ___) |
// |_| |_/_/   \_\_| \_|____/|_____|_____|_| \_\____/
//

func removeFromCartHandler(
	writer http.ResponseWriter,
	request *http.Request,
	user xdb.User,
	cartId int64,
) {
	var videoId string

	internal.WithTimer("getting video ID from query parameters", func() {
		var queryParameters url.Values = request.URL.Query()
		videoId = queryParameters.Get("v")
	})

	if videoId == "" {
		http.Error(
			writer,
			"Video ID not provided",
			http.StatusBadRequest,
		)
		return
	}

	internal.WithTimer("removing video from cart", func() {
		err = xdb.RemoveItemFromCart(db, cartId, videoId)
	})

	if err != nil {
		http.Error(
			writer,
			"Could not remove item from cart",
			http.StatusInternalServerError,
		)
		return
	}
}

func addToCartHandler(
	writer http.ResponseWriter,
	request *http.Request,
	user xdb.User,
	cartId int64,
) {
	var videoId string

	internal.WithTimer("getting video ID from query parameters", func() {
		var queryParameters url.Values = request.URL.Query()
		videoId = queryParameters.Get("v")
	})

	if videoId == "" {
		http.Error(
			writer,
			"Video ID not provided",
			http.StatusBadRequest,
		)
		return
	}

	internal.WithTimer("adding video to cart", func() {
		err = xdb.AddItemToCart(db, cartId, videoId)
	})

	if err != nil {
		http.Error(
			writer,
			"Could not add item to cart",
			http.StatusInternalServerError,
		)
		return
	}
}

func getCartHandler(
	writer http.ResponseWriter,
	request *http.Request,
	user xdb.User,
	cartId int64,
) {
	var cartItems []xyoutube.Video

	internal.WithTimer("getting cart items", func() {
		cartItems, err = xdb.GetItemsFromCart(db, cartId)
	})

	if err != nil {
		http.Error(
			writer,
			"Could not get items from cart",
			http.StatusInternalServerError,
		)
		return
	}

	var responseBody []byte

	internal.WithTimer("marshalling cart items to JSON", func() {
		response := GetCartResponseBody{CartItems: cartItems}
		responseBody, err = json.Marshal(response)
	})

	if err != nil {
		http.Error(
			writer,
			"Could not marshal cart items to JSON",
			http.StatusInternalServerError,
		)
		return
	}

	writer.Write(responseBody)
}
