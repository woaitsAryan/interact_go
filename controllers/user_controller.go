package controllers

import (
	"errors"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetViews(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	viewsArr, count, err := utils.GetProfileViews(parsedUserID)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"viewsArr": viewsArr,
		"count":    count,
	})
}

func GetMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var user models.User
	initializers.DB.
		Preload("Profile").
		Preload("Profile.Achievements").
		First(&user, "id = ?", userID)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"user":    user,
	})
}

func GetUser(c *fiber.Ctx) error {
	username := c.Params("username")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var user models.User
	initializers.DB.Preload("Profile").First(&user, "username = ?", username)

	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No user of this username found."}
	}

	if user.ID.String() != loggedInUserID {
		routines.UpdateProfileViews(&user)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"user":    user,
	})
}

func GetMyLikes(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var likes []models.Like
	if err := initializers.DB.
		Find(&likes, "user_id = ?", loggedInUserID).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var likeIDs []string
	for _, like := range likes {
		if like.PostID != nil {
			likeIDs = append(likeIDs, like.PostID.String())
		} else if like.ProjectID != nil {
			likeIDs = append(likeIDs, like.ProjectID.String())
		} else if like.CommentID != nil {
			likeIDs = append(likeIDs, like.CommentID.String())
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"likes":   likeIDs,
	})
}

func GetMyOrgMemberships(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	populate := c.Query("populate", "false")

	var memberships []models.OrganizationMembership

	if populate == "true" {
		if err := initializers.DB.
			Preload("Organization").
			Preload("Organization.User").
			Find(&memberships, "user_id = ?", loggedInUserID).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	} else {
		if err := initializers.DB.
			Find(&memberships, "user_id = ?", loggedInUserID).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":      "success",
		"message":     "User Found",
		"memberships": memberships,
	})
}

func UpdateMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var reqBody schemas.UserUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	// if err := helpers.Validate[schemas.UserUpdateSchema](reqBody); err != nil {
	// 	return &fiber.Error{Code: 400, Message: err.Error()}
	// }

	oldProfilePic := user.ProfilePic
	oldCoverPic := user.CoverPic

	// picName, err := utils.SaveFile(c, "profilePic", "user/profilePics", true, 500, 500)
	picName, err := utils.UploadFile(c, "profilePic", helpers.UserProfileClient, 500, 500)
	if err != nil {
		return err
	}
	reqBody.ProfilePic = &picName

	// coverName, err := utils.SaveFile(c, "coverPic", "user/coverPics", true, 900, 400)
	coverName, err := utils.UploadFile(c, "coverPic", helpers.UserCoverClient, 900, 400)
	if err != nil {
		return err
	}
	reqBody.CoverPic = &coverName

	// updateUserValue := reflect.ValueOf(&reqBody).Elem()
	// userValue := reflect.ValueOf(&user).Elem()

	// for i := 0; i < updateUserValue.NumField(); i++ {
	// 	field := updateUserValue.Type().Field(i)
	// 	fieldName := field.Name

	// 	if fieldValue := updateUserValue.Field(i); fieldValue.IsValid() && !fieldValue.IsZero() {
	// 		userField := userValue.FieldByName(fieldName)
	// 		if userField.IsValid() && userField.CanSet() {
	// 			userField.Set(fieldValue)
	// 		}
	// 	}
	// }

	if reqBody.Name != nil {
		user.Name = *reqBody.Name
	}
	if reqBody.Bio != nil {
		user.Bio = *reqBody.Bio
	}
	if reqBody.Tagline != nil {
		user.Tagline = *reqBody.Tagline
	}
	if reqBody.ProfilePic != nil && *reqBody.ProfilePic != "" {
		user.ProfilePic = *reqBody.ProfilePic
	}
	if reqBody.CoverPic != nil && *reqBody.CoverPic != "" {
		user.CoverPic = *reqBody.CoverPic
	}
	if reqBody.Tags != nil {
		user.Tags = *reqBody.Tags
	}
	if reqBody.Links != nil {
		user.Links = *reqBody.Links
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if *reqBody.ProfilePic != "" {
		go routines.DeleteFromBucket(helpers.UserProfileClient, oldProfilePic)
	}

	if *reqBody.CoverPic != "" {
		go routines.DeleteFromBucket(helpers.UserCoverClient, oldCoverPic)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
		"user":    user,
	})
}

func DeactivateMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	//TODO send email for verification

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	user.Active = false

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "User deactivated successfully",
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	var reqBody struct {
		Password        string `json:"password"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	if reqBody.NewPassword != reqBody.ConfirmPassword {
		return &fiber.Error{Code: 400, Message: "Passwords do not match."}
	}

	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var user models.User
	initializers.DB.First(&user, "id = ?", loggedInUserID)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return &fiber.Error{Code: 400, Message: "Incorrect Password."}
	}

	//TODO send email for verification

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), 10)

	if err != nil {
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	user.Password = string(hash)
	user.PasswordChangedAt = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return CreateSendToken(c, user, 200, "Password updated successfully")
}

func UpdateEmail(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var reqBody struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	var emailCheckUser models.User
	if err := initializers.DB.First(&emailCheckUser, "email = ?", reqBody.Email).Error; err == nil {
		return &fiber.Error{Code: 400, Message: "Email Address Already In Use."}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	user.Email = reqBody.Email
	user.Verified = false

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
	})
}

func UpdatePhoneNo(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var reqBody struct {
		PhoneNo string `json:"phoneNo"  validate:"e164"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	var phoneNoCheckUser models.User
	if err := initializers.DB.First(&phoneNoCheckUser, "phone_no = ?", reqBody.PhoneNo).Error; err == nil {
		return &fiber.Error{Code: 400, Message: "Phone Number Already In Use."}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	user.PhoneNo = reqBody.PhoneNo

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
	})
}

func UpdateResume(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	resumePath, err := utils.UploadResume(c)
	if err != nil {
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	oldResume := user.Resume

	user.Resume = resumePath

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	if oldResume != "" {
		go routines.DeleteFromBucket(helpers.UserResumeBucket, oldResume)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
		"resume":  resumePath,
	})
}

func Deactive(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	user.Active = false
	user.DeactivatedAt = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Account Deactived",
	})
}
