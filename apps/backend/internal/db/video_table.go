package db

import (
	"database/sql"
	"fmt"
	xyoutube "maestro/internal/youtube"
	"strings"
)

func CreateVideos(db *sql.DB, videos []xyoutube.Video) error {
	var b strings.Builder
	// TODO: Do we want to update videos on primary key conflict ?
	_, err = b.WriteString(
		`insert into video (youtube_id, title, channel_title, description, 
published_at, youtube_link, duration_seconds, view_count) values`,
	)

	for index, video := range videos {
		s := fmt.Sprintf(
			`('%s', '%s', '%s', '%s', '%s', '%s', %d, %d),`,
			video.Id,
			strings.ReplaceAll(video.Title, "'", "''"),
			strings.ReplaceAll(video.ChannelTitle, "'", "''"),
			strings.ReplaceAll(video.Description, "'", "''"),
			// TODO: Need to convert this string to a SQLite datetime
			video.PublishedAt,
			video.Link,
			video.DurationSeconds,
			video.ViewCount,
		)

		if index == len(videos)-1 {
			s = strings.TrimSuffix(s, ",")
		}

		b.WriteString(s)
	}

	b.WriteString(" on conflict do nothing;")
	statement := b.String()
	_, err = db.Exec(statement)

	if err != nil {
		return err
	}

	return nil
}
