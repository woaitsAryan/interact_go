package controllers

import (
	"github.com/Pratham-Mishra04/interact/cache"
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

	postInCache, err := cache.GetPost(postID)

	if err == nil {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "",
			"post":    postInCache,
		})
	}

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("RePost").Preload("User").First(&post, "id = ?", parsedPostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	cache.SetPost(postID, &post)

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
	if err := paginatedDB.
		Preload("RePost").
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
	if err := paginatedDB.Preload("RePost").Preload("User").Where("user_id = ?", loggedInUserID).Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   posts,
	})
}

func GetMyLikedPosts(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var postLikes []models.Like
	if err := initializers.DB.Where("user_id = ? AND post_id IS NOT NULL", loggedInUserID).Find(&postLikes).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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

	parsedID, err := uuid.Parse(c.GetRespHeader("loggedInUserID"))
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}

	var user models.User
	if err := initializers.DB.Where("id=?", parsedID).First(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}
	if !user.Verified {
		return &fiber.Error{Code: 401, Message: config.VERIFICATION_ERROR}
	}

	if err := helpers.Validate[schemas.PostCreateSchema](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
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

	if reqBody.RePostID != "" {
		parsedRePostID, err := uuid.Parse(reqBody.RePostID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Post ID in rePost"}
		}
		newPost.RePostID = &parsedRePostID
	}

	result := initializers.DB.Create(&newPost)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
	}

	if reqBody.TaggedUserIDS != nil {
		for _, userID := range reqBody.TaggedUserIDS {
			parsedUserID, err := uuid.Parse(userID)
			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid User ID in tagged users"}
			}
			userTag := models.UserPostTag{
				UserID: parsedUserID,
				PostID: newPost.ID,
			}
			result := initializers.DB.Create(&userTag)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").Preload("RePost").Preload("RePost.User").First(&post, "id = ?", newPost.ID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var reqBody schemas.PostUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	if reqBody.Content != "" {
		post.Content = reqBody.Content
	}
	if len(reqBody.Tags) != 0 {
		post.Tags = reqBody.Tags
	}

	post.Edited = true

	if err := initializers.DB.Save(&post).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if reqBody.TaggedUserIDS != nil {
		for _, userID := range reqBody.TaggedUserIDS {
			parsedUserID, err := uuid.Parse(userID)
			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid User ID in tagged users"}
			}
			userTag := models.UserPostTag{
				UserID: parsedUserID,
				PostID: post.ID,
			}
			result := initializers.DB.Create(&userTag)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: result.Error}
			}
		}
	}

	cache.RemovePost(postID)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Post updated successfully",
		"post":    post,
	})
}

func DeletePost(c *fiber.Ctx) error {
	//TODO Handle what happens when the post to be deleted is a post of a repost
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	for _, image := range post.Images {
		err := utils.DeleteFile("post", image)
		if err != nil {
			initializers.Logger.Warnf("Error while deleting post pic", err)
		}
	}

	var messages []models.Message
	if err := initializers.DB.Find(&messages, "post_id=?", parsedPostID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}

	for _, message := range messages {
		if err := initializers.DB.Delete(&message).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}

	if err := initializers.DB.Delete(&post).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Post deleted successfully",
	})
}
