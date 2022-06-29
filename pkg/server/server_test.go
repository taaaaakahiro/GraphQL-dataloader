package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/config"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/entity"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/generated"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/model"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/handler"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/loader"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/persistence"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
	"go.uber.org/zap"
)

func TestServer(t *testing.T) {
	//create logger
	logger, err := zap.NewProduction()
	if err != nil {
		t.Errorf("failed to setup loggger: %s\n", err)
	}
	defer logger.Sync()
	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		t.Errorf("failed to load config %s\n", err)
	}

	// init mysql
	logger.Info("connect to mysql", zap.String("DSN", cfg.DB.DSN))
	sqlSetting := &config.SQLDBSetting{
		SqlDSN:              cfg.DB.DSN,
		SqlMaxOpenConns:     cfg.DB.MaxOpenConns,
		SqlMaxIdleConns:     cfg.DB.MaxIdleConns,
		SqlConnsMaxLifetime: cfg.DB.ConnsMaxLifetime,
	}
	mysqlDatabase, err := io.NewDatabase(sqlSetting)
	if err != nil {
		t.Errorf("failed to create mysql db repository: %s\n", err)
	}
	repositories, err := persistence.NewReopsitories(mysqlDatabase)
	assert.NoError(t, err)

	// start server
	loaders := loader.NewLoaders(repositories)

	// init to start http server
	// init gql server
	query := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Repo:    repositories,
				Loaders: loaders,
			},
		}))
	registry := handler.NewHandler(logger, repositories, query, "v1.0-test")
	s := NewServer(registry, &Config{Log: logger})
	testServer := httptest.NewServer(s.Mux)
	defer testServer.Close()

	// test API
	t.Run("check /healthz", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/healthz")
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})
	t.Run("check /version", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/version")
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		//read body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)
		assert.NotEmpty(t, body)

		var data interface{}
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err)
		ver := data.(map[string]interface{})["version"].(string)
		assert.Equal(t, ver, "v1.0-test")
	})

	t.Run("check /user/list", func(t *testing.T) {
		res, err := http.Get(testServer.URL + "/user/list")
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		//read body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)
		assert.NotEmpty(t, body)

		var users []entity.User
		err = json.Unmarshal(body, &users)
		assert.NoError(t, err)
		assert.Greater(t, len(users), 0)

		assert.Contains(t, users, entity.User{
			Id:   1,
			Name: "Hoge",
		})
		assert.Contains(t, users, entity.User{
			Id:   2,
			Name: "Fuga",
		})
	})

	t.Run("check /message/1", func(t *testing.T) {

	})
	t.Run("check /message/2", func(t *testing.T) {

	})

	// test GraphQL
	t.Run("check /gql", func(t *testing.T) {

	})
	t.Run("check /query", func(t *testing.T) {

	})
	t.Run("check graphQL query", func(t *testing.T) {
		t.Run("query users", func(t *testing.T) {
			type resp struct {
				Users []model.User `json:"users"`
			}

			t.Run("get id, name", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
						query getusers {
							users {
								id
								name
							}
						}
					`,
				}
				req, _ := json.Marshal(graphReq)
				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotNil(t, gqlRes.Data)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Users)

				want := []model.User{ // 取得データの順番を考慮
					{ID: "2", Name: "Fuga"},
					{ID: "1", Name: "Hoge"},
				}
				if diff := cmp.Diff(want, data.Users); len(diff) != 0 {
					t.Errorf("users mismatch (-want + got):\n%s", diff)
				}
			})

			t.Run("get id", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `query getUsers {
						users {
							id
						}
					}`,
				}
				req, err := json.Marshal(graphReq)
				assert.NoError(t, err)
				res, err := http.Post(testServer.URL+"/query", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close() //忘れ注意
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, body)

				var data resp
				err = json.Unmarshal(gqlRes.Data, &data)
				assert.NoError(t, err)
				assert.NotNil(t, data.Users)

				want := []model.User{
					{ID: "2"},
					{ID: "1"},
				}
				if diff := cmp.Diff(want, data.Users); len(diff) != 0 {
					t.Errorf("users mismatch (-want +got):\n%s", diff)
				}
			})

			t.Run("get name", func(t *testing.T) {
				graphReq := graphql.RawParams{
					Query: `
						query getUsers {
							users {
								name
							}
						}
					`,
				}
				req, _ := json.Marshal(graphReq)
				res, err := http.Post(testServer.URL+"/graph", "application/json", bytes.NewReader(req))
				assert.NoError(t, err)
				// read body
				body, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				assert.NoError(t, err)
				assert.NotEmpty(t, body)

				var gqlRes graphql.Response
				err = json.Unmarshal(body, &gqlRes)
				assert.NoError(t, err)
				assert.Nil(t, gqlRes.Errors)
				assert.NotEmpty(t, gqlRes.Data)

			})
		})
		t.Run("query messages", func(t *testing.T) {
			t.Run("get messages user=1 with user id, name", func(t *testing.T) {

			})
			t.Run("get messages user=2 without id", func(t *testing.T) {

			})
			t.Run("get messages user=2 without users", func(t *testing.T) {

			})
			t.Run("get messages not exist user", func(t *testing.T) {

			})
			t.Run("invalid param", func(t *testing.T) {

			})
			t.Run("lack params", func(t *testing.T) {

			})

		})
		t.Run("query all messages", func(t *testing.T) {
			t.Run("get all messages with user id, name", func(t *testing.T) {

			})
			t.Run("get messages without id", func(t *testing.T) {

			})
			t.Run("get messages with users", func(t *testing.T) {

			})

		})
		t.Run("multiple modules", func(t *testing.T) {
			t.Run("get user=1", func(t *testing.T) {

			})

		})
		t.Run("undefined schema", func(t *testing.T) {

		})

	})

}

func GetDatabase() *io.SQLDatabase {
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
		os.Exit(1)
	}
	return db

}
