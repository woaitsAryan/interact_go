package controllers

import (
	"log"

	"github.com/Pratham-Mishra04/interact/config"
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

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, "id = ?", parsedPostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"post":    post,
	})
}

func GetUserPosts(c *fiber.Ctx) error {
	userID := c.Params("userID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var posts []models.Post
	if err := paginatedDB.Preload("User").Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
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

	var posts []models.Post
	if err := paginatedDB.Preload("User").Where("user_id = ?", loggedInUserID).Find(&posts).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   posts,
	})
}

func GetMyLikedPosts(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var postLikes []models.UserPostLike
	if err := initializers.DB.Where("user_id = ?", loggedInUserID).Find(&postLikes).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	var postIDs []string
	for _, post := range postLikes {
		postIDs = append(postIDs, post.PostID.String())
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   postIDs,
	})
}

func AddPost(c *fiber.Ctx) error {
	var reqBody schemas.PostCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.PostCreateSchema](reqBody); err != nil {
		return err
	}

	parsedID, err := uuid.Parse(c.GetRespHeader("loggedInUserID"))
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}

	images, err := utils.SaveMultipleFiles(c, "images", "post", true, 1280, 720)
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
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, "id = ?", newPost.ID).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Post Added",
		"post":    post,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, "id = ? and user_id=?", parsedPostID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	var updatePost schemas.PostUpdateSchema
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
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Post updated successfully",
		"post":    post,
	})
}

func DeletePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, "id = ? AND user_id=?", parsedPostID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	for _, image := range post.Images {
		err := utils.DeleteFile("post", image)
		if err != nil {
			log.Printf("Error while deleting post pic: %e", err)
		}
	}

	if err := initializers.DB.Delete(&post).Error; err != nil {
		return &fiber.Error{Code: 500, Message: config.DATABASE_ERROR}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Post deleted successfully",
	})
}
