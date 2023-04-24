package config

import (
	"database/sql"
	"fmt"
	mysql2 "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbSet struct {
	master *gorm.DB
	slave  *gorm.DB
}

type Manager struct {
	dbs map[string]*dbSet
}

// NewManager initialize all db service use
func NewManager(dbConfigs map[string]*DBConf) (*Manager, error) {
	dbSetMap := make(map[string]*dbSet, len(dbConfigs))
	for name, dbConf := range dbConfigs {
		master, err := openDB(dbConf.Master, dbConf, false)
		if err != nil {
			return nil, err
		}
		slave, err := openDB(dbConf.Slave, dbConf, true)
		if err != nil {
			return nil, err
		}
		dbSetMap[name] = &dbSet{
			master: master,
			slave:  slave,
		}
	}
	return &Manager{dbs: dbSetMap}, nil
}

func (dbm *Manager) GetMasterDB(name string) *gorm.DB {
	dbInfo, exist := dbm.dbs[name]
	if !exist {
		panic(fmt.Sprintf("master db:%s not exist", name))
	}
	return dbInfo.master
}

func (dbm *Manager) GetSlaveDB(name string) *gorm.DB {
	dbInfo, exist := dbm.dbs[name]
	if !exist {
		panic(fmt.Sprintf("slave db:%s not exist", name))
	}
	return dbInfo.slave
}

// assume that all db is mysql, also can use driver to support other db, like pgsql
func openDB(dsn string, dbConf *DBConf, isSlave bool) (*gorm.DB, error) {
	dsn = attachDBDsn(dsn, dbConf.AutoCommit, isSlave)
	db, err := gorm.Open(mysql2.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	sqlDB.SetMaxIdleConns(dbConf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(dbConf.MaxLifetime)
	err = sqlDB.Ping()
	// TODO:gorm support plugin, can use it to send qps, execute time to prometheus
	return db, err

}

func attachDBDsn(dsn string, autocommit, isSlave bool) string {
	dsn += "?interpolateParams=true&parseTime=true&loc=Local&autocommit="
	if isSlave {
		dsn += "1"
	} else {
		if autocommit {
			dsn += "1"
		} else {
			dsn += "0"
		}
	}
	return dsn
}
