package loader

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetuser(t *testing.T) {
	ctx := context.Background()
	t.Run("get user=1", func(t *testing.T) {
		user, err := testLoaders.GetUser(ctx, "1")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Hoge", user.Name)
	})

	t.Run("get user=2", func(t *testing.T) {
		user, err := testLoaders.GetUser(ctx, "2")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "2", user.ID)
		assert.Equal(t, "Fuga", user.Name)
	})

	t.Run("not exist user", func(t *testing.T) {
		user, err := testLoaders.GetUser(ctx, "9999")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
