package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"server/internal/model"
)

func healthHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		mysqlOK, mysqlErr := pingMySQL(c.Request.Context(), deps.DB)
		redisOK, redisErr := pingRedis(c.Request.Context(), deps.Redis)

		status := http.StatusOK
		if !mysqlOK || !redisOK {
			status = http.StatusServiceUnavailable
		}

		body := gin.H{
			"ok":    mysqlOK && redisOK,
			"mysql": mysqlOK,
			"redis": redisOK,
		}
		if mysqlErr != nil {
			body["mysqlError"] = mysqlErr.Error()
		}
		if redisErr != nil {
			body["redisError"] = redisErr.Error()
		}

		c.JSON(status, body)
	}
}

func pingMySQL(ctx context.Context, db *gorm.DB) (bool, error) {
	if db == nil {
		return false, errors.New("mysql not initialized")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return false, err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func pingRedis(ctx context.Context, rdb *redis.Client) (bool, error) {
	if rdb == nil {
		return false, errors.New("redis not initialized")
	}
	if err := rdb.Ping(ctx).Err(); err != nil {
		return false, err
	}
	return true, nil
}

type createUserReq struct {
	Username string `json:"username" binding:"required"`
}

func createUserHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createUserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte("ChangeMe123!"), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		u := model.User{
			Username:     req.Username,
			Account:      req.Username,
			PasswordHash: string(hash),
			Nickname:     req.Username,
			AvatarURL:    "",
			Bio:          "",
			Role:         "user",
			Status:       "normal",
			BanUntil:     nil,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := deps.DB.WithContext(c.Request.Context()).Create(&u).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, u)
	}
}

func getUserHandler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var u model.User
		if err := deps.DB.WithContext(c.Request.Context()).First(&u, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, u)
	}
}
