package dbinit

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/HsimWong/ecommerce/internal/database"
	"github.com/HsimWong/ecommerce/internal/model"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func CreateTables(db ...*gorm.DB) error {
	var dbgorm *gorm.DB
	if len(db) <= 0 {
		// var err error
		dbinstance, err := database.GetDatabaseInstance()
		if err != nil {
			return err
		}
		dbgorm, err = dbinstance.GetDB()
		if err != nil {
			return err
		}

	} else {
		dbgorm = db[0]
	}
	err := AutoCreateTable(dbgorm, &model.User{})
	if err != nil {
		return err
	}

	return nil
}

// AutoCreateTable creates table from GORM-tagged struct if not exists
func AutoCreateTable(db *gorm.DB, model any) error {
	// Get table name from struct name
	tableName := strings.ToLower(reflect.TypeOf(model).Name() + "s")

	// Check if table exists
	if db.Migrator().HasTable(model) {
		log.Printf("Table %s already exists", tableName)
		return nil
	}

	// Create table with all constraints
	err := db.AutoMigrate(model)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate %s: %w", tableName, err)
	}

	log.Printf("Created table %s with schema:\n", tableName)

	// Print table schema (debug)
	var schema strings.Builder
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(model)
	for _, field := range stmt.Schema.Fields {
		schema.WriteString(fmt.Sprintf("- %s (%s)\n", field.DBName, field.Tag))
	}
	log.Println(schema.String())

	return nil
}
