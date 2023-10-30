package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mahendraintelops/mahen-project/test-service/pkg/rest/server/daos/clients/sqls"
	"github.com/mahendraintelops/mahen-project/test-service/pkg/rest/server/models"
	"github.com/mahendraintelops/mahen-project/test-service/pkg/rest/server/services"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	"strconv"
)

type TestController struct {
	testService *services.TestService
}

func NewTestController() (*TestController, error) {
	testService, err := services.NewTestService()
	if err != nil {
		return nil, err
	}
	return &TestController{
		testService: testService,
	}, nil
}

func (testController *TestController) CreateTest(context *gin.Context) {
	// validate input
	var input models.Test
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// trigger test creation
	testCreated, err := testController.testService.CreateTest(&input)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, testCreated)
}

func (testController *TestController) ListTests(context *gin.Context) {
	// trigger all tests fetching
	tests, err := testController.testService.ListTests()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, tests)
}

func (testController *TestController) FetchTest(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger test fetching
	test, err := testController.testService.GetTest(id)
	if err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(context.Request.Context())
		currentSpan.SetAttributes(attribute.String("test.id", strconv.FormatInt(test.Id, 10)))
	}

	context.JSON(http.StatusOK, test)
}
