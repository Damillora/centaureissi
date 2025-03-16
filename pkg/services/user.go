package services

import (
	"errors"

	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/Damillora/centaureissi/pkg/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (cs *CentaureissiService) CreateUser(model models.UserCreateModel) (*schema.User, error) {
	existsUser, err := cs.repository.ExistsUserByUsername(model.Username)
	if err != nil {
		return nil, err
	}
	if existsUser {
		return nil, errors.New("user already exists")
	}

	passwd, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &schema.User{
		ID:       uuid.NewString(),
		Username: model.Username,
		Password: string(passwd),
	}
	err = cs.repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Create INBOX mailbox to prevent headaches
	_, err = cs.CreateMailbox(user.ID, "INBOX")
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (cs *CentaureissiService) GetUserById(id string) (*schema.User, error) {
	user, err := cs.repository.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (cs *CentaureissiService) GetUserByUsername(username string) (*schema.User, error) {
	user, err := cs.repository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (cs *CentaureissiService) UpdateUserProfile(id string, model models.UserUpdateModel) (*schema.User, error) {
	user, err := cs.repository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	user.Username = model.Username

	err = cs.repository.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (cs *CentaureissiService) UpdateUserPassword(id string, model models.UserUpdatePasswordModel) (*schema.User, error) {
	user, err := cs.repository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user.Password != "" {
		verifyErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(model.OldPassword))
		if verifyErr != nil {
			return nil, verifyErr
		}

		passwd, err := bcrypt.GenerateFromPassword([]byte(model.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(passwd)
	}

	err = cs.repository.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
