package usecases_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/AlecSmith96/faceit-user-service/internal/usecases"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Updating a user", func() {
	var w *httptest.ResponseRecorder
	var requestBody *usecases.UpdateUserRequestBody
	var userID string

	var updateUserResponse *entities.User
	var updateUserErr error
	var updateUserCallCount int

	var changelogWriterErr error
	var changelogWriterCallCount int

	BeforeEach(func() {
		requestBody = &usecases.UpdateUserRequestBody{
			FirstName: "alec",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "some-password",
			Email:     "alec@email.com",
			Country:   "UK",
		}

		userID = uuid.New().String()

		updateUserResponse = &entities.User{
			ID:        uuid.New(),
			FirstName: "alec",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "some-password",
			Email:     "alec@email.com",
			Country:   "UK",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		updateUserErr = nil
		updateUserCallCount = 1

		changelogWriterErr = nil
		changelogWriterCallCount = 1
	})

	JustBeforeEach(func() {
		w = httptest.NewRecorder()
		requestBodyJSON, err := json.Marshal(requestBody)
		Expect(err).ToNot(HaveOccurred())

		mockUserUpdater.EXPECT().UpdateUser(
			gomock.AssignableToTypeOf(ctxType),
			gomock.AssignableToTypeOf(uuid.UUID{}),
			requestBody.FirstName,
			requestBody.LastName,
			requestBody.Nickname,
			requestBody.Password,
			requestBody.Email,
			requestBody.Country,
		).Return(updateUserResponse, updateUserErr).Times(updateUserCallCount)

		mockChangelogWriter.EXPECT().PublishChangelogEntry(entities.ChangelogEntry{
			UserID:     updateUserResponse.ID,
			CreatedAt:  updateUserResponse.CreatedAt,
			ChangeType: "PUT",
		}).Return(changelogWriterErr).Times(changelogWriterCallCount)

		req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/user/%s", userID), bytes.NewReader(requestBodyJSON))
		Expect(err).ToNot(HaveOccurred())
		r.ServeHTTP(w, req)
	})

	It("should return the updated user", func() {
		Expect(w.Code).To(Equal(http.StatusOK))
		var user usecases.CreateUserResponseBody
		err := json.NewDecoder(w.Body).Decode(&user)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.ID).To(Equal(updateUserResponse.ID.String()))
		Expect(user.FirstName).To(Equal(updateUserResponse.FirstName))
		Expect(user.LastName).To(Equal(updateUserResponse.LastName))
		Expect(user.Nickname).To(Equal(updateUserResponse.Nickname))
		Expect(user.Password).To(Equal(updateUserResponse.Password))
		Expect(user.Email).To(Equal(updateUserResponse.Email))
		Expect(user.Country).To(Equal(updateUserResponse.Country))
		Expect(user.CreatedAt.UTC()).To(Equal(updateUserResponse.CreatedAt.UTC()))
		Expect(user.UpdatedAt.UTC()).To(Equal(updateUserResponse.UpdatedAt.UTC()))
	})

	When("the userID isnt a valid uuid", func() {
		BeforeEach(func() {
			userID = "invalid-uuid"
			requestBody = &usecases.UpdateUserRequestBody{}
			updateUserCallCount = 0
			changelogWriterCallCount = 0
		})

		It("should return a 400 Bad Request", func() {
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("the request fails to validate", func() {
		BeforeEach(func() {
			requestBody = &usecases.UpdateUserRequestBody{}
			updateUserCallCount = 0
			changelogWriterCallCount = 0
		})

		It("should return a 400 Bad Request", func() {
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("the userUpdater adapter returns ErrUserNotFound", func() {
		BeforeEach(func() {
			updateUserErr = entities.ErrUserNotFound
			changelogWriterCallCount = 0
		})

		It("should return a 400 Bad Request", func() {
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("the userUpdater adapter returns generic error", func() {
		BeforeEach(func() {
			updateUserErr = errors.New("an error occurred")
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
