package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type UserRepository interface {
	UserExists(ctx context.Context, email string) (time.Duration, bool, error)
	Create(ctx context.Context, u User) (time.Duration, string, error)
	GetByID(ctx context.Context, id string) (time.Duration, User, error)
	GetAllUser(ctx context.Context) (time.Duration, []User, error)
}

type SQLUserStore struct {
	db *sql.DB
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewSQLUserStore(db *sql.DB) UserRepository {
	return &SQLUserStore{db: db}
}

// store.go
func (us *SQLUserStore) UserExists(ctx context.Context, email string) (time.Duration, bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	timeStart := time.Now()
	err := us.db.QueryRowContext(ctx, query, email).Scan(&exists)
	dur := time.Since(timeStart)

	if err != nil {
		return dur, false, err
	}
	return dur, exists, nil
}

// ⚠️Todo: Remove tx 
func (us *SQLUserStore) Create(ctx context.Context, user User) (time.Duration, string, error) {
	// ctx := context.TODO()
	timeStart := time.Now()
	// tx, err := us.db.BeginTx(ctx, nil)
	// if err != nil {
	// 	return "", err
	// }

	// defer func(){
	// 	if err != nil {
	// 		tx.Rollback()
	// 	}
	// }()

	//    _, err = tx.ExecContext(ctx, "INSERT INTO users (id, name, email) VALUES ($1, $2 , $3)", user.Id, user.Name, user.Email)
	//    if err != nil {
	// 	return "", err
	//    }

	// todo: Remove later (Pretending)
	//    _, err = tx.ExecContext(ctx, "INSERT INTO audit_logs (action) VALUES ($1)",
	//     "USER_CREATED")
	// 	if err != nil {
	// 		return  "", err
	// 	}

	// 	err = tx.Commit()
	// 	if err != nil {
	// 		return "", err
	// 	}

	//🔨 us.db is replaced with tx here temporarily
	_, err := us.db.ExecContext(ctx, "Insert into users (id, name, email) Values ($1, $2, $3)", user.Id, user.Name, user.Email)
	dur := time.Since(timeStart)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			fmt.Println("Message:", err.Message)
		}
		return dur, "", err
	}

	return dur, user.Id, nil
}

func (us *SQLUserStore) GetByID(ctx context.Context, id string) (time.Duration, User, error) {
	var u User
	timeStart := time.Now()
	result, err := us.db.QueryContext(ctx, "SELECT * FROM users WHERE id = $1", id)
	dur := time.Since(timeStart)

	if err != nil {
		return dur, User{}, err
	}

	if err := result.Scan(&u.Id, &u.Name, &u.Email); err != nil {
		if err == sql.ErrNoRows {
			return dur, User{}, fmt.Errorf("userById %s: no such user", id)
		}
		return dur, User{}, err
	}

	return dur, u, nil
}

func (us *SQLUserStore) GetAllUser(ctx context.Context) (time.Duration, []User, error) {

	var users []User

	timeStart := time.Now()
	result, err := us.db.QueryContext(ctx, "SELECT * FROM users")
	dur := time.Since(timeStart)

	if err != nil {
		if err == context.DeadlineExceeded {
			return dur, nil, err
		}
		return dur, nil, fmt.Errorf("Unable to fetch record")
	}
	defer result.Close()

	for result.Next() {
		var user User
		if err := result.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			return dur, nil, err
		}
		users = append(users, user)
	}
	if err := result.Err(); err != nil {
		return dur, nil, fmt.Errorf("Error: %v", err)
	}

	return dur, users, nil
}
