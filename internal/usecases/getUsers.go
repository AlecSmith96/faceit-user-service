package usecases

import (
	"context"
	"github.com/AlecSmith96/faceit-user-service/internal/entities"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type UserGetter interface {
	GetPaginatedUsers(ctx context.Context, firstName, lastName, nickname, email, country string, pageInfo entities.PageInfo) ([]entities.User, string, error)
}

// GetUsersRequestBody represents the request body for getting users
// @Description Optional search criteria for getting users
type GetUsersRequestBody struct {
	// FirstName represents the user's first name
	FirstName string `json:"first_name"`
	// LastName represents the user's last name
	LastName string `json:"last_name"`
	// Nickname represents the user's nickname
	Nickname string `json:"nickname"`
	// Email represents the user's email address
	Email string `json:"email"`
	// Country represents the user's country
	Country string `json:"country"`
	// PageInfo represents the pagination information for the request
	PageInfo PageInfo `json:"page_info"`
}

// GetUsersResponseBody represents the response body for getting users
// @Description List of users matching search criteria and pagination info
type GetUsersResponseBody struct {
	// Users represents the users matching the search criteria
	Users []UserResponse `json:"users"`
	// PageInfo represents the pagination information for the request
	PageInfo PageInfo `json:"page_info"`
}

// PageInfo represents the pagination info for a request
// @Description Provides page size and the token used to get the next page of users
type PageInfo struct {
	// NextPageToken represents the token used to get the next page of results
	NextPageToken string `json:"next_page_token"`
	// PageSize represents the number of results per page, default is 10
	PageSize int `json:"page_size"`
}

// UserResponse represents the response body of a user
// @Description Information of an individual user in the list of users
type UserResponse struct {
	// ID represents the user's unique identifier
	ID string `json:"id"`
	// FirstName represents the user's first name
	FirstName string `json:"first_name"`
	// LastName represents the user's last name
	LastName string `json:"last_name"`
	// Nickname represents the user's nickname
	Nickname string `json:"nickname"`
	// Password represents the user's password
	Password string `json:"password"`
	// Email represents the user's email address
	Email string `json:"email"`
	// Country represents the user's country
	Country string `json:"country"`
	// CreatedAt represents the timestamp when the user was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt represents the timestamp when the user was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// NewGetUsers Get Users
// @Summary Get a list of users
// @Description Gets  list of users based on optional search criteria
// @Tags users
// @Accept json
// @Produce json
// @Param user body GetUsersRequestBody true "Get Users Request Body"
// @Success 200 {object} GetUsersResponseBody
// @Failure 400
// @Failure 500
// @Router /users [get]
func NewGetUsers(userGetter UserGetter, changelogWriter ChangelogWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request GetUsersRequestBody
		err := c.ShouldBindJSON(&request)
		if err != nil {
			slog.Warn("unable to bind request", "err", err)
			c.Status(http.StatusBadRequest)
			return
		}

		if request.PageInfo.PageSize == 0 {
			request.PageInfo.PageSize = 10
		}

		pageInfo := entities.PageInfo{
			NextPageToken: request.PageInfo.NextPageToken,
			PageSize:      request.PageInfo.PageSize,
		}

		users, nextPageToken, err := userGetter.GetPaginatedUsers(
			c.Request.Context(),
			request.FirstName,
			request.LastName,
			request.Nickname,
			request.Email,
			request.Country,
			pageInfo,
		)
		if err != nil {
			slog.Error("getting paginated users", "err", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		usersResponse := make([]UserResponse, 0)
		for _, user := range users {
			usersResponse = append(usersResponse, UserResponse{
				ID:        user.ID.String(),
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Nickname:  user.Nickname,
				Password:  user.Password,
				Email:     user.Email,
				Country:   user.Country,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			})
		}

		response := GetUsersResponseBody{
			Users: usersResponse,
			PageInfo: PageInfo{
				NextPageToken: nextPageToken,
				PageSize:      request.PageInfo.PageSize,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}
