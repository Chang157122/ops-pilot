package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"opsPilot/internal/common/util"
	"opsPilot/internal/pkg/log"
	"os"
	"reflect"
	"strings"
	"time"
)

var db *gorm.DB

type Model struct {
	ID         int       `gorm:"primary_key" json:"id"`
	CreateOn   time.Time `gorm:"type:date;column:create_on" json:"create_on"`
	ModifiedOn time.Time `gorm:"type:date;column:modified_on" json:"modified_on"`
}

func init() {
	var (
		err error
		//dbType      = "sqlite3"
		dbFile        string
		dbFilePathDir string // sqlite3的存放路径
		tablePrefix   = "t_"
	)
	if dbFile = os.Getenv("DB_PATH"); dbFile == "" {
		dbFile = "sqlite3.db"
	} else {
		s := strings.Split(dbFile, "/")
		// 获取数组长度
		c := len(s)
		for k, v := range s {
			if c-1 == k {
				break
			}
			dbFilePathDir += "/" + v
		}
		util.CheckDirIsExist(dbFilePathDir)
	}

	db, err = gorm.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)

	// 创建表
	createTable()
}

func CloseDB() {
	defer db.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now()
		if createTimeField, ok := scope.FieldByName("create_on"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("modified_on"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}

	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// updateTimeStampForUpdateCallback will set `ModifyTime` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("modified_on", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func createTable() {
	models := []interface{}{
		&AuthLogin{},
	}

	for _, model := range models {
		if !db.HasTable(model) {
			table := db.CreateTable(model)
			if table.Error != nil {
				log.Logger.Errorf("create table failed! err: %v", table.Error)
			}
		}
	}
}
