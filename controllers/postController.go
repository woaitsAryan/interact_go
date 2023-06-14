package controllers

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetPost(c *fiber.Ctx) error {

	postID := c.Params("postID")

	postUserSelectedDB := utils.PostSelectConfig(initializers.DB.Preload("User"))

	var post models.Post
	if err := postUserSelectedDB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	// user, err := helpers.Filter(post.User, []string{"username", "name", "profilePic"})
	// if err != nil {
	// 	return err
	// }

	// filteredUser, ok := user.(models.User)
	// if !ok {
	// 	return &fiber.Error{Code: 500, Message: "Failed to assert user type"}
	// }

	// post.User = filteredUser

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"post":    post,
	})
}

func GetUserPosts(c *fiber.Ctx) error {
	userID := c.Params("userID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	postUserSelectedDB := utils.PostSelectConfig(paginatedDB.Preload("User"))

	var posts []models.Post
	if err := postUserSelectedDB.Where("user_id = ?", userID, false).Find(&posts).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   posts,
	})
}

func GetMyPosts(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	postUserSelectedDB := utils.PostSelectConfig(paginatedDB.Preload("User"))

	var posts []models.Post
	if err := postUserSelectedDB.Where("user_id = ?", loggedInUserID).Find(&posts).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   posts,
	})
}

func AddPost(c *fiber.Ctx) error {
	var reqBody schemas.PostCreateScheam
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.PostCreateScheam](reqBody); err != nil {
		return err
	}

	parsedID, err := uuid.Parse(c.GetRespHeader("loggedInUserID"))
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}

	images, err := utils.SaveMultipleFiles(c, "images", "posts/images", true, 1280, 720)
	if err != nil {
		return err
	}

	newPost := models.Post{
		UserID:  parsedID,
		Content: reqBody.Content,
		Images:  images,
		Tags:    reqBody.Tags,
	}

	result := initializers.DB.Create(&newPost)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating post"}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Post Added",
		"post":    newPost,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	var post models.Post
	initializers.DB.First(&post, "id = ?", postID)

	var updatePost schemas.PostUpdateScheam
	if err := c.BodyParser(&updatePost); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	if updatePost.Content != "" {
		post.Content = updatePost.Content
	}
	if len(updatePost.Tags) != 0 {
		post.Tags = updatePost.Tags
	}

	post.Edited = true

	if err := initializers.DB.Save(&post).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Post updated successfully",
		"post":    post,
	})
}

func DeletePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	var post models.Post
	initializers.DB.First(&post, "id = ?", postID)

	if err := initializers.DB.Delete(&post).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Post deleted successfully",
	})
}
