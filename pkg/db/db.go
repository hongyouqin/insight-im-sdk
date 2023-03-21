package db

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 本地数据库
type DataBase struct {
	userId  string
	dbDir   string
	db      *gorm.DB
	rwMutex sync.RWMutex
}

func (d *DataBase) initDB() error {
	d.rwMutex.Lock()
	defer d.rwMutex.Unlock()

	dbName := d.dbDir + "im_" + d.userId + ".db"
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return err
	}
	d.db = db

	//创建表

	return nil
}
