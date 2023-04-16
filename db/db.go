package db

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func New() (*Db, error) {
	conf := getConfig()

	dsn := url.URL{
		User:   url.UserPassword(conf.Username, conf.Password),
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Path:   conf.Database,
		RawQuery: (&url.Values{
			"sslmode": []string{
				"disable",
			},
			"TimeZone": []string{
				"Asia/Jakarta",
			},
		}).Encode(),
	}
	db, err := gorm.Open(
		postgres.Open(dsn.String()),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	return &Db{
		db,
	}, nil
}
