package db

// import (
// 	"database/sql"
// 	"fmt"
// )

// func CreateCart(db *sql.DB, userId int64) error {
// 	statement := fmt.Sprintf(`insert into cart (user_id) values(%d);`, userId)
// 	_, err = db.Exec(statement)
//
// 	return err
// }

// func GetUserCartId(db *sql.DB, userId int64) (int64, error) {
// 	statement := fmt.Sprintf("select id from cart where user_id=%d", userId)
// 	var row *sql.Row = db.QueryRow(statement)
// 	var cartId int64
// 	err = row.Scan(&cartId)
//
// 	return cartId, err
// }
