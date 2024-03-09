package controllers

import (
	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/Pratham-Mishra04/interact/utils/select_fields"
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
	if err := initializers.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("RePost").
		Preload("RePost.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("RePost.TaggedUsers", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.ShorterUser)
		}).
		Preload("TaggedUsers", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.ShorterUser)
		}).
		First(&post, "id = ?", parsedPostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go cache.SetPost(postID, &post)

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
		Preload("RePost.User").
		Preload("RePost.TaggedUsers").
		Preload("User").
		Preload("TaggedUsers").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
	if err := paginatedDB.Preload("User").
		Preload("RePost").
		Preload("RePost.User").
		Preload("TaggedUsers").
		Preload("RePost.TaggedUsers").
		Where("user_id = ?", loggedInUserID).Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
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

	parsedLoggedInUserID, err := uuid.Parse(c.GetRespHeader("loggedInUserID"))
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}

	var user models.User
	if err := initializers.DB.Where("id=?", parsedLoggedInUserID).First(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	if !user.Verified {
		return &fiber.Error{Code: 401, Message: config.VERIFICATION_ERROR}
	}

	if err := helpers.Validate[schemas.PostCreateSchema](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	// images, err := utils.SaveMultipleFiles(c, "images", "post", true, 1280, 720)
	images, err := utils.UploadMultipleImages(c, "images", helpers.PostClient, 1280, 720)
	if err != nil {
		return err
	}

	newPost := models.Post{
		UserID:  parsedLoggedInUserID,
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
		newPost.IsRePost = true
	}

	var taggedUsers []models.User

	if reqBody.TaggedUsernames != nil {
		for _, username := range reqBody.TaggedUsernames {
			var user models.User
			if err := initializers.DB.First(&user, "username=?", username).Error; err == nil {
				taggedUsers = append(taggedUsers, user)
			}
		}

		newPost.TaggedUsers = taggedUsers
	}

	result := initializers.DB.Create(&newPost)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	if len(taggedUsers) > 0 {
		for _, user := range taggedUsers {
			go routines.SendTaggedNotification(user.ID, parsedLoggedInUserID, &newPost.ID, nil)
		}
	}

	//TODO6 convert to routine
	routines.GetBlurHashesForPost(c, "images", &newPost)

	if err := initializers.DB.Preload("User").
		Preload("RePost").
		Preload("RePost.User").
		Preload("TaggedUsers").
		First(&newPost).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	orgMemberID := c.GetRespHeader("orgMemberID")
	orgID := c.Params("orgID")
	if orgMemberID != "" && orgID != "" {
		parsedOrgID, err := uuid.Parse(orgID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
		}

		parsedOrgMemberID, err := uuid.Parse(orgMemberID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid User ID."}
		}
		go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 6, &newPost.ID, nil, nil, nil, nil, nil, nil, nil, nil, nil, "")
	}
	if reqBody.RePostID != "" {
		go routines.IncrementReposts(*newPost.RePostID)
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Post Added",
		"post":    newPost,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").Preload("TaggedUsers").First(&post, "id = ? and user_id=?", parsedPostID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var reqBody schemas.PostUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	if reqBody.Content != "" {
		post.Content = reqBody.Content
	}
	if reqBody.Tags != nil {
		post.Tags = *reqBody.Tags
	}

	post.Edited = true

	var newTaggedUsers []models.User
	var usersToRemove []models.User

	if reqBody.TaggedUsernames != nil { //TODO7 not working
		// Create a map to store existing tagged users for quick comparison
		existingTaggedUsers := make(map[uuid.UUID]models.User)
		for _, user := range post.TaggedUsers {
			existingTaggedUsers[user.ID] = user
		}

		for _, username := range reqBody.TaggedUsernames {
			var user models.User
			if err := initializers.DB.First(&user, "username=?", username).Error; err == nil {
				newTaggedUsers = append(newTaggedUsers, user)
			}
		}

		// Compare and find users to remove
		for _, existingUser := range post.TaggedUsers {
			if _, exists := existingTaggedUsers[existingUser.ID]; !exists {
				usersToRemove = append(usersToRemove, existingUser)
			}
		}

		// Update the TaggedUsers with new users
		post.TaggedUsers = newTaggedUsers
	}

	tx := initializers.DB.Begin()

	if err := tx.Save(&post).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for _, userToRemove := range usersToRemove {
		tx.Model(&post).Association("TaggedUsers").Delete(userToRemove)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if len(newTaggedUsers) > 0 {
		for _, user := range newTaggedUsers {
			go routines.SendTaggedNotification(user.ID, parsedLoggedInUserID, &post.ID, nil)
		}
	}

	orgMemberID := c.GetRespHeader("orgMemberID")
	orgID := c.Params("orgID")
	if orgMemberID != "" && orgID != "" {
		parsedOrgID, err := uuid.Parse(orgID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
		}

		parsedOrgMemberID, err := uuid.Parse(orgMemberID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid User ID."}
		}
		go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 8, &post.ID, nil, nil, nil, nil, nil, nil, nil, nil, nil, "")
	}

	go cache.RemovePost(postID)

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

	orgMemberID := c.GetRespHeader("orgMemberID")
	orgID := c.Params("orgID")

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, "id = ? AND user_id=?", parsedPostID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "DeletePost")
		}
	}()

	// Delete all the shared messages
	var messages []models.Message
	if err := tx.Find(&messages, "post_id=?", parsedPostID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	for _, message := range messages {
		if err := tx.Delete(&message).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	// Delete all the tagged users
	if err := tx.Model(&post).Association("TaggedUsers").Clear(); err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	// Edit all the posts where this is a repost
	var posts []models.Post
	if err := tx.Where("re_post_id=?", post.ID).Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for _, p := range posts {
		p.RePostID = nil
		if err := tx.Save(&p).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	// Delete the post
	if err := tx.Delete(&post).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for _, image := range post.Images {
		go routines.DeleteFromBucket(helpers.PostClient, image)
	}

	if orgMemberID != "" && orgID != "" {
		parsedOrgID, err := uuid.Parse(orgID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
		}

		parsedOrgMemberID, err := uuid.Parse(orgMemberID)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid User ID."}
		}
		go routines.MarkOrganizationHistory(parsedOrgID, parsedOrgMemberID, 7, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, post.Content)
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Post deleted successfully",
	})
}
