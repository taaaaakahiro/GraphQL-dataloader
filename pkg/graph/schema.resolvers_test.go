package graph

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/config"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/generated"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/model"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/loader"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/persistence"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
	"go.uber.org/zap"
)

func TestQueryResolver_Users(t *testing.T) {
	t.Run("get all users", func(t *testing.T) {
		resolver := getQueryResolver()
		ctx := context.Background()
		users, err := resolver.Users(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, users)
		assert.Len(t, users, 2)

		want := []*model.User{
			{ID: "2", Name: "Fuga"},
			{ID: "1", Name: "Hoge"},
		}
		if diff := cmp.Diff(want, users); len(diff) != 0 {
			t.Error("users mismatch (-want + got)", diff)
		}

	})
}

func TestMessageResolver_User(t *testing.T) {
	resolver := getResolver()
	ctx := context.Background()
	msgResolver := resolver.Message()

	t.Run("get user=1", func(t *testing.T) {
		message := model.Message{
			UserID: "1",
		}
		user, err := msgResolver.User(ctx, &message)
		assert.NoError(t, err)
		assert.NotEmpty(t, user)
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Hoge", user.Name)
	})
	t.Run("get user=2", func(t *testing.T) {
		message := model.Message{
			UserID: "2",
		}
		user, err := msgResolver.User(ctx, &message)
		assert.NoError(t, err)
		assert.NotEmpty(t, user)
		assert.Equal(t, "2", user.ID)
		assert.Equal(t, "Fuga", user.Name)
	})
	t.Run("get not exist user", func(t *testing.T) {
		message := model.Message{
			UserID: "9999",
		}
		user, err := msgResolver.User(ctx, &message)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestMutationResolver_CreateMessage(t *testing.T) {
	resolver := getMutationResolver()
	ctx := context.Background()

	t.Run("create message user=1", func(t *testing.T) {
		input := model.NewMessage{
			UserID:  "1",
			Message: "new message 1",
		}
		message, err := resolver.CreateMessage(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, message)
		id, _ := strconv.Atoi(message.ID)
		assert.Greater(t, id, 1)
		assert.Equal(t, input.UserID, message.UserID)
		assert.Equal(t, input.Message, message.Message)

	})

}

func getQueryResolver() generated.QueryResolver {
	return getResolver().Query()
}

func getResolver() *Resolver {
	mysqlDatabase := getDatabase()
	repositories, err := persistence.NewReopsitories(mysqlDatabase)
	if err != nil {
		log.Println("failed to new repositories", zap.Error(err))
	}

	loaders := loader.NewLoaders(repositories)
	resolver := &Resolver{
		Repo:    repositories,
		Loaders: loaders,
	}
	return resolver
}

func getDatabase() *io.SQLDatabase {
	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	sqlSetting := &config.SQLDBSetting{
		SqlDSN:              cfg.DB.DSN,
		SqlMaxOpenConns:     cfg.DB.MaxOpenConns,
		SqlMaxIdleConns:     cfg.DB.MaxIdleConns,
		SqlConnsMaxLifetime: cfg.DB.ConnsMaxLifetime,
	}

	db, err := io.NewDatabase(sqlSetting)
	if err != nil {
		log.Println(err.Error())
	}
	return db
}

func getMutationResolver() generated.MutationResolver {
	return getResolver().Mutation()
}
