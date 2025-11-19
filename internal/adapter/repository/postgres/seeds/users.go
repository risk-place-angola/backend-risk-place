package seeds

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func SeedUsers(ctx context.Context, db *sql.DB) error {
	var passwordHash = "#Pwd1234"

	generateFromPassword, err := bcrypt.GenerateFromPassword([]byte(passwordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	insertUsersQuery := `
		INSERT INTO users (id, name, email, password, phone, latitude, longitude, email_verified, country, province, municipality, neighborhood)
		VALUES
		('7aa9c0c0-14d3-4af9-a631-956c12a1f100', 'lopes Estevão', 'lopes@example.com', $1, '+244923111111', -8.839987, 13.289437, TRUE, 'Angola', 'Luanda', 'Talatona', 'Benfica'),
		('f55c21ea-18e3-4fc9-99c3-d03b234bc110', 'João Silva', 'joao@example.com', $2, '+244923222222', -8.915120, 13.242380, TRUE, 'Angola', 'Luanda', 'Viana', 'Zango 2'),
		('a11fd8dc-55d0-4e07-b0db-23f659ed3201', 'Maria João', 'maria@example.com', $3, '+244923333333', -8.828765, 13.247865, TRUE, 'Angola', 'Luanda', 'Talatona', 'Gamek'),
		('cc772485-bb16-4584-a8a4-3fd366478931', 'Carlos Domingos', 'carlos@example.com', $4, '+244923444444', -8.842560, 13.300120, TRUE, 'Angola', 'Luanda', 'Samba', 'Morro Bento'),
		('51bfcbfd-c896-4a7a-ae6a-79a0df2aab30', 'Ana Ferreira', 'ana@example.com', $5, '+244923555555', -8.903290, 13.312540, TRUE, 'Angola', 'Luanda', 'Cacuaco', 'Sequele')
		ON CONFLICT (id) DO NOTHING
	`

	_, err = db.ExecContext(ctx, insertUsersQuery,
		string(generateFromPassword),
		string(generateFromPassword),
		string(generateFromPassword),
		string(generateFromPassword),
		string(generateFromPassword),
	)
	if err != nil {
		return err
	}

	insertUserRolesQuery := `
		INSERT INTO user_roles (user_id, role_id)
		SELECT id, (SELECT id FROM roles WHERE name='citizen' LIMIT 1)
		FROM users
		WHERE id IN (
			'7aa9c0c0-14d3-4af9-a631-956c12a1f100',
			'f55c21ea-18e3-4fc9-99c3-d03b234bc110',
			'a11fd8dc-55d0-4e07-b0db-23f659ed3201',
			'cc772485-bb16-4584-a8a4-3fd366478931',
			'51bfcbfd-c896-4a7a-ae6a-79a0df2aab30'
		)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err = db.ExecContext(ctx, insertUserRolesQuery)
	return err
}
