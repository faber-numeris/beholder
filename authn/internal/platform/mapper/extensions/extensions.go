package extensions

import (
	"time"

	"github.com/faber-numeris/beholder/authn/internal/adapters/outbound/postgres/gen"
	"github.com/faber-numeris/beholder/authn/internal/core/domain"
	foundation "github.com/faber-numeris/foundation/beholder/api"
)

func StringToULID(s string) foundation.ULID {
	return s
}

func StringToOptString(s string) *string {
	return &s
}

func OptStringToString(os *string) (string, error) {
	if os != nil {
		return *os, nil
	}
	return "", nil
}

func StringToOptstring(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func TimeToTimePtr(t time.Time) (*time.Time, error) {
	return &t, nil
}

func StringToByteSlice(s string) ([]byte, error) {
	return []byte(s), nil
}

func UserTypeToApiUserType(t domain.UserType) foundation.UserType {
	switch t {
	case domain.UserTypeUser:
		return foundation.UserTypeUSER
	case domain.UserTypeServiceAccount:
		return foundation.UserTypeSERVICEACCOUNT
	case domain.UserTypeDevice:
		return foundation.UserTypeDEVICE
	default:
		return foundation.UserType(t)
	}
}

func SQLCUserToUserProfile(user gen.User) *domain.UserProfile {
	if user.FirstName == "" && user.LastName == "" && user.Locale == "" && user.Timezone == "" {
		return nil
	}
	return &domain.UserProfile{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Locale:    user.Locale,
		Timezone:  user.Timezone,
	}
}

func UserProfileToSQLCUser(profile *domain.UserProfile) (firstName, lastName, locale, timezone string) {
	if profile == nil {
		return
	}
	return profile.FirstName, profile.LastName, profile.Locale, profile.Timezone
}
