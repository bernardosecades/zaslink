package repositories

import (
	"database/sql"
	"fmt"
	"github.com/bernardosecades/sharesecret/models"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

type MySqlSecretRepository struct {
	SQL *sql.DB
}

func NewMySqlSecretRepository(dbName string, dbUser string, dbPass string, dbHost string, dbPort string) SecretRepository {
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
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

func (repository *MySqlSecretRepository) GetSecret(id string) (models.Secret, error) {

	res := repository.SQL.QueryRow("SELECT * FROM secret WHERE id = ? AND viewed = 0", id)

	var secret models.Secret
	err := res.Scan(&secret.Id, &secret.Content, &secret.Viewed, &secret.CustomPwd)

	if err != nil {
		return models.Secret{}, err
	}

	return secret, nil
}

func (repository *MySqlSecretRepository) CreateSecret(content string, customPwd bool) (models.Secret, error) {

	u := uuid.Must(uuid.NewV4(), nil)
	id := u.String()

	secret := models.Secret{Id: id, Content: content, Viewed: false, CustomPwd: customPwd}

	_, err := repository.SQL.Exec("INSERT INTO secret (id, content, viewed, custom_pwd) VALUES (?, ?, ?, ?)", secret.Id, secret.Content, secret.Viewed, secret.CustomPwd)

	if err != nil {
		return models.Secret{}, err
	}

	return secret, nil
}

func (repository *MySqlSecretRepository) UpdateToViewed(id string) error {

	_, err := repository.SQL.Exec("UPDATE secret SET viewed = 1 WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *MySqlSecretRepository) HasSecretWithCustomPwd(id string) (bool, error) {

	secret, err := repository.GetSecret(id)
	if err != nil {
		return false, err
	}

	return secret.CustomPwd, nil
}
