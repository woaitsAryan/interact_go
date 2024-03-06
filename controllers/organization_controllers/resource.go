package organization_controllers

import (
	"fmt"
	"path"

	"github.com/Pratham-Mishra04/interact/cache"
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetOrgResourceBuckets(c *fiber.Ctx) error {
	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	var organization models.Organization
	if err := initializers.DB.Preload("User").Where("id = ?", parsedOrgID).First(&organization).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No Organization of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var resourceBuckets []models.ResourceBucket
	if err := initializers.DB.Where("organization_id = ?", parsedOrgID).Find(&resourceBuckets).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":          "success",
		"organization":    organization,
		"resourceBuckets": resourceBuckets,
	})
}

func GetResourceBucketFiles(c *fiber.Ctx) error {
	resourceBucketID := c.Params("resourceBucketID")
	parsedResourceBucketID, err := uuid.Parse(resourceBucketID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Resource Bucket ID."}
	}

	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	resourceBucketInCache, err := cache.GetResourceBucket(resourceBucketID)
	if err == nil {
		return c.Status(200).JSON(fiber.Map{
			"status":        "success",
			"resourceFiles": resourceBucketInCache.ResourceFiles,
		})
	}

	var resourceBucket models.ResourceBucket
	if err := initializers.DB.Preload("ResourceFiles").Preload("ResourceFiles.User").Where("id=? AND organization_id = ?", parsedResourceBucketID, parsedOrgID).First(&resourceBucket).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource Bucket does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go cache.SetResourceBucket(resourceBucket.ID.String(), &resourceBucket)

	return c.Status(200).JSON(fiber.Map{
		"status":        "success",
		"message":       "Resource Bucket added",
		"resourceFiles": resourceBucket.ResourceFiles,
	})
}

func ServeResourceFile(c *fiber.Ctx) error {
	parsedResourceFileID, err := uuid.Parse(c.Params("resourceFileID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Resource File ID."}
	}

	var resourceFile models.ResourceFile
	if err := initializers.DB.Where("id=?", parsedResourceFileID).First(&resourceFile).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource File does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	buffer, err := helpers.ResourceClient.GetBucketFile(resourceFile.Path)
	if err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.SERVER_ERROR, LogMessage: err.Error(), Err: err}
	}

	c.Set("Content-Type", "application/octet-stream")

	// Set the Content-Disposition header to prompt the user to download the file
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", resourceFile.Path))

	// Send the file content to the frontend
	return c.Send(buffer.Bytes())
}

func AddResourceBucket(c *fiber.Ctx) error {
	var reqBody schemas.ResourceBucketCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	resourceBucket := models.ResourceBucket{
		OrganizationID: parsedOrgID,
		Title:          reqBody.Title,
		Description:    reqBody.Description,
		ViewAccess:     reqBody.ViewAccess,
		EditAccess:     reqBody.EditAccess,
	}

	if err := initializers.DB.Create(&resourceBucket).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":         "success",
		"message":        "Resource Bucket added",
		"resourceBucket": resourceBucket,
	})
}

func AddResourceFile(c *fiber.Ctx) error {
	var reqBody schemas.ResourceFileCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedResourceBucketID, err := uuid.Parse(c.Params("resourceBucketID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Resource Bucket ID."}
	}

	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid User ID."}
	}

	var resourceBucket models.ResourceBucket
	if err := initializers.DB.Where("id=? AND organization_id = ?", parsedResourceBucketID, parsedOrgID).First(&resourceBucket).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource Bucket does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var fileExtension string
	var link string

	if reqBody.Link == "" {
		link, err = utils.UploadResourceFile(c)
		if err != nil {
			if err.Error() == "size-exceeded" {
				return &fiber.Error{Code: 400, Message: "File too large"}
			}
			return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, LogMessage: err.Error(), Err: err}
		}

		fileExtension = path.Ext(link)

		if len(fileExtension) > 0 {
			// Remove the leading dot from the extension
			fileExtension = fileExtension[1:]
		}
	} else {
		link = reqBody.Link
		fileExtension = ""
	}

	resourceFile := models.ResourceFile{
		ResourceBucketID: resourceBucket.ID,
		UserID:           parsedUserID,
		Title:            reqBody.Title,
		Description:      reqBody.Description,
		Path:             link,
		Type:             fileExtension,
	}

	if reqBody.Link == "" {
		resourceFile.FileUploaded = true
	}

	if err := initializers.DB.Create(&resourceFile).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.IncrementResourceBucketFiles(resourceBucket.ID)
	go cache.RemoveResourceBucket(resourceBucket.ID.String())

	initializers.DB.Preload("User").First(&resourceFile)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":       "success",
		"message":      "Resource File added",
		"resourceFile": resourceFile,
	})
}

func EditResourceBucket(c *fiber.Ctx) error {
	var reqBody schemas.ResourceBucketEditSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedResourceBucketID, err := uuid.Parse(c.Params("resourceBucketID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Resource Bucket ID."}
	}

	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Organization ID."}
	}

	var resourceBucket models.ResourceBucket
	if err := initializers.DB.Where("id=? AND organization_id = ?", parsedResourceBucketID, parsedOrgID).First(&resourceBucket).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource Bucket does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if reqBody.Title != "" {
		resourceBucket.Title = reqBody.Title
	}
	if reqBody.Description != nil {
		resourceBucket.Description = *reqBody.Description
	}
	if reqBody.ViewAccess != "" {
		resourceBucket.ViewAccess = reqBody.ViewAccess
	}
	if reqBody.EditAccess != "" {
		resourceBucket.EditAccess = reqBody.EditAccess
	}

	if err := initializers.DB.Save(&resourceBucket).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go cache.RemoveResourceBucket(resourceBucket.ID.String())

	return c.Status(200).JSON(fiber.Map{
		"status":         "success",
		"message":        "Resource Bucket Edited",
		"resourceBucket": resourceBucket,
	})
}

func EditResourceFile(c *fiber.Ctx) error {
	var reqBody schemas.ResourceFileCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedResourceFileID, err := uuid.Parse(c.Params("resourceFileID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Resource File ID."}
	}

	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Member ID."}
	}

	var resourceFile models.ResourceFile
	if err := initializers.DB.Where("id=? AND user_id=?", parsedResourceFileID, parsedUserID).First(&resourceFile).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource File does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if reqBody.Title != "" {
		resourceFile.Title = reqBody.Title
	}
	if reqBody.Description != "" {
		resourceFile.Description = reqBody.Description
	}

	if err := initializers.DB.Save(&resourceFile).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go cache.RemoveResourceBucket(resourceFile.ResourceBucketID.String())

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "Resource File added",
		"resourceFile": resourceFile,
	})
}

func DeleteResourceBucket(c *fiber.Ctx) error {
	//TODO14 add OTP here
	parsedResourceBucketID, err := uuid.Parse(c.Params("resourceBucketID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Resource Bucket ID."}
	}

	parsedOrgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Organization ID."}
	}

	var resourceBucket models.ResourceBucket
	if err := initializers.DB.Preload("ResourceFiles").Where("id=? AND organization_id = ?", parsedResourceBucketID, parsedOrgID).First(&resourceBucket).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource Bucket does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	files := resourceBucket.ResourceFiles

	if err := initializers.DB.Delete(&resourceBucket).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for _, file := range files {
		go routines.DeleteFromBucket(helpers.ResourceClient, file.Path)
	}

	go cache.RemoveResourceBucket(resourceBucket.ID.String())

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Resource Bucket deleted",
	})
}

func DeleteResourceFile(c *fiber.Ctx) error {
	parsedResourceFileID, err := uuid.Parse(c.Params("resourceFileID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Resource Bucket ID."}
	}

	parsedUserID, err := uuid.Parse(c.GetRespHeader("orgMemberID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Member ID."}
	}

	var resourceFile models.ResourceFile
	if err := initializers.DB.Where("id=? AND user_id=?", parsedResourceFileID, parsedUserID).First(&resourceFile).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource File does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	path := resourceFile.Path

	var resourceBucket models.ResourceBucket
	if err := initializers.DB.Where("id=?", resourceFile.ResourceBucketID).First(&resourceBucket).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource Bucket does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&resourceFile).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.DecrementResourceBucketFiles(resourceBucket.ID)
	go routines.DeleteFromBucket(helpers.ResourceClient, path)
	go cache.RemoveResourceBucket(resourceFile.ResourceBucketID.String())

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Resource File deleted",
	})
}
