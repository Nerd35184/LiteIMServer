package repository

import (
	"errors"
	"fmt"
	"server/conf"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	COUNT_SELECT = " count(1) as count "
)

type MysqlModel interface {
	TableName() string
	Indexes() []*MySqlIndex
}

type CountModel struct {
	Count int64
}

type MaxModel struct {
	M int64
}

type MySqlIndex struct {
	Unique bool
	Cols   []string
}

type Mysql struct {
	DB *gorm.DB
}

func (mysql *Mysql) CreateIndex(tableName string, unique bool, indexName string, cols ...string) error {
	if len(cols) == 0 {
		return errors.New("len(cols) == 0")
	}
	colsStr := strings.Join(cols, ",")
	indexType := "INDEX"
	if unique {
		indexType = "UNIQUE INDEX"
	}
	sql := fmt.Sprintf("CREATE %s %s ON %s (%s)", indexType, indexName, tableName, colsStr)
	return mysql.DB.Exec(sql).Error
}

func InitMysqlDb(mysqlConfig *conf.MysqlConfig) (*Mysql, error) {
	url := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.DbName)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	mysqlModels := []MysqlModel{
		&AddContactRequest{},
		&ContactInfo{},
		&DialogueSess{},
		&MessageInfo{},
		&UserInfo{},
		&UserSessionInfo{},
	}
	objects := make([]interface{}, 0, len(mysqlModels))
	for _, mysqlModel := range mysqlModels {
		objects = append(objects, mysqlModel)
	}
	err = db.AutoMigrate(
		objects...,
	)
	if err != nil {
		panic(err)
	}
	mysqlDb := &Mysql{
		DB: db,
	}
	// for _, mysqlModel := range mysqlModels {
	// 	for _, index := range mysqlModel.Indexes() {
	// 		err := mysqlDb.CreateIndex(mysqlModel.TableName(), index.Unique, strings.Join(index.Cols, ","), index.Cols...)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }
	return mysqlDb, nil
}
