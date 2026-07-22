package postgresadapter

import (
	"context"

	"github.com/faber-numeris/beholder/authn/internal/adapters/outbound/postgres/gen"
)

// fakeQuerier is a hand-written test double for gen.Querier. Each field defaults to nil;
// tests set only the functions exercised by the scenario under test.
type fakeQuerier struct {
	confirmUserRegistrationFn        func(ctx context.Context, userid string) error
	createUserFn                     func(ctx context.Context, arg gen.CreateUserParams) (gen.User, error)
	createUserConfirmationFn         func(ctx context.Context, arg gen.CreateUserConfirmationParams) (gen.UserConfirmation, error)
	deleteExpiredUserConfirmationsFn func(ctx context.Context) error
	deleteUserFn                     func(ctx context.Context, id string) error
	deleteUserConfirmationFn         func(ctx context.Context, userid string) error
	getUserFn                        func(ctx context.Context, id string) (gen.User, error)
	getUserByEmailFn                 func(ctx context.Context, email string) (gen.User, error)
	getUserConfirmationByTokenFn     func(ctx context.Context, token string) (gen.UserConfirmation, error)
	getUserConfirmationByUserIDFn    func(ctx context.Context, userid string) (gen.UserConfirmation, error)
	listUsersFn                      func(ctx context.Context, arg gen.ListUsersParams) ([]gen.User, error)
	updatePasswordFn                 func(ctx context.Context, arg gen.UpdatePasswordParams) error
	updateUserFn                     func(ctx context.Context, arg gen.UpdateUserParams) (gen.User, error)
}

func (f *fakeQuerier) ConfirmUserRegistration(ctx context.Context, userid string) error {
	return f.confirmUserRegistrationFn(ctx, userid)
}

func (f *fakeQuerier) CreateUser(ctx context.Context, arg gen.CreateUserParams) (gen.User, error) {
	return f.createUserFn(ctx, arg)
}

func (f *fakeQuerier) CreateUserConfirmation(ctx context.Context, arg gen.CreateUserConfirmationParams) (gen.UserConfirmation, error) {
	return f.createUserConfirmationFn(ctx, arg)
}

func (f *fakeQuerier) DeleteExpiredUserConfirmations(ctx context.Context) error {
	return f.deleteExpiredUserConfirmationsFn(ctx)
}

func (f *fakeQuerier) DeleteUser(ctx context.Context, id string) error {
	return f.deleteUserFn(ctx, id)
}

func (f *fakeQuerier) DeleteUserConfirmation(ctx context.Context, userid string) error {
	return f.deleteUserConfirmationFn(ctx, userid)
}

func (f *fakeQuerier) GetUser(ctx context.Context, id string) (gen.User, error) {
	return f.getUserFn(ctx, id)
}

func (f *fakeQuerier) GetUserByEmail(ctx context.Context, email string) (gen.User, error) {
	return f.getUserByEmailFn(ctx, email)
}

func (f *fakeQuerier) GetUserConfirmationByToken(ctx context.Context, token string) (gen.UserConfirmation, error) {
	return f.getUserConfirmationByTokenFn(ctx, token)
}

func (f *fakeQuerier) GetUserConfirmationByUserID(ctx context.Context, userid string) (gen.UserConfirmation, error) {
	return f.getUserConfirmationByUserIDFn(ctx, userid)
}

func (f *fakeQuerier) ListUsers(ctx context.Context, arg gen.ListUsersParams) ([]gen.User, error) {
	return f.listUsersFn(ctx, arg)
}

func (f *fakeQuerier) UpdatePassword(ctx context.Context, arg gen.UpdatePasswordParams) error {
	return f.updatePasswordFn(ctx, arg)
}

func (f *fakeQuerier) UpdateUser(ctx context.Context, arg gen.UpdateUserParams) (gen.User, error) {
	return f.updateUserFn(ctx, arg)
}
