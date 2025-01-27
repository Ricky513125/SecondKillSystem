package data

import (
	"SecKill/conf"
	"SecKill/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // _ to avoid unused package, just use the init
	"log"
	"time"
)

var Db *gorm.DB

// init the mysql, connect
func initMysql(config conf.AppConfig) {
	fmt.Println("Load Database Service configuration...")

	// set the parameters of connection
	dbType := config.App.Database.Type
	usr := config.App.Database.User
	pwd := config.App.Database.Password
	address := config.App.Database.Address
	dbName := config.App.Database.DbName
	dbLink := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		usr, pwd, address, dbName)
	// use this to form different database link

	// create a connection to a database.
	// It needs to connect again for the service in docker starts later.
	fmt.Println("Init Database Service connections...")
	var err error
	// func Open(dialect string, args ...interface{}) (db *DB, err error)
	maxRetries := 5
	for attempt := 0; attempt < maxRetries; attempt++ {
		Db, err = gorm.Open(dbType, dbLink)
		if err == nil {
			break // successfully connect
		}
		log.Println("Attempt %d failed: %v. Retrying in 5 seconds...", attempt, err)
		time.Sleep(5 * time.Second)
	}
	//for Db, err = gorm.Open(dbType, dbLink); err != nil; Db, err = gorm.Open(dbType, dbLink) {
	//	//
	//	log.Println("Failed to connect database: ", err.Error())
	//	log.Println("Reconnecting in 5 seconds...")
	//	time.Sleep(5 * time.Second)
	//}
	if err != nil {
		log.Fatalf("Failed to connect after %d attempts: %v", maxRetries, err)
	}

	// the max open conns allow : small -> wait long; big -> waste space
	// the max idle conns : to keep the connection in the pool to save the spending on the next connection
	Db.DB().SetMaxOpenConns(config.App.Database.MaxOpen)
	Db.DB().SetMaxIdleConns(config.App.Database.MaxIdle)

	// init the database
	user := model.User{}      // don't change the user
	coupon := &model.Coupon{} // change the coupon's amount

	// create tables = [user, coupon]
	tables := []interface{}{user, coupon}

	for _, table := range tables {
		if !Db.HasTable(table) { // check if exists
			Db.AutoMigrate(table) // auto migrate to create the table
		}
	}

	if config.App.FlushAllForTest {
		println("FlushAllForTest is true. Delete records of all tables.")
		for _, table := range tables {
			Db.Delete(table) // delete the content, keep the table
		}
	}

	// create unique index
	Db.Model(user).AddUniqueIndex("username_index", "username")
	Db.Model(coupon).AddUniqueIndex("coupon_index", "username", "coupon_name")

	println("---Mysql connection is initialized.---")
	// 添加外键的demo代码
	// Db.Model(credit_card).
	//	 AddForeignKey("owner_id", "users(id)", "RESTRICT", "RESTRICT").
	//	 AddUniqueIndex("unique_owner", "owner_id")
}
