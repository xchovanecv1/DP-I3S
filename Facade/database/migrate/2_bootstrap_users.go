package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
	guuid "github.com/google/uuid"
)

func init() {

	id := guuid.New()

	bootstrapAdminAccount := `
	INSERT INTO accounts (id, uuid, email, pass, name, active, roles)
	VALUES (DEFAULT, '` + id.String() + `', 'admin@boot.io', '$2a$10$m/NKr.O9keCXAoNZckI80ukbtj8ZHZaOShWxHQuQQbyJaiA4fBPhW', 'Admin Boot', true, '{admin}')
	`

	id = guuid.New()

	bootstrapUserAccount := `
	INSERT INTO accounts (id, uuid, email, pass, name, active)
	VALUES (DEFAULT, '` + id.String() + `', 'user@boot.io', '$2a$10$VPhigWo/C2C.ILqTqCJCJuYasHA3NaqFDfPgNwOEkTWrRPJp8Zs9S', 'User Boot', true)
	`

	up := []string{
		bootstrapAdminAccount,
		bootstrapUserAccount,
	}

	down := []string{
		`TRUNCATE accounts CASCADE`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("add bootstrap accounts")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("truncate accounts cascading")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
