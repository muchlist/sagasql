package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

const (
	dbname              = "testsaga"
	postgreUserHost     = "PG_USER_HOST"
	postgreUserPort     = "PG_USER_PORT"
	postgreUserUName    = "PG_USER_UNAME"
	postgreUserPassword = "PG_USER_PASSWORD"
)

var (
	// DB global variable agar database dapat diakses di dao
	DB *pgxpool.Pool
)

// InitDB menginisiasi database
// responsenya digunakan untuk memutus koneksi apabila main program dihentikan
func InitDB() *pgxpool.Pool {

	host := os.Getenv(postgreUserHost)         // localhost
	port := os.Getenv(postgreUserPort)         // 5432
	user := os.Getenv(postgreUserUName)        // postgres
	password := os.Getenv(postgreUserPassword) // postgres

	// databaseUrl := "postgres://username:password@localhost:5432/database_name"
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	var err error
	DB, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	fmt.Println("Connected!")

	return DB
}
