package services

import (
	"github.com/mahendraintelops/mahen-project/test-service/pkg/rest/server/daos"
	"github.com/mahendraintelops/mahen-project/test-service/pkg/rest/server/models"
)

type TestService struct {
	testDao *daos.TestDao
}

func NewTestService() (*TestService, error) {
	testDao, err := daos.NewTestDao()
	if err != nil {
		return nil, err
	}
	return &TestService{
		testDao: testDao,
	}, nil
}

func (testService *TestService) CreateTest(test *models.Test) (*models.Test, error) {
	return testService.testDao.CreateTest(test)
}

func (testService *TestService) ListTests() ([]*models.Test, error) {
	return testService.testDao.ListTests()
}

func (testService *TestService) GetTest(id int64) (*models.Test, error) {
	return testService.testDao.GetTest(id)
}
