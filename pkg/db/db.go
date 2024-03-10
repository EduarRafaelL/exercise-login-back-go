package db

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb" // Importa el driver de SQL Server
)

// InitializeDatabase abre una conexión a la base de datos y la retorna.
func InitializeDatabase(dbDriver, dbSource string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
		return nil, err
	}

	// Prueba la conexión con db.Ping()
	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping the database: %v", err)
		return nil, err
	}

	log.Println("Database connection established successfully")
	return db, nil
}
