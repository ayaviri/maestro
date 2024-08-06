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
		"insert into search_result (search_id, video_youtube_id) values",
	)

	for index, video := range videos {
		s := fmt.Sprintf(
			`(%d, "%s"),`, searchId, video.Id,
		)

		if index == len(videos)-1 {
			s = strings.TrimSuffix(s, ",")
		}

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
