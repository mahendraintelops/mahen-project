package daos

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/mahendraintelops/mahen-project/test-service/pkg/rest/server/daos/clients/sqls"
	"github.com/mahendraintelops/mahen-project/test-service/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
)

type TestDao struct {
	sqlClient *sqls.MySQLClient
}

func migrateTests(r *sqls.MySQLClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS tests(
		ID int NOT NULL AUTO_INCREMENT,
        
		Name VARCHAR(100) NOT NULL,
	    PRIMARY KEY (ID)
	);
	`
	_, err := r.DB.Exec(query)
	return err
}

func NewTestDao() (*TestDao, error) {
	sqlClient, err := sqls.InitMySQLDB()
	if err != nil {
		return nil, err
	}
	err = migrateTests(sqlClient)
	if err != nil {
		return nil, err
	}
	return &TestDao{
		sqlClient,
	}, nil
}

func (testDao *TestDao) CreateTest(m *models.Test) (*models.Test, error) {
	insertQuery := "INSERT INTO tests(Name) values(?)"
	res, err := testDao.sqlClient.DB.Exec(insertQuery, m.Name)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return nil, sqls.ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.Id = id
	log.Debugf("test created")
	return m, nil
}

func (testDao *TestDao) ListTests() ([]*models.Test, error) {
	selectQuery := "SELECT * FROM tests"
	rows, err := testDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var tests []*models.Test
	for rows.Next() {
		m := models.Test{}
		if err = rows.Scan(&m.Id, &m.Name); err != nil {
			return nil, err
		}
		tests = append(tests, &m)
	}
	if tests == nil {
		tests = []*models.Test{}
	}
	log.Debugf("test listed")
	return tests, nil
}

func (testDao *TestDao) GetTest(id int64) (*models.Test, error) {
	selectQuery := "SELECT * FROM tests WHERE Id = ?"
	row := testDao.sqlClient.DB.QueryRow(selectQuery, id)

	m := models.Test{}
	if err := row.Scan(&m.Id, &m.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}
	log.Debugf("test retrieved")
	return &m, nil
}
