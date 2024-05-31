package usecases_test

import (
	"bytes"
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

var _ = Describe("Getting a list of users", func() {
	var w *httptest.ResponseRecorder
	var requestBody *usecases.GetUsersRequestBody

	var getPaginatedUsersResponse []entities.User
	var nextPageToken string
	var getPaginatedUsersErr error
	var getPaginatedUsersCallCount int

	BeforeEach(func() {
		requestBody = &usecases.GetUsersRequestBody{
			FirstName: "alec",
			LastName:  "smith",
			Nickname:  "alecsmith",
			Email:     "alec@email.com",
			Country:   "UK",
			PageInfo: usecases.PageInfo{
				PageSize: 20,
			},
		}

		getPaginatedUsersResponse = []entities.User{
			{
				ID:        uuid.New(),
				FirstName: "alec",
				LastName:  "smith",
				Nickname:  "alecsmith",
				Password:  "some-password",
				Email:     "alec@email.com",
				Country:   "UK",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				FirstName: "john",
				LastName:  "smith",
				Nickname:  "johnsmith",
				Password:  "some-password",
				Email:     "john@email.com",
				Country:   "UK",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		nextPageToken = "some-page-token"
		getPaginatedUsersErr = nil
		getPaginatedUsersCallCount = 1
	})

	JustBeforeEach(func() {
		w = httptest.NewRecorder()
		requestBodyJSON, err := json.Marshal(requestBody)
		Expect(err).ToNot(HaveOccurred())

		pageInfo := entities.PageInfo{
			NextPageToken: requestBody.PageInfo.NextPageToken,
			PageSize:      requestBody.PageInfo.PageSize,
		}

		mockUserGetter.EXPECT().GetPaginatedUsers(
			gomock.AssignableToTypeOf(ctxType),
			requestBody.FirstName,
			requestBody.LastName,
			requestBody.Nickname,
			requestBody.Email,
			requestBody.Country,
			pageInfo,
		).Return(getPaginatedUsersResponse, nextPageToken, getPaginatedUsersErr).Times(getPaginatedUsersCallCount)

		req, err := http.NewRequest("GET", "http://localhost:8080/users", bytes.NewReader(requestBodyJSON))
		Expect(err).ToNot(HaveOccurred())
		r.ServeHTTP(w, req)
	})

	It("should return the list of users", func() {
		Expect(w.Code).To(Equal(http.StatusOK))
		var resp usecases.GetUsersResponseBody
		err := json.NewDecoder(w.Body).Decode(&resp)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Users).To(HaveLen(2))
		Expect(resp.PageInfo.PageSize).To(Equal(20))
		Expect(resp.PageInfo.NextPageToken).To(Equal("some-page-token"))
	})
})
