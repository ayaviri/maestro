package db

import (
	"database/sql"
	"fmt"
	xyoutube "maestro/internal/youtube"
	"strings"
)

func CreateSearchResults(db *sql.DB, searchId int64, videos []xyoutube.Video) error {
	var b strings.Builder
	_, err = b.WriteString(
		"insert into search_result (search_id, video_youtube_id) values ",
	)

	for _, video := range videos {
		s := fmt.Sprintf(
			`(%d, "%s")`, searchId, video.Id,
		)
		b.WriteString(s)
	}

	b.WriteString(";")
	statement := b.String()
	_, err = db.Exec(statement)

	if err != nil {
		return err
	}

	return nil
}
