package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type UserStore interface {
	UserExists(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, u User) (string, error)
	GetByID(id string) (User, error)
	GetAllUser() ([]User, error)
}

type SQLUserStore struct {
	db *sql.DB
}


type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func NewSQLUserStore(db *sql.DB) UserStore {
	return &SQLUserStore{db: db}
}

// store.go
func (us *SQLUserStore) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := us.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// refactored to show demo transaction logic 
func (us *SQLUserStore) Create(ctx context.Context, user User) (string, error) {
	// ctx := context.TODO()
	tx, err := us.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	
	defer func(){
		if err != nil {
			tx.Rollback()
		}
	}()


   _, err = tx.ExecContext(ctx, "INSERT INTO users (id, name, email) VALUES ($1, $2 , $3)", user.Id, user.Name, user.Email)
   if err != nil {
	return "", err
   }

   // todo: Remove later (Pretending)
   _, err = tx.ExecContext(ctx, "INSERT INTO audit_logs (action) VALUES ($1)",
    "USER_CREATED")
	if err != nil {
		return  "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}
	
   return  user.Id, nil
}

func (us *SQLUserStore) GetByID(id string) (User, error) {
	var u User
	result := us.db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	
	if err := result.Scan(&u.Id, &u.Name, &u.Email); err != nil {
        if err == sql.ErrNoRows {
            return User{}, fmt.Errorf("userById %s: no such user", id)
        }
		return User{}, err
	}

	return u, nil
}

func (us *SQLUserStore) GetAllUser() ([]User, error){
	
	var users []User
	result, err := us.db.Query("SELECT * FROM users"); if err != nil {
		return nil, fmt.Errorf("Unable to fetch record")
	}
	defer result.Close()

	for result.Next() {
		var user User
		if err := result.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			slog.Info("Unable to scan a row")
		}
		users = append(users, user)
	}
	if err := result.Err(); err != nil {
        return nil, fmt.Errorf("Error: %v", err)
    }

	return users, nil
}