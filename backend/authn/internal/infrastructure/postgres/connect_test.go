package postgres

import (
	"context"
	"testing"

	"github.com/faber-numeris/beholder/authn/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestConnectErrorsOnBadDSN(t *testing.T) {
	mockCfg := mocks.NewMockIAppConfig(t)
	mockCfg.EXPECT().DBHost().Return("!invalid!").Maybe()
	mockCfg.EXPECT().DBPort().Return(0).Maybe()
	mockCfg.EXPECT().DBUser().Return("").Maybe()
	mockCfg.EXPECT().DBPassword().Return("").Maybe()
	mockCfg.EXPECT().DBName().Return("").Maybe()
	mockCfg.EXPECT().DBSSLMode().Return("").Maybe()

	pool, err := Connect(context.Background(), mockCfg)
	assert.Error(t, err)
	assert.Nil(t, pool)
}
