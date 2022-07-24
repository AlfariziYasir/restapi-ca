package postgres

import (
	"context"
	"fmt"

	"restapi/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client interface {
	Conn() *gorm.DB
	Close() error
}

func NewClientPG() (Client, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.Cfg().DBHost,
		config.Cfg().DBPort,
		config.Cfg().DBUser,
		config.Cfg().DBName,
		config.Cfg().DBPass,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	test, _ := db.DB()
	err = test.PingContext(context.Background())
	if err != nil {
		return nil, err
	}

	return &client{db}, nil
}

func NewClient() (Client, error) {
	return NewClientPG()
}

type client struct {
	db *gorm.DB
}

func (c *client) Conn() *gorm.DB { return c.db }
func (c *client) Close() error {
	db, _ := c.db.DB()
	return db.Close()
}
