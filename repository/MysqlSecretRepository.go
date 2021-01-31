package repository

import (
	"database/sql"

	"github.com/bernardosecades/sharesecret/types"
	_ "github.com/go-sql-driver/mysql"

	uuid "github.com/satori/go.uuid"

	"fmt"
	"time"
)

const formatDate = "2006-01-02 15:04:05"

type MySqlSecretRepository struct {
	SQL *sql.DB
}

func NewMySqlSecretRepository(dbName string, dbUser string, dbPass string, dbHost string, dbPort string) *MySqlSecretRepository {
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
	d, err := sql.Open("mysql", dbSource)
	if err != nil {
		panic(err)
	}

	return &MySqlSecretRepository{SQL: d}
}

func (r *MySqlSecretRepository) GetSecret(id string) (types.Secret, error) {

	res := r.SQL.QueryRow("SELECT * FROM secret WHERE id = ? AND expired_at > ?", id, time.Now().UTC().Format(formatDate))

	var secret types.Secret
	err := res.Scan(&secret.ID, &secret.Content, &secret.CustomPwd, &secret.CreatedAt, &secret.ExpiredAt)

	if err != nil {
		return types.Secret{}, err
	}

	return secret, nil
}

func (r *MySqlSecretRepository) CreateSecret(content string, customPwd bool, expire time.Time) (types.Secret, error) {

	u := uuid.Must(uuid.NewV4(), nil)
	id := u.String()

	secret := types.Secret{ID: id, Content: content, CustomPwd: customPwd, CreatedAt: time.Now().UTC(), ExpiredAt: expire}

	_, err := r.SQL.Exec("INSERT INTO secret (id, content, custom_pwd, created_at, expired_at) VALUES (?, ?, ?, ?, ?)", secret.ID, secret.Content, secret.CustomPwd, secret.CreatedAt.Format(formatDate), secret.ExpiredAt.Format(formatDate))

	if err != nil {
		return types.Secret{}, err
	}

	return secret, nil
}

func (r *MySqlSecretRepository) RemoveSecret(id string) error {

	_, err := r.SQL.Exec("DELETE FROM secret WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *MySqlSecretRepository) HasSecretWithCustomPwd(id string) (bool, error) {

	secret, err := r.GetSecret(id)
	if err != nil {
		return false, err
	}

	return secret.CustomPwd, nil
}

func (r *MySqlSecretRepository) RemoveSecretsExpired() (int64, error) {

	re, err := r.SQL.Exec("DELETE FROM secret WHERE expired_at <= ?", time.Now().UTC().Format(formatDate))
	if err != nil {
		return 0, err
	}

	return re.RowsAffected()
}
