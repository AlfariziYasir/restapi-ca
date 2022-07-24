package migration

import (
	"errors"
	"restapi/internal/app/model"
	"restapi/internal/db/postgres"
)

func Up() error {
	pg, err := postgres.NewClient()
	if err != nil {
		return err
	}

	err = pg.Conn().AutoMigrate(
		&model.User{},
	)
	return ignoreErrNoChange(err)
}

func Drop() error {
	pg, err := postgres.NewClient()
	if err != nil {
		return err
	}

	err = pg.Conn().Migrator().DropTable(
		&model.User{},
	)
	return ignoreErrNoChange(err)
}

func ignoreErrNoChange(err error) error {
	if err != nil && err != errors.New("no change") {
		return err
	}

	return nil
}
