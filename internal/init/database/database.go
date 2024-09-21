package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mojganchakeri/whatsapp-manager/internal/init/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	GetConnection(ctx context.Context) *gorm.DB
	GetDSN() string
}

type database struct {
	dsn        string
	dbConnOnce sync.Once
	db         *gorm.DB
}

func New(cfg config.Config) Database {
	host := cfg.GetConfig().Database.Host
	port := cfg.GetConfig().Database.Port
	user := cfg.GetConfig().Database.Username
	password := cfg.GetConfig().Database.Password
	db := cfg.GetConfig().Database.Database

	return &database{
		dsn: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Tehran",
			host, user, password, db, port),
	}
}

func (d *database) GetConnection(ctx context.Context) *gorm.DB {
	if d.db == nil {
		d.dbConnOnce.Do(func() {
			var err error
			d.db, err = gorm.Open(postgres.Open(d.dsn), &gorm.Config{
				NowFunc: func() time.Time {
					ti, _ := time.LoadLocation("Asia/Tehran")
					return time.Now().In(ti)
				},
			})
			if err != nil {
				panic(err)
			}

			db, err := d.db.DB()
			if err != nil {
				panic(err)
			}
			db.SetMaxIdleConns(10)
			db.SetMaxOpenConns(25)
			db.SetConnMaxLifetime(time.Hour)
		})
	}
	return d.db.WithContext(ctx)
}

func (d *database) GetDSN() string {
	return d.dsn
}
