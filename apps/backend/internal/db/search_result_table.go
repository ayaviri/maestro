package db

import (
	"database/sql"
	"fmt"
	xyoutube "maestro/internal/youtube"
	"strings"

	"github.com/google/uuid"
)

func CreateSearchResults(db *sql.DB, searchId string, videos []xyoutube.Video) error {
	var searchResultId string
	var b strings.Builder
	_, err = b.WriteString(
		"insert into search_result (id, search_id, video_youtube_id) values",
	)

	for index, video := range videos {
		searchResultId = uuid.NewString()
		s := fmt.Sprintf(
			`('%s', '%s', '%s'),`, searchResultId, searchId, video.Id,
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
