package authorization

import (
	"time"

	"EverythingSuckz/fsb/config"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type AuthorizedUser struct {
	UserID     int64 `gorm:"primaryKey"`
	ApprovedAt time.Time
}

var db *gorm.DB

// Init initializes the authorization database.
func Init() error {
	var err error

	db, err = gorm.Open(sqlite.Open("authorization.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&AuthorizedUser{}); err != nil {
		return err
	}

	// Automatically authorize the administrator.
	if config.ValueOf.AdminID != 0 {
		if err := Approve(config.ValueOf.AdminID); err != nil {
			return err
		}
	}

	return nil
}

// IsAuthorized checks whether a user has been approved.
func IsAuthorized(userID int64) bool {
	if db == nil {
		return false
	}

	var user AuthorizedUser
	result := db.First(&user, "user_id = ?", userID)

	return result.Error == nil
}

// Approve adds a user to the authorized list.
func Approve(userID int64) error {
	if db == nil {
		return gorm.ErrInvalidDB
	}

	user := AuthorizedUser{
		UserID:     userID,
		ApprovedAt: time.Now(),
	}

	return db.Save(&user).Error
}

// Revoke removes a user from the authorized list.
func Revoke(userID int64) error {
	if db == nil {
		return gorm.ErrInvalidDB
	}

	return db.Delete(&AuthorizedUser{}, "user_id = ?", userID).Error
}
