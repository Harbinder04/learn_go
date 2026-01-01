package internal

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type UserStore struct {
	db *sql.DB
}

type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) Create(user User) (string, error) {
   _, err := us.db.Exec("INSERT INTO users (id, name, email) VALUES ($1, $2 , $3)", user.Id, user.Name, user.Email)
   if err != nil {
	return "", fmt.Errorf("Failed to insert: %w", err)
   }

   return  user.Id, nil
}

func (us *UserStore) GetByID(id string) (User, error) {
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

func (us *UserStore) GetAllUser() ([]User, error){
	
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