package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
)

type User struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(UninitializedID, UninitializedVersion, fullName, phoneNumber)

}

func (u *User) Validate() error {
	fullNameLen := len([]rune(u.FullName))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf("invalid `Full Name` len: %d: %w", fullNameLen, core_errors.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("invalid `Phone Number` len: %d: %w", phoneNumberLen, core_errors.ErrInvalidArgument)
		}
		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid `Phone Number` format: %w", core_errors.ErrInvalidArgument)
		}
	}
	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func NewUserPatch(fullName Nullable[string], phoneNumber Nullable[string]) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("'FullName' can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

func (u *User) ApplyPatch(p UserPatch) error {
	//dвалидируем патч
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}
	//создаём временного пользователя
	tmpUser := *u
	//применяем патч к временному пользвателю
	if p.FullName.Set {
		tmpUser.FullName = *p.FullName.Value
	}
	if p.PhoneNumber.Set {
		tmpUser.PhoneNumber = p.PhoneNumber.Value
	}
	// валидируем пропатченного временного пользователя
	if err := tmpUser.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}
	//пропатченный пользователь прошёл валидацию перезаписываем значение временного пользователя в постоянного
	*u = tmpUser
	return nil
}
