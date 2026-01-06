package service

import "github.com/gin-gonic/gin"

func (o *OpenAPIService) AddRouters(router *gin.Engine) {
	router.GET("/api/category", func(c *gin.Context) {
		resp := o.getAllCategory(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.GET("/api/articles", func(c *gin.Context) {
		resp := o.getApprovedArticles(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.GET("/api/read-news", func(c *gin.Context) {
		resp := o.readNews(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.POST("/api/comment", func(c *gin.Context) {
		resp := o.createComment(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.GET("/api/all-comments", func(c *gin.Context) {
		resp := o.getNewsComments(c)
		c.JSON(resp.StatusCode, resp)
	})
}

// AddRouters add api end points specific to this service
func (s *RESTService) AddRouters(router *gin.Engine) {
	router.POST("/api/auth/create", func(c *gin.Context) {
		resp := s.createUser(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.POST("/api/auth/login", func(c *gin.Context) {
		resp := s.validateLogin(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.POST("/api/auth/resetpwd", func(c *gin.Context) {
		resp := s.resetPassword(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.PUT("/api/auth/update", func(c *gin.Context) {
		resp := s.updateUser(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.GET("/api/auth/users", func(c *gin.Context) {
		resp := s.getAllUsers(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.POST("/api/verify-user", func(c *gin.Context) {
		resp := s.userVerify(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.POST("/api/verify-otp", func(c *gin.Context) {
		resp := s.validateOtp(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.POST("/api/send-otp", func(c *gin.Context) {
		resp := s.sendOtp(c)
		c.JSON(resp.StatusCode, resp)
	})
	// category service
	router.POST("/api/category-service", func(c *gin.Context) {
		resp := s.createCategory(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.PUT("/api/category-service", func(c *gin.Context) {
		resp := s.updateCategory(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.DELETE("/api/category-service", func(c *gin.Context) {
		resp := s.deleteCategory(c)
		c.JSON(resp.StatusCode, resp)
	})

	//comment service
	router.GET("/api/active-comment", func(c *gin.Context) {
		resp := s.approveComment(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.GET("/api/disable-comment", func(c *gin.Context) {
		resp := s.disableComment(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.GET("/api/approval-due-comments", func(c *gin.Context) {
		resp := s.approvalDueCommentList(c)
		c.JSON(resp.StatusCode, resp)
	})

	// article service
	router.GET("/api/draft-article-list", func(c *gin.Context) {
		resp := s.unApprovedArticleList(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.GET("/api/publish-article", func(c *gin.Context) {
		resp := s.approveArticle(c)
		c.JSON(resp.StatusCode, resp)
	})
	router.GET("/api/draft-article", func(c *gin.Context) {
		resp := s.draftArticle(c)
		c.JSON(resp.StatusCode, resp)
	})
}
