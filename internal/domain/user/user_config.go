package user

import shared "go-cs/internal/pkg/domain"
import domain_message "go-cs/internal/domain/pkg/message"

type UserConfig struct {
	shared.AggregateRoot

	Id        int64
	UserId    int64
	Key       string
	Value     string
	CreatedAt int64
	UpdatedAt int64
}

// setValue
func (u *UserConfig) SetValue(value string, oepr shared.Oper) {
	oldValue := u.Value

	u.Value = value

	u.AddMessage(oepr, &domain_message.PersonalSetUserConfig{
		UserId:   u.UserId,
		Key:      u.Key,
		NewValue: u.Value,
		OldValue: oldValue,
	})
}
