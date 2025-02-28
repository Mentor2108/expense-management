package data

import (
	"backend-service/util"
	"context"
)

func (db *Database) InitialiseDatabaseTables(ctx context.Context) error {
	// log := util.GetGlobalLogger(ctx)
	// if err := db.createUserTable(ctx); err != nil {
	// 	return err
	// }
	// log.Println("scrape_job table successfully created")
	return nil
}

func (db *Database) createUserTable(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS scrape_job(
		id VARCHAR(26) PRIMARY KEY CONSTRAINT ulid_size	CHECK (char_length(id) = 26),
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(24) NOT NULL,
		created_on timestamp default NOW()
	);`
	if _, err := db.Pool.Exec(ctx, createTableSQL); err != nil {
		util.GetGlobalLogger(ctx).Println("Failed to execute create query", err)
		return err
	}
	return nil
}
