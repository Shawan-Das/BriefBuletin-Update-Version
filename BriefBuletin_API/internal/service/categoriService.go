package service

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	auth "github.com/rest/api/internal/dbmodel/db_query"
	"github.com/rest/api/internal/model"
)

func (o *OpenAPIService) getAllCategory(c *gin.Context) APIResponse {
	ctx := context.Background()
	db := o.dbConn.GetPool()
	qtx := auth.New(db)

	data, err := qtx.GetAllCategory(ctx)
	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Category List", data)
}

func (s *RESTService) createCategory(c *gin.Context) APIResponse {
	user, _ := GetLoggedInUser(c)
	if user.Role != "ADMIN" {
		return BuildResponse400("Permission denied")
	}
	var input model.CreateCategory

	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	if input.Name == "" || input.Slug == "" {
		return BuildResponse400("Category Name/Slug is missing")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)
	// check for name duplicate
	check, err := qtx.CheckCategorySlugExists(ctx, auth.CheckCategorySlugExistsParams{
		Slug: input.Slug,
		Name: input.Name,
	})
	if check || err == nil {
		return BuildResponse400("Category already exists")
	}

	data, err := qtx.CreateCategory(ctx, auth.CreateCategoryParams{
		Name: input.Name,
		Slug: input.Slug,
	})

	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Created", data)
}

func (s *RESTService) updateCategory(c *gin.Context) APIResponse {
	var input model.UpdateCategory
	user, _ := GetLoggedInUser(c)
	if user.Role != "ADMIN" {
		return BuildResponse400("Permission denied")
	}

	if !parseInput(c, &input) {
		return BuildResponse400("Invalid input provided")
	}

	if input.Id == 0 || input.Name == "" || input.Slug == "" {
		return BuildResponse400("Category Id/Name/Slug is missing")
	}

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)
	// check for name duplicate
	check, err := qtx.CheckCategoryExists(ctx, input.Id)
	if !check || err == nil {
		return BuildResponse400("Category doesn't exists")
	}

	data, err := qtx.UpdateCategory(ctx, auth.UpdateCategoryParams{
		Name: input.Name,
		Slug: input.Slug,
		ID:   input.Id,
	})

	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}
	return BuildResponse200("Updated", data)
}

func (s *RESTService) deleteCategory(c *gin.Context) APIResponse {
	user, _ := GetLoggedInUser(c)
	if user.Role != "ADMIN" {
		return BuildResponse400("Permission denied")
	}
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)

	ctx := context.Background()
	db := s.dbConn.GetPool()
	qtx := auth.New(db)
	// check for name duplicate
	check, err := qtx.CheckCategoryExists(ctx, int32(id))
	if !check || err == nil {
		return BuildResponse400("Category doesn't exists")
	}
	// check category in use
	exist, err := qtx.CheckCategoryInUse(ctx, int32(id))
	if exist || err == nil {
		return BuildResponse400("Category in use. Not allowed to delete")
	}

	err = qtx.DeleteCategory(ctx, int32(id))
	if err != nil {
		return BuildResponse500("Something went wrong", err.Error())
	}

	return BuildResponse200("Deleted", nil)
}
