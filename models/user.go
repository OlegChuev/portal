package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User is a generated model from buffalo-auth, it serves as the base for username/password authentication.
type User struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`

	Password             string `json:"-" db:"-"`
	PasswordConfirmation string `json:"-" db:"-"`
}

// Create wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (user *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	user.Email = strings.ToLower(user.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}

	user.PasswordHash = string(ph)

	return tx.ValidateAndCreate(user)
}

// String is not required by pop and may be deleted
func (user User) String() string {
	ju, _ := json.Marshal(user)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (user *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error

	return validate.Validate(
		// check if string present
		&validators.StringIsPresent{
			Field: user.Email,
			Name:  "Email",
		},
		// check if string present
		&validators.StringIsPresent{
			Field: user.PasswordHash,
			Name:  "PasswordHash",
		},
		// check to see if the email address is already taken:
		&validators.FuncValidator{
			Field:   user.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", user.Email)
				q = q.Where("id != ?", user.ID)

				b, err = q.Exists(user)

				if err != nil {
					return false
				}

				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (user *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error

	return validate.Validate(
		&validators.StringIsPresent{
			Field: user.Password,
			Name:  "Password",
		},
		&validators.StringsMatch{
			Name:    "Password",
			Field:   user.Password,
			Field2:  user.PasswordConfirmation,
			Message: "Password does not match confirmation",
		},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
