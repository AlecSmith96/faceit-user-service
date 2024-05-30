package usecases_test

import (
	"bytes"
	"errors"
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

var _ = Describe("Creating a user", func() {
	var w *httptest.ResponseRecorder
	var requestBody *usecases.CreateUserRequestBody

	var createUserResponse *entities.User
	var createUserErr error
	var createUserCallCount int

	var changelogWriterErr error
	var changelogWriterCallCount int

	BeforeEach(func() {
		requestBody = &usecases.CreateUserRequestBody{
			FirstName: "alec",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Password:  "some-password",
			Email:     "alec@email.com",
			Country:   "UK",
		}

		createUserResponse = &entities.User{
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

		createUserErr = nil
		createUserCallCount = 1

		changelogWriterErr = nil
		changelogWriterCallCount = 1
	})

	JustBeforeEach(func() {
		w = httptest.NewRecorder()
		requestBodyJSON, err := json.Marshal(requestBody)
		Expect(err).ToNot(HaveOccurred())

		mockUserCreator.EXPECT().CreateUser(
			gomock.AssignableToTypeOf(ctxType),
			requestBody.FirstName,
			requestBody.LastName,
			requestBody.Nickname,
			requestBody.Password,
			requestBody.Email,
			requestBody.Country,
		).Return(createUserResponse, createUserErr).Times(createUserCallCount)

		mockChangelogWriter.EXPECT().PublishChangelogEntry(entities.ChangelogEntry{
			UserID:     createUserResponse.ID,
			CreatedAt:  createUserResponse.CreatedAt,
			ChangeType: "POST",
		}).Return(changelogWriterErr).Times(changelogWriterCallCount)

		req, err := http.NewRequest("POST", "http://localhost:8080/user", bytes.NewReader(requestBodyJSON))
		Expect(err).ToNot(HaveOccurred())
		r.ServeHTTP(w, req)
	})

	It("should return the created user", func() {
		Expect(w.Code).To(Equal(http.StatusOK))
		var user usecases.CreateUserResponseBody
		err := json.NewDecoder(w.Body).Decode(&user)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.ID).To(Equal(createUserResponse.ID.String()))
		Expect(user.FirstName).To(Equal(createUserResponse.FirstName))
		Expect(user.LastName).To(Equal(createUserResponse.LastName))
		Expect(user.Nickname).To(Equal(createUserResponse.Nickname))
		Expect(user.Password).To(Equal(createUserResponse.Password))
		Expect(user.Email).To(Equal(createUserResponse.Email))
		Expect(user.Country).To(Equal(createUserResponse.Country))
		Expect(user.CreatedAt.UTC()).To(Equal(createUserResponse.CreatedAt.UTC()))
		Expect(user.UpdatedAt.UTC()).To(Equal(createUserResponse.UpdatedAt.UTC()))
	})

	When("the request fails to validate", func() {
		BeforeEach(func() {
			requestBody = &usecases.CreateUserRequestBody{}
			createUserCallCount = 0
			changelogWriterCallCount = 0
		})

		It("should return a 400 Bad Request", func() {
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("the userCreator adapter returns ErrEmailAlreadyUsed", func() {
		BeforeEach(func() {
			createUserErr = entities.ErrEmailAlreadyUsed
			changelogWriterCallCount = 0
		})

		It("should return a 400 Bad Request", func() {
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("the userCreator adapter returns generic error", func() {
		BeforeEach(func() {
			createUserErr = errors.New("an error occurred")
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
