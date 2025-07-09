package rest

import (
	"fmt"
	"net/http"
	"platform-exercise/internal/config"
	"platform-exercise/internal/entities"
	"platform-exercise/internal/middlewares"
	"platform-exercise/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserRoutes struct {
	cfg         *config.Config
	userService *services.UserService
}

func NewUserRoutes(cfg *config.Config, userService *services.UserService) *UserRoutes {
	return &UserRoutes{
		cfg:         cfg,
		userService: userService,
	}
}

func (u *UserRoutes) ImportRoutes(e *gin.Engine) {
	users := e.Group("/v1/users")
	{
		users.POST("/", u.createUser)
		users.POST("/login", u.logInUser)
	}
	protected := e.Group("/v1/users/")
	protected.Use(middlewares.AuthMiddleware(u.userService.AuthService))
	{
		protected.GET("/:id", u.getUser)
		protected.PATCH("/:id", u.updateUser)
		protected.DELETE("/:id", u.deleteUser)
	}
}

// @id      createUser
// @Summary Create a new user
// @Tags    user
// @Accept  json
// @Produce json
// @Param   request body CreateUserPayload true "The user payload"
// @Success 201 {object} UserResponse "New user created successfully"
// @Failure 400 {object} UserResponse "Payload schema validation failed"
// @Failure 500 {object} UserResponse "Unexpected internal server errors"
// @Router  /v1/users [post]
func (r *UserRoutes) createUser(c *gin.Context) {
	ctx := c.Request.Context()

	var input CreateUserPayload

	var response UserResponse
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorAsString := err.Error()
		response.Error = &errorAsString
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := r.userService.NewUser(ctx, &entities.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Birthday: input.Birthday,
	})

	if err != nil {
		errorAsString := err.Error()
		if strings.HasPrefix(errorAsString, "invalid password") {
			response.Error = &errorAsString
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		log.Warn().Msgf("error:%v", err)
		serviceUnavailable := "service unavailable"
		response.Error = &serviceUnavailable
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = newUser
	c.JSON(http.StatusCreated, response)
}

// @id      createUser
// @Summary Create a new user
// @Tags    user
// @Accept  json
// @Produce json
// @Success 202 {object} UserResponse "New user created successfully"
// @Failure 401 {object} UserResponse "Unauthorized"
// @Failure 500 {object} UserResponse "Unexpected internal server errors"
// @Router  /v1/users/:id [get]
func (r *UserRoutes) getUser(c *gin.Context) {
	ctx := c.Request.Context()

	var response UserResponse
	resourceId := c.Param("id")
	userId, _ := c.Get("userid")
	if resourceId != userId {
		forbidden := "forbidden"
		response.Error = &forbidden
		c.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	user, err := r.userService.GetUser(ctx, resourceId)

	if err != nil {
		log.Warn().Msgf("error:%v", err)
		serviceUnavailable := "service unavailable"
		response.Error = &serviceUnavailable
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	if user == nil {
		notFound := "not found"
		response.Error = &notFound
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	response.Data = user
	c.JSON(http.StatusOK, response)
}

// @id      logInUser
// @Summary Log user in
// @Tags    user
// @Accept  json
// @Produce json
// @Param   request body LoginUserPayload true "The log in user payload"
// @Success 200 {object} LogInUserResponse "successful log user in"
// @Failure 400 {object} LogInUserResponse "Payload schema validation failed"
// @Failure 401 {object} LogInUserResponse "Payload schema validation failed"
// @Failure 500 {object} LogInUserResponse "Unexpected internal server errors"
// @Router  /v1/users [post]
func (r *UserRoutes) logInUser(c *gin.Context) {
	ctx := c.Request.Context()

	var input LogInUserPayload

	var response LogInUserResponse
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorAsString := err.Error()
		response.Error = &errorAsString
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	token, err := r.userService.LogInUser(ctx, &entities.User{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		log.Warn().Msgf("error:%v", err)
		serviceUnavailable := "service unavailable"
		response.Error = &serviceUnavailable
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	if token == "" {
		unauthorized := "unauthorized"
		response.Error = &unauthorized
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	response.Data = map[string]string{
		"token": token,
	}
	c.JSON(http.StatusOK, response)
}

// @id      updateUser
// @Summary update some user properties
// @Tags    user
// @Accept  json
// @Produce json
// @Param   request body UserPatchPayload true "users properties to update"
// @Success 200 {object} UserResponse "successful update user"
// @Failure 400 {object} UserResponse "Payload schema validation failed"
// @Failure 401 {object} UserResponse "Unauthorized"
// @Failure 500 {object} LogInUserResponse "Unexpected internal server errors"
// @Router  /v1/users/:id [patch]
func (r *UserRoutes) updateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var input UserPatchPayload

	var response UserResponse
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorAsString := err.Error()
		response.Error = &errorAsString
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	resourceId := c.Param("id")
	userId, _ := c.Get("userid")
	if resourceId != userId {
		forbidden := "forbidden"
		response.Error = &forbidden
		c.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	user, err := r.userService.UpdateUser(ctx, &entities.User{
		ID:       &resourceId,
		Name:     input.Name,
		Birthday: input.Birthday,
	})

	if err != nil {
		log.Warn().Msgf("error:%v", err)
		serviceUnavailable := "service unavailable"
		response.Error = &serviceUnavailable
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	if user == nil {
		notFound := "not found"
		response.Error = &notFound
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	response.Data = user
	c.JSON(http.StatusOK, response)
}

// @id      deleteUser
// @Summary delete user
// @Tags    user
// @Accept  json
// @Produce json
// @Success 204 {object} UserResponse "successful deleted user"
// @Failure 400 {object} UserResponse "Payload schema validation failed"
// @Failure 401 {object} UserResponse "Unauthorized"
// @Failure 500 {object} LogInUserResponse "Unexpected internal server errors"
// @Router  /v1/users/:id [delete]
func (r *UserRoutes) deleteUser(c *gin.Context) {
	ctx := c.Request.Context()

	var response UserResponse
	resourceId := c.Param("id")
	userId, _ := c.Get("userid")

	if resourceId != userId {
		forbidden := "forbidden"
		response.Error = &forbidden
		c.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}
	err := r.userService.DeleteUser(ctx, resourceId)

	if err != nil {
		fmt.Println("err:", err)
		if strings.Contains(err.Error(), "record not found") {
			notFound := "not found"
			response.Error = &notFound
			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return
		}
		log.Warn().Msgf("error:%v", err)
		serviceUnavailable := "service unavailable"
		response.Error = &serviceUnavailable
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusNoContent, response)
}
