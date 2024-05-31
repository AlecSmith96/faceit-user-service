package usecases_test

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("checking service readiness", func() {
	var w *httptest.ResponseRecorder
	var checkConnectionErr error

	BeforeEach(func() {
		checkConnectionErr = nil
	})

	JustBeforeEach(func() {
		w = httptest.NewRecorder()

		mockReadinessChecker.EXPECT().CheckConnection().Return(checkConnectionErr).Times(1)

		req, err := http.NewRequest("GET", "http://localhost:8080/health/readiness", nil)
		Expect(err).ToNot(HaveOccurred())
		r.ServeHTTP(w, req)
	})

	It("should return a 200 OK", func() {
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	When("the database connection cant be established", func() {
		BeforeEach(func() {
			checkConnectionErr = errors.New("an error occurred")
		})

		It("should return a 500 Internal Server Error", func() {
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
