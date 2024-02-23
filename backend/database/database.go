package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	_ "github.com/lib/pq"
)

type migration struct {
	filename string;
	query string;
}

var db *sql.DB = nil;
func GetDbConnection() *sql.DB {
	if db == nil {
		db = InitDbConnection()
	}
	return db
}

func InitDbConnection() *sql.DB {
	psql_conenction, found := os.LookupEnv("TBL_PSQL_CONNECTION")

	if !found {
		log.Println("Could not find \"TBL_PSQL_CONNECTION\" variable")
		return nil;
	}

	db, err := sql.Open("postgres", psql_conenction)
	if err != nil {
		log.Println("Could not open PostgreSQL connection: ", err)
		return nil;
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Error on starting transaction: ", err)
		return nil
	}

	err = migrateDb(tx)
	if err != nil {
		log.Println("Error during migration")
		err = tx.Rollback()

		if err != nil {
			log.Println("Error rolling back the migration: ", err)
		}
		return nil
	}
	err = tx.Commit()
	if err != nil {
		log.Println("Error during commiting transtaction: ", err)
		return nil
	}

	return db
}

func fetchMigrationsFromFolder(migrationPath string) ([]migration, error) {
	files, err := os.ReadDir(migrationPath)
	migrations := []migration{}

	if err != nil {
		return migrations, err
	}

	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".sql") {
			filePath := filepath.Join(migrationPath, fileName)
			query, err := os.ReadFile(filePath)
			if err != nil {
				return migrations, err
			}
			migrations = append(migrations,
				migration {
					filename: fileName,
					query: string(query),
				})
		}
	}

	sort.Slice(migrations,
		func (i, j int) bool {
			return migrations[i].filename < migrations[j].filename
		})

	return migrations, nil
}

func migrateDb(tx *sql.Tx) error {
	migrations, err := fetchMigrationsFromFolder("./sql/")

	if err != nil {
		return err
	}

	for _, m := range migrations {
		log.Printf("Applying \"%v\" from %v", m.query, m.filename)
		_, err = tx.Exec(m.query)

		if err != nil {
			log.Println("Error!")
			return err
		}
	}

	return nil
}
