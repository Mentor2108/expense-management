package data

import (
	"backend-service/defn"
	"backend-service/util"
	"context"
	"fmt"
	"strings"
)

type UserRepository struct {
	db Database
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: GetDatabaseConnection()}
}

func (repo *UserRepository) Create(ctx context.Context, userData defn.UserProfile) (map[string]interface{}, *util.CustomError) {
	// log := util.GetGlobalLogger(ctx)

	userData.ID = util.ULID()

	row := repo.db.Pool.QueryRow(ctx, "INSERT INTO user_profile (id, email, password) VALUES ($1, $2, $3) RETURNING id, email", userData.ID, userData.Email, userData.Password)
	err := row.Scan(&userData.ID, &userData.Email)
	// log.Printf("userData: %+v", userData)
	// rows, err := repo.db.Pool.Query(ctx, "INSERT INTO user_profile (id, email, password) VALUES ($1, $2, $3)",
	// 	userData.ID, userData.Email, userData.Password)
	// if err != nil {
	// 	cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseCreateOperationFailed, defn.ErrDatabaseCreateOperationFailed, map[string]string{
	// 		"error": err.Error(),
	// 	})
	// 	// log.Println(cerr)
	// 	return nil, cerr
	// }
	// defer rows.Close()

	// // Get column names dynamically
	// fieldDescriptions := rows.FieldDescriptions()
	// columns := make([]string, len(fieldDescriptions))
	// for i, fd := range fieldDescriptions {
	// 	columns[i] = fd.Name
	// }

	// log.Println(columns)

	// // Read the first and only row
	// if !rows.Next() {
	// 	cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseCreateOperationFailed, defn.ErrDatabaseCreateOperationFailed, map[string]string{
	// 		"error": "no rows found",
	// 	})
	// 	// log.Println(cerr)
	// 	return nil, cerr
	// }

	// values, err := rows.Values() // Get all values in a slice
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseCreateOperationFailed, defn.ErrDatabaseCreateOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}

	newUser := map[string]interface{}{
		"id":    userData.ID,
		"email": userData.Email,
	}

	// // Map column names to values
	// newUser := make(map[string]interface{})
	// for i, column := range columns {
	// 	newUser[column] = values[i]
	// }

	return newUser, nil
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (map[string]interface{}, *util.CustomError) {
	// log := util.GetGlobalLogger(ctx)

	rows, err := repo.db.Pool.Query(ctx, "SELECT * from user_profile where email = $1", email)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}
	defer rows.Close()

	// Fetch column names dynamically
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}

	// Read the first row
	if !rows.Next() {
		cerr := util.NewCustomError(ctx, defn.ErrCodeNoDataFound, defn.ErrNoDataFound)
		// log.Println(cerr)
		return nil, cerr
	}

	values, err := rows.Values() // Get all values in a slice
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}

	// Store values dynamically in map
	userData := make(map[string]interface{})
	for i, column := range columns {
		userData[column] = values[i]
	}
	return userData, nil
}

func (repo *UserRepository) GetByID(ctx context.Context, id string) (map[string]interface{}, *util.CustomError) {
	// log := util.GetGlobalLogger(ctx)

	rows, err := repo.db.Pool.Query(ctx, "SELECT * from user_profile where id = $1", id)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}
	defer rows.Close()

	// Fetch column names dynamically
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}

	// Read the first row
	if !rows.Next() {
		cerr := util.NewCustomError(ctx, defn.ErrCodeNoDataFound, defn.ErrNoDataFound)
		// log.Println(cerr)
		return nil, cerr
	}

	values, err := rows.Values() // Get all values in a slice
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}

	// Store values dynamically in map
	userData := make(map[string]interface{})
	for i, column := range columns {
		if column == "password" {
			continue
		}
		userData[column] = values[i]
	}
	return userData, nil
}

func (repo *UserRepository) UpdateById(ctx context.Context, id string, updates map[string]interface{}) (map[string]interface{}, *util.CustomError) {
	// log := util.GetGlobalLogger(ctx)

	if len(updates) == 0 {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseUpdateOperationFailed, defn.ErrDatabaseUpdateOperationFailed, map[string]string{
			"error": "no fields found for updating",
		})
		// log.Println(cerr)
		return nil, cerr
	}

	setClauses := []string{}
	args := []interface{}{id}
	argIndex := 2

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf("UPDATE user_profile SET %s WHERE id = $1 RETURNING *", strings.Join(setClauses, ", "))
	rows, err := repo.db.Pool.Query(ctx, query, args...)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseUpdateOperationFailed, defn.ErrDatabaseUpdateOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}
	defer rows.Close()

	// Get column names dynamically
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}

	// Read the first and only row
	if !rows.Next() {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseUpdateOperationFailed, defn.ErrDatabaseUpdateOperationFailed, map[string]string{
			"error": "no rows found",
		})
		// log.Println(cerr)
		return nil, cerr
	}

	values, err := rows.Values() // Get all values in a slice
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseUpdateOperationFailed, defn.ErrDatabaseUpdateOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}

	// Map column names to values
	updatedUser := make(map[string]interface{})
	for i, column := range columns {
		if column == "password" {
			continue
		}
		updatedUser[column] = values[i]
	}

	return updatedUser, nil
}
