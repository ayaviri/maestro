package db

import (
	"database/sql"
	"fmt"
	xyoutube "maestro/internal/youtube"
)

func AddItemToCart(db *sql.DB, cartId int64, videoId string) error {
	statement := fmt.Sprintf(
		`insert into cart_item (cart_id, video_youtube_id) values(%d, "%s");`,
		cartId, videoId,
	)
	_, err = db.Exec(statement)

	return err
}

func RemoveItemFromCart(db *sql.DB, cartId int64, videoId string) error {
	statement := fmt.Sprintf(
		`delete from cart_item where cart_id=%d and video_youtube_id="%s"`,
		cartId, videoId,
	)
	_, err = db.Exec(statement)

	return err
}

func GetItemsFromCart(db *sql.DB, cartId int64) ([]xyoutube.Video, error) {
	query := fmt.Sprintf(
		`select video.* from video join cart_item on video.youtube_id = 
        cart_item.video_youtube_id where cart_item.cart_id = %d;`,
		cartId,
	)
	fmt.Println(query)
	var rows *sql.Rows
	rows, err = db.Query(query)
	cartItems := make([]xyoutube.Video, 0)

	if err != nil {
		return cartItems, nil
	}

	for rows.Next() {
		var cartItem xyoutube.Video
		err = rows.Scan(
			&cartItem.Id,
			&cartItem.Title,
			&cartItem.ChannelTitle,
			&cartItem.Description,
			&cartItem.PublishedAt,
			&cartItem.Link,
			&cartItem.DurationSeconds,
			&cartItem.ViewCount,
		)

		if err != nil {
			return cartItems, nil
		}

		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}
