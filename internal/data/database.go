package data

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

var (
	DbName     = "snape"
	DiskDbPath = fmt.Sprintf("%s.db", DbName)
	DbURI      = fmt.Sprintf("sqlite://./%s", DiskDbPath)
	MemDbURI   = fmt.Sprintf("file:%s?mode=memory&cache=shared", DbName)
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func MigrateDB(dbURI string) {
	dbDriver, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		log.Println(err)
		return
	}

	migrations, err := migrate.NewWithSourceInstance("iofs", dbDriver, dbURI)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Running migrations")
	err = migrations.Up()
	if err != nil && err.Error() != "no change" {
		log.Println(err)
		return
	}
}

func OpenDB(dbPath string) *sql.DB {
	log.Println("Open DB handle to:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func saveDB(db *sql.DB, dbDiskPath string) {
	log.Println("Writing DB to file:", dbDiskPath)
	os.Remove(dbDiskPath)
	statement, err := db.Prepare("vacuum main into ?")
	if err != nil {
		log.Fatal(err)
	}
	statement.Close()

	result, err := statement.Exec(dbDiskPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
}

func loadDB(dbDiskPath string) *sql.DB {
	log.Println("loading saved DB file:", dbDiskPath)
	db := OpenDB(MemDbURI)
	// dbDisk := OpenDB(dbDiskPath)
	// err := backupDB(dbDisk, db)
	// if err != nil {
	// 	panic(err)
	// }
	return db
}

// func backupDB(sourceDB, destDB *sql.DB) error {
// 	destConn, err := destDB.Conn(context.Background())
// 	if err != nil {
// 		return err
// 	}
//
// 	srcConn, err := sourceDB.Conn(context.Background())
// 	if err != nil {
// 		return err
// 	}
//
// 	// Get the raw Conn instance from the underlying SQLite driver,
// 	// which exposes the Backup() functions.
// 	return destConn.Raw(func(destConn interface{}) error {
// 		return srcConn.Raw(func(srcConn interface{}) error {
// 			srcSQLiteConn, ok := srcConn.(*sqlite3.SQLiteConn)
// 			if !ok {
// 				return fmt.Errorf("can't convert source connection to SQLiteConn")
// 			}
//
// 			destSQLiteConn, ok := destConn.(*sqlite3.SQLiteConn)
// 			if !ok {
// 				return fmt.Errorf("can't convert destination connection to SQLiteConn")
// 			}
//
// 			backup, err := destSQLiteConn.Backup("main", srcSQLiteConn, "main")
// 			if err != nil {
// 				return fmt.Errorf("error initializing SQLite backup: %w", err)
// 			}
//
// 			done, err := backup.Step(-1)
// 			if !done {
// 				return fmt.Errorf("step of -1, but not done")
// 			}
// 			if err != nil {
// 				return fmt.Errorf("error in stepping backup: %w", err)
// 			}
//
// 			err = backup.Finish()
// 			if err != nil {
// 				return fmt.Errorf("error finishing backup: %w", err)
// 			}
//
// 			return err
// 		})
// 	})
// }
