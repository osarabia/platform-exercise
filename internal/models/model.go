package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const TableNameUser = "users"

type User struct {
	gorm.Model        // includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"size:100;not null"`
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Birthday   *time.Time
}

func (u *User) validatePassword() error {
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		uppercase = regexp.MustCompile(`[A-Z]`)
		lowercase = regexp.MustCompile(`[a-z]`)
		number    = regexp.MustCompile(`[0-9]`)
		special   = regexp.MustCompile(`[!@#~$%^&*()+|_.,<>?/\\[\]{}-]`)
	)

	switch {
	case !uppercase.MatchString(u.Password):
		return errors.New("password must include at least one uppercase letter")
	case !lowercase.MatchString(u.Password):
		return errors.New("password must include at least one lowercase letter")
	case !number.MatchString(u.Password):
		return errors.New("password must include at least one digit")
	case !special.MatchString(u.Password):
		return errors.New("password must include at least one special character")
	}

	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	err = u.validatePassword()
	if err != nil {
		return fmt.Errorf("%s: %s", "invalid password:", err)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) TableName() string {
	return TableNameUser
}
