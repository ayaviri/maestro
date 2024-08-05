package db

import (
	"database/sql"
	"fmt"
	xyoutube "maestro/internal/youtube"
	"strings"
)

func CreateVideos(db *sql.DB, videos []xyoutube.Video) error {
	var b strings.Builder
	_, err = b.WriteString(
		`insert into video (youtube_id, title, channel_title, description, 
published_at, youtube_link, duration_seconds, view_count) values `,
	)

	for _, video := range videos {
		s := fmt.Sprintf(
			`("%s", "%s", "%s", "%s", "%s", "%s", %d, %d) `,
			video.Id,
			video.Title,
			video.ChannelTitle,
			video.Description,
			// TODO: Need to convert this string to a SQLite datetime
			video.PublishedAt,
			video.Link,
			video.DurationSeconds,
			video.ViewCount,
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
