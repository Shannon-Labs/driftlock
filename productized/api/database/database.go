package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"driftlock/productized/api/models"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(dsn string, debug bool) {
	var err error
	var config = &gorm.Config{}
	
	if debug {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	err = DB.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.TenantSettings{},
		&models.Anomaly{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connection established and migrations completed")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database not initialized")
	}
	return DB
}

// SeedDB adds initial data to the database
func SeedDB() {
	// Create default admin user if not exists
	var count int64
	DB.Model(&models.User{}).Where("email = ?", "admin@driftlock.example").Count(&count)
	
	if count == 0 {
		adminUser := &models.User{
			Email:    "admin@driftlock.example",
			Name:     "Admin User",
			Role:     "admin",
		}
		
		// Hash the password
		password := "defaultAdminPassword123!" // In production, use a strong password
		// Note: In a real implementation, you would hash this password
		
		adminUser.Password = password
		result := DB.Create(adminUser)
		
		if result.Error != nil {
			log.Printf("Failed to create admin user: %v", result.Error)
		} else {
			log.Printf("Created admin user with ID: %d", adminUser.ID)
		}
	}
}