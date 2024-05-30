package usecases_test

import (
	"errors"
	"fmt"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Deleting a user", func() {
	var w *httptest.ResponseRecorder

	var userID string

	var deleteUserErr error
	var deleteUserCallCount int

	var changelogWriterErr error
	var changelogWriterCallCount int

	BeforeEach(func() {
		userID = uuid.New().String()

		deleteUserErr = nil
		deleteUserCallCount = 1

		changelogWriterErr = nil
		changelogWriterCallCount = 1
	})

	JustBeforeEach(func() {
		w = httptest.NewRecorder()

		mockUserDeleter.EXPECT().DeleteUser(
			gomock.AssignableToTypeOf(ctxType),
			gomock.AssignableToTypeOf(uuid.UUID{}),
		).Return(deleteUserErr).Times(deleteUserCallCount)

		mockChangelogWriter.EXPECT().PublishChangelogEntry(gomock.AssignableToTypeOf(entities.ChangelogEntry{})).
			Return(changelogWriterErr).Times(changelogWriterCallCount)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/user/%s", userID), nil)
		Expect(err).ToNot(HaveOccurred())
		r.ServeHTTP(w, req)
	})

	It("should return the created user", func() {
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	When("the request fails to validate", func() {
		BeforeEach(func() {
			userID = "invalid-uuid"
			deleteUserCallCount = 0
			changelogWriterCallCount = 0
		})

		It("should return a 400 Bad Request", func() {
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("the userDeleter adapter returns ErrUserNotFound", func() {
		BeforeEach(func() {
			deleteUserErr = entities.ErrUserNotFound
			changelogWriterCallCount = 0
		})

		It("should return a 400 Bad Request", func() {
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("the userDeleter adapter returns generic error", func() {
		BeforeEach(func() {
			deleteUserErr = errors.New("an error occurred")
			changelogWriterCallCount = 0
		})

		It("should return a 500 Internal Server Error", func() {
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	When("the changelog fails to send", func() {
		BeforeEach(func() {
			changelogWriterErr = errors.New("an error occurred")
		})

		It("should return a 500 Internal Server Error", func() {
			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})
})
