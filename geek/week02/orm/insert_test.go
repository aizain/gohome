package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInserter_buildByReflect(t *testing.T) {
	type Addr struct {
		Desc string
	}
	type User struct {
		Name  string
		Age   int64
		Money float64
		Empty string
		Addr
	}

	user := User{
		Name:  "zain",
		Age:   18,
		Money: 100,
		Addr: Addr{
			Desc: "北京",
		},
	}
	tests := []struct {
		name      string
		inserter  *Inserter[User]
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "build insert one",
			inserter: &Inserter[User]{
				values: []*User{
					&user,
				},
			},
			wantQuery: &Query{
				SQL:  "INSERT INTO `user`(name,age,money,empty,desc) VALUES(?,?,?,?,?)",
				Args: []any{user.Name, user.Age, user.Money, user.Empty, user.Addr.Desc},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.inserter.buildByReflect()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantQuery, got)
		})
	}
}
