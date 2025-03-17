package services

import (
	"errors"

	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/google/uuid"
)

func (cs *CentaureissiService) CounterUidValidity(id string) (uint32, error) {
	return cs.repository.CounterUidValidity(id)
}
func (cs *CentaureissiService) CounterMailboxUid(id string) (uint32, error) {
	return cs.repository.CounterMailboxUid(id)
}
func (cs *CentaureissiService) CounterMailboxUidNext(id string) (uint32, error) {
	num, err := cs.repository.CounterMailboxUid(id)
	if err != nil {
		return 0, err
	}

	return num + 1, nil
}

func (cs *CentaureissiService) IncrementUidValidity(id string) (uint32, error) {
	return cs.repository.IncrementUidValidity(id)
}
func (cs *CentaureissiService) IncrementMailboxUid(id string) (uint32, error) {
	return cs.repository.IncrementMailboxUid(id)
}

func (cs *CentaureissiService) ListMailboxesByUserId(id string) ([]*schema.Mailbox, error) {
	mboxes, err := cs.repository.ListMailboxesByUserId(id)
	if err != nil {
		return nil, err
	}
	return mboxes, err
}

func (cs *CentaureissiService) GetMailboxByUserIdAndName(userId string, name string) (*schema.Mailbox, error) {
	mbox, err := cs.repository.GetMailboxByUserIdAndName(userId, name)
	if err != nil {
		return nil, err
	}
	if mbox == nil {
		return nil, nil
	}
	return mbox, nil
}

func (cs *CentaureissiService) CreateMailbox(userId string, name string) (*schema.Mailbox, error) {

	existsUser, err := cs.repository.ExistsMailboxByUserIdAndName(userId, name)
	if err != nil {
		return nil, err
	}
	if existsUser {
		return nil, errors.New("mailbox already exists")
	}

	uidValidity, err := cs.IncrementUidValidity(userId)
	if err != nil {
		return nil, err
	}

	mbox := &schema.Mailbox{
		Id:          uuid.NewString(),
		UserId:      userId,
		Name:        name,
		UidValidity: uidValidity,
		Subscribed:  false,
	}
	err = cs.repository.CreateMailbox(mbox)
	if err != nil {
		return nil, err
	}

	return mbox, err
}

func (cs *CentaureissiService) UpdateMailboxSubscribeStatus(userId string, name string, subscribed bool) (*schema.Mailbox, error) {
	mbox, err := cs.repository.GetMailboxByUserIdAndName(userId, name)
	if err != nil {
		return nil, err
	}
	if mbox == nil {
		// Ignore unsub status for deleted mailbox
		return nil, nil
	}

	mbox.Subscribed = subscribed
	err = cs.repository.UpdateMailbox(mbox)
	if err != nil {
		return nil, err
	}

	return mbox, nil
}

func (cs *CentaureissiService) UpdateMailboxName(userId string, name string, newName string) (*schema.Mailbox, error) {
	mbox, err := cs.repository.GetMailboxByUserIdAndName(userId, name)
	if err != nil {
		return nil, err
	}
	if mbox == nil {
		// Ignore unsub status for deleted mailbox
		return nil, nil
	}
	mbox.Name = newName

	err = cs.repository.UpdateMailboxOldNameIndex(mbox.Id, name)
	if err != nil {
		return nil, err
	}
	err = cs.repository.UpdateMailbox(mbox)
	if err != nil {
		return nil, err
	}

	return mbox, nil
}

func (cs *CentaureissiService) DeleteMailbox(userId string, name string) error {
	existingMbox, err := cs.repository.GetMailboxByUserIdAndName(userId, name)
	if err != nil {
		return err
	}
	if existingMbox == nil {
		return errors.New("mailbox does not exists")
	}
	err = cs.repository.DeleteMailbox(existingMbox.Id)
	if err != nil {
		return err
	}

	return nil
}
