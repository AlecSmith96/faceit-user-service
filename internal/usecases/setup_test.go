package usecases_test

import (
	"context"
	"github.com/AlecSmith96/faceit-user-service/internal/drivers"
	mock_usecases "github.com/AlecSmith96/faceit-user-service/mocks"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
	"net/http"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func TestHandleUsers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Users Test Suite")
}

var (
	r                    *gin.Engine
	mockChangelogWriter  *mock_usecases.MockChangelogWriter
	mockUserCreator      *mock_usecases.MockUserCreator
	mockUserUpdater      *mock_usecases.MockUserUpdater
	mockUserDeleter      *mock_usecases.MockUserDeleter
	mockUserGetter       *mock_usecases.MockUserGetter
	mockReadinessChecker *mock_usecases.MockReadinessChecker
)

var _ = BeforeSuite(func() {
	// Put gin in test mode
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(GinkgoT())
	mockChangelogWriter = mock_usecases.NewMockChangelogWriter(ctrl)
	mockUserCreator = mock_usecases.NewMockUserCreator(ctrl)
	mockUserUpdater = mock_usecases.NewMockUserUpdater(ctrl)
	mockUserDeleter = mock_usecases.NewMockUserDeleter(ctrl)
	mockUserGetter = mock_usecases.NewMockUserGetter(ctrl)
	mockReadinessChecker = mock_usecases.NewMockReadinessChecker(ctrl)

	r = drivers.NewRouter(
		mockChangelogWriter,
		mockUserGetter,
		mockUserCreator,
		mockUserDeleter,
		mockUserUpdater,
		mockReadinessChecker,
	)

	go func() {
		defer GinkgoRecover()
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}
	}()
})
