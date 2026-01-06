package service

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	auth "github.com/rest/api/internal/dbmodel/db_query"
	"github.com/rest/api/internal/model"
)

func (o *OpenAPIService) createComment(c *gin.Context) APIResponse {
	var i model.CreateComment
	if !parseInput(c, &i) {
		return BuildResponse400("Invalid input provided")
	}

	if i.ArticleID == 0 {
		return BuildResponse400("Article id is missing")
	}
	if i.UserName == "" || i.UserEmail == "" {
		return BuildResponse400("User name/email is missing")
	}

	ctx := context.Background()
	db := o.dbConn.GetPool()
	qtx := auth.New(db)

	data, err := qtx.CreateCommentWithDefaults(ctx, auth.CreateCommentWithDefaultsParams{
		ArticleID: i.ArticleID,
		UserName:  i.UserName,
		UserEmail: i.UserEmail,
		Content:   i.Content,
	})
	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Category List", data)
}

func (o *OpenAPIService) getNewsComments(c *gin.Context) APIResponse {
	id := c.Query("article_id")
	if id == "" {
		return BuildResponse400("Article id is missing")
	}
	articleID, _ := strconv.Atoi(id)

	ctx := context.Background()
	db := o.dbConn.GetPool()
	qtx := auth.New(db)

	data, err := qtx.GetArticleDetails(ctx, int32(articleID))
	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	} else if data.ID == 0 {
		return BuildResponse400("Article not found")
	}
	comments, err := qtx.GetApprovedCommentsByArticle(ctx, int32(articleID))
	if len(comments) == 0 {
		return BuildResponse200("No comments found for this article", comments)
	} else if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Comments for the article", comments)
}

func (s *RESTService) approveComment(c *gin.Context) APIResponse {
	user, _ := GetLoggedInUser(c)
	if user.Role != "ADMIN" && user.Role != "EDITOR" {
		return BuildResponse400("Permission denied")
	}
	commentId := c.Query("comment_id")
	id, _ := strconv.Atoi(commentId)

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	_, err := qtx.ApproveComment(ctx, int32(id))

	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Done!", nil)
}

func (s *RESTService) disableComment(c *gin.Context) APIResponse {
	user, _ := GetLoggedInUser(c)
	if user.Role != "ADMIN" && user.Role != "EDITOR" {
		return BuildResponse400("Permission denied")
	}
	commentId := c.Query("comment_id")
	id, _ := strconv.Atoi(commentId)

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	_, err := qtx.DisableComment(ctx, int32(id))

	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Done!", nil)
}

func (s *RESTService) approvalDueCommentList(c *gin.Context) APIResponse {
	user, _ := GetLoggedInUser(c)
	if user.Role != "ADMIN" && user.Role != "EDITOR" {
		return BuildResponse400("Permission denied")
	}
	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	data, err := qtx.ApprovalDueListOfComments(ctx)
	if len(data) == 0 {
		return BuildResponse200("No comments due to approve", nil)
	} else if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Comment approval due list", data)
}
