package db

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

const StatusCreated = "CREATED"
const StatusReceived = "RECEIVED"
const StatusFinished = "FINISHED"

// Creates a new row in the job table and returns its ID
func CreateNewJob(db *sql.DB) (string, error) {
	jobId := uuid.NewString()
	statement := fmt.Sprintf(
		`insert into job (id, status) values('%s', '%s');`, jobId, StatusCreated,
	)
	_, err = db.Exec(statement)

	if err != nil {
		return "", err
	}

	return jobId, nil
}

func UpdateJobStatus(db *sql.DB, jobId string, status string) error {
	statement := fmt.Sprintf(
		`update job set status = '%s' where id = '%s'`, status, jobId,
	)
	_, err = db.Exec(statement)

	return err
}

func AddJobPayload(db *sql.DB, jobId string, payload string) error {
	statement := fmt.Sprintf(
		`update job set response_payload = '%s' where id = '%s'`, payload, jobId,
	)
	_, err = db.Exec(statement)

	return err
}
