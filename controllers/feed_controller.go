package controllers

import (
	"sort"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/Pratham-Mishra04/interact/utils/select_fields"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetFeed(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var followings []models.FollowFollower
	if err := initializers.DB.Model(&models.FollowFollower{}).Where("follower_id = ?", loggedInUserID).Find(&followings).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	followingIDs := make([]uuid.UUID, len(followings))
	for i, following := range followings {
		followingIDs[i] = following.FollowedID
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	var posts []models.Post
	if err := paginatedDB.
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
		Joins("JOIN users ON posts.user_id = users.id AND users.active = ?", true).
		Where("is_flagged=?", false).
		Where("user_id = ? OR user_id IN (?)", loggedInUserID, followingIDs).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.IncrementPostImpression(posts)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"feed":   posts,
	})
}

type CombinedFeedItem interface {
	GetCreatedAt() time.Time
}

// Post implements CombinedFeedItem interface
type PostAlias models.Post

// GetCreatedAt is a method for PostAlias to satisfy the interface
func (p PostAlias) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// Announcement implements CombinedFeedItem interface
type AnnouncementAlias models.Announcement

// GetCreatedAt is a method for AnnouncementAlias to satisfy the interface
func (a AnnouncementAlias) GetCreatedAt() time.Time {
	return a.CreatedAt
}

// Poll implements CombinedFeedItem interface
type PollAlias models.Poll

// GetCreatedAt is a method for PollAlias to satisfy the interface
func (p PollAlias) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func GetCombinedFeed(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	var followings []models.FollowFollower
	if err := initializers.DB.Preload("Followed").Where("follower_id = ?", parsedUserID).Find(&followings).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	followingIDs := make([]uuid.UUID, len(followings))
	for i, following := range followings {
		followingIDs[i] = following.FollowedID
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	var posts []models.Post
	if err := paginatedDB.
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
		Where("is_flagged=?", false).
		Joins("JOIN users ON posts.user_id = users.id AND users.active = ?", true).
		Where("user_id = ? OR user_id IN (?)", parsedUserID, followingIDs).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var memberOrgIDs []uuid.UUID
	var followingOrgIDs []uuid.UUID

	var orgMemberships []models.OrganizationMembership
	if err := initializers.DB.Find(&orgMemberships, "user_id = ?", parsedUserID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for _, orgMembership := range orgMemberships {
		memberOrgIDs = append(memberOrgIDs, orgMembership.OrganizationID)
	}

	for _, followFollower := range followings {
		if followFollower.Followed.OrganizationStatus {
			var organization models.Organization
			if err := initializers.DB.Where("user_id = ?", followFollower.Followed.ID).Find(&organization).Error; err == nil {
				followingOrgIDs = append(followingOrgIDs, organization.ID)
			}
		}
	}

	paginatedDB = API.Paginator(c)(initializers.DB)

	var announcements []models.Announcement
	if err := paginatedDB.
		Preload("Organization").
		Preload("Organization.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("TaggedUsers", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.ShorterUser)
		}).
		Where("(organization_id IN (?) OR (is_open=true AND organization_id IN (?)))", memberOrgIDs, followingOrgIDs).
		Order("created_at DESC").
		Find(&announcements).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	paginatedDB = API.Paginator(c)(initializers.DB)

	db := paginatedDB.
		Preload("Organization").
		Preload("Organization.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User)
		}).
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("options.created_at DESC")
		}).
		Preload("Options.VotedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select(select_fields.User).Limit(3)
		}).
		Where("(organization_id IN (?) OR (is_open=true AND organization_id IN (?)))", memberOrgIDs, followingOrgIDs)

	var polls []models.Poll
	if err := db.Order("created_at DESC").Find(&polls).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	// Combine posts, announcements and polls, and sort them by created_at
	var combinedFeed []CombinedFeedItem
	for _, a := range announcements {
		combinedFeed = append(combinedFeed, AnnouncementAlias(a))
	}
	for _, p := range polls {
		combinedFeed = append(combinedFeed, PollAlias(p))
	}
	for _, p := range posts {
		combinedFeed = append(combinedFeed, PostAlias(p))
	}

	// Sort the combined news feed by created_at
	sort.Slice(combinedFeed, func(i, j int) bool {
		return combinedFeed[i].GetCreatedAt().After(combinedFeed[j].GetCreatedAt())
	})

	go routines.IncrementPostImpression(posts)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"feed":   combinedFeed,
	})
}
