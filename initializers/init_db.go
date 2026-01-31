package initializers

import (
	"fmt"
	"os"
	"time"

	"github.com/AKSHAY-1505/auth-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, username, password, database, port)

	Log.Infof("[SERVER] Initializing connection to PostgresDB Host: %s", host)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		Log.Fatalf("[SERVER] Unable to connect to database host: %s", host)
	}

	Log.Info("[SERVER] Successfully connected to PostgresDB")

	configureConnectionPool(db)
	migrateDB(db)

	DB = db
}

func configureConnectionPool(db *gorm.DB) {
	sqlDB, _ := db.DB()

	Log.Info("[SERVER] Configuring Database Connection Pool")

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// SetConnMaxIdleTime sets the maximum amount of time a connection can be kept idle.
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	Log.Info("[SERVER] Successfully Configured Database Connection Pool")
}

func migrateDB(db *gorm.DB) {
	// Auto-migrate the User struct
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		Log.Fatalf("[SERVER] Failed to migrate user table in database")
	}

	Log.Info("[SERVER] Successfully migrated all tables in database")
}
