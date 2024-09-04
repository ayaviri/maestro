package db

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

const StatusCreated = "CREATED"
const StatusReceived = "RECEIVED"
const StatusFinished = "FINISHED"

func CreateNewJobWithId(db *sql.DB, userId string, jobId string) error {
	statement := fmt.Sprintf(
		`insert into job (id, app_user_id, status) values('%s', '%s', '%s');`,
		jobId,
		userId,
		StatusCreated,
	)
	_, err = db.Exec(statement)

	return err
}

func CreateNewJob(db *sql.DB, userId string) (string, error) {
	jobId := uuid.NewString()
	err = CreateNewJobWithId(db, userId, jobId)

	return jobId, err
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

func GetJobPayload(db *sql.DB, jobId string) (string, error) {
	var responsePayload string
	query := fmt.Sprintf(`select response_payload from job where id = '%s'`, jobId)
	var row *sql.Row = db.QueryRow(query)
	err = row.Scan(&responsePayload)

	return responsePayload, err
}
