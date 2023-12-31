package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

var (
	db *gorm.DB = nil
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func setup() error {
	username := getEnv("PSQLUSER", "postgres")
	password := getEnv("PSQLPASS", "") //3L5JaSDTDC"
	dbName := getEnv("PSQLDB", "postgres")
	dbHost := getEnv("PSQLHOST", "localhost")
	dbPort := getEnv("PSQLPORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", dbHost, username, dbName, password, dbPort)
	//dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	fmt.Println("DSN: ", dsn)
	dbase, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Println("We cant op open a DATABASE")
		return err
	}

	db = dbase.Debug()

	//db.AutoMigrate(&models.Location{})

	return nil
}

func TryGetDDB() (*gorm.DB, error) {
	if db == nil {
		err := setup()
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func GetDB() *gorm.DB {
	/*if db == nil {
		err := setup()
		if err != nil {
			log.Fatalln(err)
		}
	}*/
	db, err := TryGetDDB()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
