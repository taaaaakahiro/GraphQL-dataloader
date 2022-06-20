package persistence

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/graph-gophers/dataloader"
	"github.com/stretchr/testify/assert"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/entity"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/model"
)

func TestUserRepo_Listusers(t *testing.T) {
	t.Run("get all users", func(t *testing.T) {
		users, err := userRepo.ListUsers()

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.NotEmpty(t, users)
		assert.Len(t, users, 2)

		want := []entity.User{
			{Id: 2, Name: "Fuga"},
			{Id: 1, Name: "Hoge"},
		}

		if diff := cmp.Diff(want, users); len(diff) != 0 {
			t.Errorf("Users mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestUserRepo_GetUsers(t *testing.T) {
	t.Run("get user=1", func(t *testing.T) {
		ctx := context.Background()
		keys := dataloader.NewKeysFromStrings([]string{"1"})
		result := userRepo.GetUsers(ctx, keys)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 1)

		assert.Nil(t, result[0].Error)
		user := result[0].Data.(*model.User)
		assert.NotNil(t, user)
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Hoge", user.Name)
	})

}
