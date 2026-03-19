package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Converts (*)Account or [](*)Account to its json equivalent marshalled bytes
func AccountDTO(a any) ([]byte, error) {
	account, ok := a.(Account)
	if ok {
		return json.Marshal(account)
	}

	accountPtr, ok := a.(*Account)
	if ok {
		return json.Marshal(*accountPtr)
	}

	accountArr, ok := a.([]Account)
	if ok {
		return json.Marshal(accountArr)
	}

	accountPtrArr, ok := a.([]*Account)
	if ok {
		var accounts []Account
		for _, v := range accountPtrArr {
			accounts = append(accounts, *v)
		}
		return json.Marshal(accounts)
	}

	return nil, errors.New("failed to marshal: not an Account object")
}
