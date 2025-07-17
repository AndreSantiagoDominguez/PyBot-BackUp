package connections

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

type ConnPostgreSQL struct {
	DB *sql.DB
	Err string
}

func GetDBPool() *ConnPostgreSQL {
	error := ""

	// Conectar a la base de datos
	db, err := sql.Open("postgres", os.Getenv("URL_POSTGRES"))
	if err != nil {
		error = fmt.Sprintf("error al abrir la base de datos: %v", err)
	}

	// Configuración del pool de conexiones
	db.SetMaxOpenConns(10)

	// Probar conexión
	if err := db.Ping(); err != nil {
		db.Close()
		error = fmt.Sprintf("error al verificar la conexión a la base de datos: %v", err)
	}

	return &ConnPostgreSQL{DB: db, Err: error}
}

func (conn *ConnPostgreSQL) ExecutePreparedQuery(query string, values ...interface{}) (sql.Result, error) {
	stmt, err := conn.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta preparada: %w", err)
	}

	return result, nil
}

func (conn *ConnPostgreSQL) FetchRows(query string, values ...interface{}) (*sql.Rows, error) {
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT: %w", err)
	}

	return rows, nil
}

func (conn *ConnPostgreSQL) QueryRowScan(query string, dest ...interface{}) error {
	return conn.DB.QueryRow(query).Scan(dest...)
}