package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/graph-gophers/dataloader"
	errs "github.com/pkg/errors"

	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/entity"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/repository"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/model"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
)

type UserRepo struct {
	database *io.SQLDatabase
}

var _ repository.IUserRepository = (*UserRepo)(nil)

func NewUserRepository(db *io.SQLDatabase) *UserRepo {
	return &UserRepo{
		database: db,
	}
}

func (r UserRepo) User(userId int) (entity.User, error) {
	user := entity.User{}

	query := "SELECT id, name FROM user WHERE id = ?"
	stmtOut, err := r.database.Prepare(query)
	if err != nil {
		return user, err
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(userId).Scan(&user.Id, &user.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Println("row not found")
			return user, err
		default:
			return user, err
		}
	}
	return user, nil
}

func (r UserRepo) ListUsers() ([]entity.User, error) {
	query := "SELECT id, name FROM user ORDER BY id DESC"
	stmtOut, err := r.database.Prepare(query)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		return nil, errs.WithStack(err)
	}

	users := make([]entity.User, 0)
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		user := entity.User{}

		err = rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return nil, errs.WithStack(err)
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil

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
