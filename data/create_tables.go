package data

import (
	"backend-service/util"
	"context"
)

func (db *Database) InitialiseDatabaseTables(ctx context.Context) error {
	log := util.GetGlobalLogger(ctx)
	if err := db.createUserProfileTable(ctx); err != nil {
		return err
	}
	log.Println("user_profile table successfully created")
	return nil
}

func (db *Database) createUserProfileTable(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS user_profile(
		id VARCHAR(26) PRIMARY KEY CONSTRAINT ulid_size	CHECK (char_length(id) = 26),
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(127) NOT NULL,
		created_on timestamp default NOW()
	);`
	if _, err := db.Pool.Exec(ctx, createTableSQL); err != nil {
		util.GetGlobalLogger(ctx).Println("Failed to execute create query", err)
		return err
	}
	return nil
}
