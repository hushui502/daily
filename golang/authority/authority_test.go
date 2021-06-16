package authority

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var dsn string
	if os.Getenv("env") == "testing" {
		fmt.Println("preparing testing config...")
		dsn = fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/db_test?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("ROOT_PASSWORD"))
	} else {
		dsn = "root:@tcp(127.0.0.1:3306)/db_test?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

//TODO: TEST CASE
//func TestCreateRole(t *testing.T) {
//	auth := New(Options{
//		TablesPrefix: "authority_",
//		DB: db,
//	})
//
//
//}
