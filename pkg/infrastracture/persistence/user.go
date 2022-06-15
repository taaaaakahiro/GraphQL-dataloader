package persistence

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/graph-gophers/dataloader"

	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/entity"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/repository"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/model"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
)

type UserRepo struct {
	database *io.SQLDatabase
}

var _ repository.IUserRepository = (*UserRepo)(nil)

func NewUserReopsitory(db *io.SQLDatabase) *UserRepo {
	return &UserRepo{
		database: db,
	}
}

// GetUsers ref) https://gqlgen.com/reference/dataloaders/
func (r UserRepo) GetUsers(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
	output := make([]*dataloader.Result, len(keys))

	userIds := make([]interface{}, len(keys))
	for i, key := range keys {
		userId, err := strconv.Atoi(key.String())
		if err != nil {
			log.Printf("%+v", err)
			err := fmt.Errorf("user error %s", err.Error())
			output[0] = &dataloader.Result{Data: nil, Error: err}
			return output
		}
		userIds[i] = userId
	}
	query := "SELECT id, name FROM user WHERE id IN (?" + strings.Repeat(",?", len(userIds)-1) + ");"
	stmtOut, err := r.database.Prepare(query)
	if err != nil {
		log.Printf("%+v", err)
		err := fmt.Errorf("user error %s", err.Error())
		output[0] = &dataloader.Result{Data: nil, Error: err}
		return output
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userIds...)
	if err != nil {
		err := fmt.Errorf("user error %s", err.Error())
		output[0] = &dataloader.Result{Data: nil, Error: err}
		return output
	}

	userById := map[string]*model.User{}
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		user := entity.User{}

		err = rows.Scan(&user.Id, &user.Name)
		if err != nil {
			log.Printf("%+v", err)
			err := fmt.Errorf("user error %s", err.Error())
			output[0] = &dataloader.Result{Data: nil, Error: err}
			return output
		}
		modelUser := model.User{
			ID:   strconv.Itoa(user.Id),
			Name: user.Name,
		}
		userById[modelUser.ID] = &modelUser
	}
	if err = rows.Err(); err != nil {
		log.Printf("%+v", err)
		err := fmt.Errorf("user error %s", err.Error())
		output[0] = &dataloader.Result{Data: nil, Error: err}
		return output
	}

	for index, userKey := range keys {
		user, ok := userById[userKey.String()]
		if ok {
			output[index] = &dataloader.Result{Data: user, Error: nil}
		} else {
			err := fmt.Errorf("user not found %s", userKey.String())
			output[index] = &dataloader.Result{Data: nil, Error: err}
			// HACK:
			// ここで、いわゆるDBの外部結合的に、該当レコードがなかったとしても、ダミー値をセットしてエラーをを返却したくない場合は、
			// 下記のようにでダミー値をセットしたDataインスタンスをセット& Errorはnilすることで、親モデルが全部エラーにならないように回避できる
			//dummy := &model.User{
			//	ID:   "",
			//	Name: "unknown",
			//}
			//output[index] = &dataloader.Result{Data: dummy, Error: nil}
		}
	}
	return output
}
