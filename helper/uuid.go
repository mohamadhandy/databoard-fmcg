package helper

import "github.com/google/uuid"

type UuidHelper interface {
	GenerateUUID() string
}

type uuidHelper struct {
	u uuid.UUID
}

func InitUuidHelper() UuidHelper {
	return &uuidHelper{}
}

func (u *uuidHelper) GenerateUUID() string {
	u.u = uuid.New()
	return u.u.String()
}
