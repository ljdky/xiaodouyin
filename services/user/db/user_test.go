package db

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	Init()
	DB = DB.Debug()
	m.Run()
}

func TestQueryUserByUsername(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryUserByUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryUserByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddUser(t *testing.T) {

	// type args struct {
	// 	ctx  context.Context
	// 	user *User
	// }

	// type test struct {
	// 	name    string
	// 	args    args
	// 	wantErr bool
	// }

	// tests := make([]test, 0, 1000000)

	users := make([]*User, 0, 1000000)

	for i := 0; i < 1000000; i++ {
		username := "test" + fmt.Sprint(i)
		user := &User{
			Username: username,
			Password: "123456",
		}

		users = append(users, user)

		// tests = append(tests, test{
		// 	name: username,

		// 	args: args{
		// 		ctx:  context.Background(),
		// 		user: user,
		// 	},
		// 	wantErr: false,
		// })
	}

	t.Run("create1M", func(t *testing.T) {
		if err := DB.WithContext(context.Background()).CreateInBatches(users, 1000); err != nil {
			t.Errorf("AddUser() error = %v, wantErr %v", err, false)
		}
	})

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		if err := AddUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
	// 			t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
	// 		}
	// 	})
	// }
}
