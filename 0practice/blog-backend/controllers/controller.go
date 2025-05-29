package controllers

import (
	"blog-backend/config"
	"blog-backend/models"
	"blog-backend/utils"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment 发表评论
func CreateComment(c *gin.Context) {
	//测试函数是否正确进入
	// utils.Logger.Infof("Into CreateComment func!")

	// 获取用户 ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		utils.Logger.Errorf("User not authenticated")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// 获取并验证文章 ID
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	// 接收输入
	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建评论对象
	comment := models.Comment{
		Content: input.Content,
		UserID:  userID,
		PostID:  uint(postID),
	}

	// 插入数据库
	if err := config.DB.Create(&comment).Error; err != nil {
		utils.Logger.Errorf("Failed to create comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 返回结果
	c.JSON(http.StatusCreated, comment)
}

// GetCommentsByPost 获取某篇文章的所有评论
func GetCommentsByPost(c *gin.Context) {
	postIDStr := c.Param("id")
	var comments []models.Comment
	if err := config.DB.Preload("User").Where("post_id = ?", postIDStr).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}
