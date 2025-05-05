package porodnica_ambulance_home

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tessyelots/porodnica-webapi/internal/db_service"
)

type PorodnicaHomeSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Porodnica]
}

func TestPorodnicaHomeSuite(t *testing.T) {
	suite.Run(t, new(PorodnicaHomeSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := this.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := this.Called(ctx, id)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := this.Called(ctx)
	return args.Error(0)
}

func (suite *PorodnicaHomeSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Porodnica]{}

	// Compile time Assert that the mock is of type db_service.DbService[Porodnica]
	var _ db_service.DbService[Porodnica] = suite.dbServiceMock

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Porodnica{
				Id: "test-porodnica",
				WaitingList: []WaitingListEntry{
					{
						Id:                 "test-entry",
						PatientId:          "test-patient",
						WaitingSince:       time.Now(),
						EstimatedLaborDate: time.Now().Add(time.Hour),
					},
				},
			},
			nil,
		)
}

func (suite *PorodnicaHomeSuite) Test_UpdateWl_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
        "id": "test-entry",
        "patientId": "test-patient",
        "estimatedDurationMinutes": 42
    }`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "porodnicaId", Value: "test-porodnica"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("POST", "/porodnica/test-porodnica/waitinglist/test-entry", strings.NewReader(json))

	sut := implPorodnicaWaitingListAPI{}

	// ACT
	sut.UpdateWaitingListEntry(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-porodnica", mock.Anything)
}
