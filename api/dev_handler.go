package api

import (
	"context"
	"fmt"
	db "github.com/aopontann/gin-sqlc/db/sqlc"
	"github.com/google/uuid"
	"net/http"
	"strconv"

	"github.com/aopontann/gin-sqlc/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateChefData(c *gin.Context) {
	numRecipe, err := strconv.Atoi(c.Query("num_recipe"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// シェフを一人登録
	jsn := []byte(`
		{
			"name": "シェフ名",
			"imageUrl": "8495f92d-e0ad-4015-8283-86cb30d1376c.webp",
			"profile": "シェフです",
			"link": [
				{
					"label": "Facebook",
					"url": "https://www.yahoo.co.jp"
				},
				{
					"label": "X",
					"url": "https://www.yahoo.co.jp"
				}
			]
		}`)
	response, err := s.q.CreateChef(context.Background(), jsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	chefId := utils.UUID2Str(response.ID)
	fmt.Printf("chefId = %s\n", chefId)

	// シェフのレシピをnumRecipe件登録
	for i := 1; i <= numRecipe; i++ {
		jsn := []byte(fmt.Sprintf(`
			{
				"chefId": "%s",
				"name": "レシピ%d",
				"servings": 4,
				"ingredient": [
					{
						"name": "材料1",
						"supplement": "補足1"
					},
					{
						"name": "材料2",
						"supplement": "補足2"
					}
				],
				"method": [
					{
						"html": "<p>煮る</p>",
						"supplement": {
							"補足1キー": "補足1バリュー"
						}
					},
					{
						"html": "<p>焼く</p>",
						"supplement": [
							{
								"補足1キー": "補足1バリュー"
							},
							{
								"補足2キー": "補足2バリュー",
								"補足3キー": "補足3バリュー"
							}
						]
					}
				],
				"imageUrl": "03ca1240-4fde-4dd0-8639-e8c026d4512b.webp",
				"introduction": "レシピ紹介文%d",
				"link": [
					"https://www.yahoo.co.jp",
					"https://www.amazon.co.jp"
				],
				"accessLevel": 1
			}`, chefId, i, i))
		_, err := s.q.CreateRecipe(context.Background(), jsn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) CreateUserData(c *gin.Context) {
	numRecipe, err := strconv.Atoi(c.Query("num_recipe"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ユーザーを一人登録
	email, _ := uuid.NewUUID()
	var createUserParams db.CreateUserParams
	createUserParams.Email = email.String() + "@gmail.com"
	createUserParams.Name = "ユーザー名"
	createUserParams.AuthServer = "test"
	createUserParams.AuthUserinfo = []byte("{}")
	_, err = s.q.CreateUser(context.Background(), createUserParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var param db.UpdateUserParams
	param.Email = createUserParams.Email
	param.Data = []byte(fmt.Sprintf(`
		{
			"name": "ユーザー名",
			"imageUrl": "1e902068-d0cf-4fa5-a957-de5c49762179.webp",
			"profile": "PROFILE_1",
			"link": [
				{
					"label": "LABEL_1",
					"url": "http://www.yahoo.co.jp"
				},
				{
					"label": "LABEL_2",
					"url": "http://www.yahoo.co.jp"
				}
			]
		}`))
	response, err := s.q.UpdateUser(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	usrId := utils.UUID2Str(response.ID)
	fmt.Printf("usrId = %s\n", usrId)

	// 一般シェフのレシピをnumRecipe件登録
	for i := 1; i <= numRecipe; i++ {
		jsn := []byte(fmt.Sprintf(`
			{
				"usrId": "%s",
				"name": "レシピ%d",
				"servings": 4,
				"ingredient": [
					{
						"name": "材料1",
						"supplement": "補足1"
					},
					{
						"name": "材料2",
						"supplement": "補足2"
					}
				],
				"method": [
					{
						"html": "<p>煮る</p>",
						"supplement": {
							"補足1キー": "補足1バリュー"
						}
					},
					{
						"html": "<p>焼く</p>",
						"supplement": [
							{
								"補足1キー": "補足1バリュー"
							},
							{
								"補足2キー": "補足2バリュー",
								"補足3キー": "補足3バリュー"
							}
						]
					}
				],
				"imageUrl": "03ca1240-4fde-4dd0-8639-e8c026d4512b.webp",
				"introduction": "レシピ紹介文%d",
				"link": [
					"https://www.yahoo.co.jp",
					"https://www.amazon.co.jp"
				],
				"accessLevel": 1
			}`, usrId, i, i))
		_, err := s.q.CreateRecipe(context.Background(), jsn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, response)
}
