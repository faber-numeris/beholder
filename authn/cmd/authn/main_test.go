package main

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/caarlos0/env/v11"
	"github.com/faber-numeris/beholder/authn/internal/app/bootstrap"
	"github.com/faber-numeris/foundation/testutils/envutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainFunction(t *testing.T) {
	t.Run("Expect an error for not setting environment variables", func(t *testing.T) {
		// Snapshot and clear the ambient environment so config parsing fails on
		// the required variables; envutils.Pop restores it once the test returns.
		require.NoError(t, envutils.Push())
		defer func() {
			require.NoError(t, envutils.Pop())
			if r := recover(); r != nil {
				convErr, ok := r.(error)
				require.True(t, ok)
				assert.ErrorIs(t, convErr, env.VarIsNotSetError{})
			}
		}()
		main()
	})

	t.Run("Verify that application initialization returns an error", func(t *testing.T) {
		var app bootstrap.App
		patches := gomonkey.NewPatches()
		patches.ApplyMethodReturn(&app, "Run", assert.AnError)
		patches.ApplyFuncReturn(bootstrap.NewApp, &app)

		defer func() {
			if r := recover(); r != nil {
				_, ok := r.(error)
				require.True(t, ok)
			}
		}()

		main()
	})

}
