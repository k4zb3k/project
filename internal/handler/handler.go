package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/k4zb3k/project/internal/apperror"
	"github.com/k4zb3k/project/internal/models"
	"github.com/k4zb3k/project/internal/service"
	"github.com/k4zb3k/project/pkg/logger"
)

type Handler struct {
	Engine  *gin.Engine
	Service *service.Service
}

func NewHandler(engine *gin.Engine, service *service.Service) *Handler {
	return &Handler{
		Engine:  engine,
		Service: service,
	}
}

// ==============================================

func (h *Handler) InitRoutes() {
	generalRout := h.Engine.Group("v1")

	auth := generalRout.Group("/auth")
	{
		auth.POST("/create", h.CreateUser)
		auth.POST("/login", h.Login)
	}

	api := generalRout.Group("/api")
	api.Use(h.TokenAuthMiddleware())
	{
		api.POST("/account", h.CreateAccount)
		api.GET("/account", h.GetAccounts)
		api.GET("/account/:id", h.GetAccountById)
		api.PUT("/account", h.UpdateAccount) // todo
		api.POST("/transaction", h.CreateTransaction)
		api.GET("/transaction", h.GetTransactions) // todo
	}
}

// ==============================================

func (h *Handler) CreateUser(c *gin.Context) {
	var u *models.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(400, apperror.ErrBadRequest)
		return
	}

	ctx := c.Request.Context()

	err = h.Service.ValidateUser(u)
	if err != nil {
		c.JSON(400, apperror.ErrInvalid)
		return
	}

	existsUser, err := h.Service.ExistsUser(u.Username)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}
	if existsUser {
		c.JSON(400, apperror.ErrRegistered)
		return
	}

	userID, err := h.Service.CreateUser(ctx, u)
	if err != nil {
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	c.JSON(201, map[string]string{
		"user_id": userID,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var u *models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		logger.Error.Println(err)
		c.JSON(400, apperror.ErrBadRequest)
		return
	}

	userID, err := h.Service.CheckUser(u)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(401, apperror.ErrUnauthorized)
		return
	}

	ts, err := h.Service.CreateToken(userID)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	saveErr := h.Service.CreateAuth(userID, ts)
	if saveErr != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	c.JSON(200, tokens)
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var acc *models.Account

	userId, ok := c.Get("user_id")
	if !ok {
		logger.Error.Println("can not get user ID from token")
		c.AbortWithStatus(500) //todo
		return
	}
	userID := userId.(string)

	if err := c.ShouldBindJSON(&acc); err != nil {
		logger.Error.Println(err)
		c.JSON(400, apperror.ErrBadRequest)
		return
	}

	acc.UserID = userID

	// проверка счета пользователя на дубликат
	existsAccount, err := h.Service.ExistsAccount(acc.Number)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}
	if existsAccount {
		c.JSON(400, apperror.ErrExistsAccount)
		return
	}

	// регистрация нового счета пользователя
	err = h.Service.CreateAccount(acc)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	c.JSON(201, "adding new account was successful")
}

func (h *Handler) GetAccounts(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		logger.Error.Println("can not get user ID from token")
		c.AbortWithStatus(500) //todo
		return
	}
	userID := userId.(string)

	accounts, err := h.Service.GetAccounts(userID)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	c.JSON(200, accounts)
}

func (h *Handler) GetAccountById(c *gin.Context) {
	id := c.Param("id")

	userId, ok := c.Get("user_id")
	if !ok {
		logger.Error.Println("can not get user ID from token")
		c.AbortWithStatus(500) //todo
		return
	}
	userID := userId.(string)

	account, err := h.Service.GetAccountById(userID, id)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	c.JSON(200, account)
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	var acc *models.Account

	userId, ok := c.Get("user_id")
	if !ok {
		logger.Error.Println("can not get user ID from token")
		c.AbortWithStatus(500) //todo
		return
	}
	userID := userId.(string)

	acc.UserID = userID

	err := c.ShouldBindJSON(&acc)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(400, apperror.ErrBadRequest)
		return
	}

	err = h.Service.UpdateAccount(acc)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	c.JSON(200, "account was Updated")
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var tr *models.Transaction

	userId, ok := c.Get("user_id")
	if !ok {
		logger.Error.Println("can not get user ID from token")
		c.AbortWithStatus(500) //todo
		return
	}
	userID := userId.(string)

	err := c.ShouldBindJSON(&tr)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(400, apperror.ErrBadRequest)
		return
	}
	ok = tr.Type == "expense" || tr.Type == "income"
	if !ok {
		logger.Error.Println("incorrect transaction type")
		c.JSON(400, apperror.ErrBadRequest)
		return
	}

	account, err := h.Service.GetAccountById(userID, tr.AccountID)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(400, apperror.ErrInternalServer)
		return
	}
	if account.ID == "" {
		logger.Error.Printf("this user doesn't have an account with id %s \n", tr.AccountID)
		c.JSON(400, apperror.ErrBadRequest)
		return
	}

	err = h.Service.CreateTransaction(tr)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	if tr.Type == "expense" { // todo
		account.Balance -= tr.Amount
	} else if tr.Type == "income" {
		account.Balance += tr.Amount
	}

	err = h.Service.UpdateAccount(&account)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	c.JSON(201, "the transaction is saved")
}

func (h *Handler) GetTransactions(c *gin.Context) {
	var tr []models.Transaction

	userId, ok := c.Get("user_id")
	if !ok {
		logger.Error.Println("can not get user ID from token")
		c.AbortWithStatus(500) //todo
		return
	}
	userID := userId.(string)

	accounts, err := h.Service.GetAccounts(userID)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(500, apperror.ErrInternalServer)
		return
	}

	for _, account := range accounts {
		transactions, err := h.Service.GetTransactions(account.ID)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(500, apperror.ErrInternalServer)
			return
		}

		tr = append(tr, transactions...)
	}
}
