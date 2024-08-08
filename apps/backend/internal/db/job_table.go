package db

import (
	"database/sql"
	"fmt"
)

const StatusCreated = "CREATED"
const StatusReceived = "RECEIVED"
const StatusFinished = "FINISHED"

// Creates a new row in the job table and returns its ID
func CreateNewJob(db *sql.DB) (int64, error) {
	statement := fmt.Sprintf(
		`insert into job (status) values("%s");`, StatusCreated,
	)
	var result sql.Result
	result, err = db.Exec(statement)

	if err != nil {
		return 0, err
	}

	var jobId int64
	jobId, err = result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return jobId, nil
}

func UpdateJobStatus(db *sql.DB, jobId int64, status string) error {
	statement := fmt.Sprintf(
		`update job set status = "%s" where id = %d`, status, jobId,
	)
	_, err = db.Exec(statement)

	return err
}
