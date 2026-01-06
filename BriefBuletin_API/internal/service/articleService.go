package service

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	auth "github.com/rest/api/internal/dbmodel/db_query"
)

func (o *OpenAPIService) getApprovedArticles(c *gin.Context) APIResponse {
	p := c.Query("page")
	lang := c.Query("lang")
	totalNews := c.Query("totalNews")

	limit := 10
	page, _ := strconv.Atoi(p)
	prevNewsCount, _ := strconv.Atoi(totalNews)

	ctx := context.Background()
	db := o.dbConn.GetPool()
	qtx := auth.New(db)

	// first try to fetch today's articles
	zone := time.FixedZone("BST", 6*60*60)
	now := time.Now().In(zone)
	today := now.Format("2006-01-02")
	var data []auth.NewsArticle
	var err error
	if lang == "bn" && page == 1 {
		data, _ = qtx.GetTodaysBnNews(ctx, StringToPgDate(today))
	} else if lang == "en" && page == 1 {
		data, _ = qtx.GetTodaysEnNews(ctx, StringToPgDate(today))
	}
	if len(data) != 0 {
		return BuildResponse200("Today's Article List", data)
	}

	if lang == "bn" {
		data, err = qtx.GetApprovedBanArticles(ctx, auth.GetApprovedBanArticlesParams{
			Offset: int32(prevNewsCount), //int32((page - 1) * limit),
			Limit:  int32(limit),
		})
		if err != nil {
			return BuildResponse500("Something went wrong", err.Error())
		}
		return BuildResponse200("Article List", data)
	}

	data, err = qtx.GetApprovedEngArticles(ctx, auth.GetApprovedEngArticlesParams{
		Offset: int32(prevNewsCount), //int32((page - 1) * limit),
		Limit:  int32(limit),
	})
	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Article List", data)
}

func (o *OpenAPIService) readNews(c *gin.Context) APIResponse {
	p := c.Query("id")
	id, _ := strconv.Atoi(p)
	ctx := context.Background()
	db := o.dbConn.GetPool()
	qtx := auth.New(db)

	err := qtx.ReadArticleCount(ctx, int32(id))
	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Article List", nil)
}

func (s *RESTService) unApprovedArticleList(c *gin.Context) APIResponse {
	_ = c.Query("page")
	// limit := 12
	// page, _ := strconv.Atoi(p)

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	// data, err := qtx.GetUnApprovedArticleList(ctx, auth.GetUnApprovedArticleListParams{
	// 	Offset: int32((page - 1) * limit),
	// 	Limit:  int32(limit),
	// })
	data, err := qtx.GetUnApprovedArticleList(ctx)
	if len(data) == 0 {
		return BuildResponse200("No article in draft list", true)
	}
	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Draft article list", data)
}

func (s *RESTService) approveArticle(c *gin.Context) APIResponse {
	user, _ := GetLoggedInUser(c)
	if user.Role != "ADMIN" && user.Role != "editor" {
		return BuildResponse400("Permission denied")
	}
	id := c.Query("id")
	if id == "" {
		return BuildResponse400("Article id is missing")
	}
	article_id, _ := strconv.Atoi(id)

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	err := qtx.ApproveArticle(ctx, int32(article_id))

	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Article Approved", nil)
}

func (s *RESTService) draftArticle(c *gin.Context) APIResponse {
	user, _ := GetLoggedInUser(c)
	if user.Role != "ADMIN" && user.Role != "editor" {
		return BuildResponse400("Permission denied")
	}
	id := c.Query("id")
	if id == "" {
		return BuildResponse400("Article id is missing")
	}
	article_id, _ := strconv.Atoi(id)

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)

	err := qtx.DraftArticle(ctx, int32(article_id))

	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Article drafted", nil)
}
