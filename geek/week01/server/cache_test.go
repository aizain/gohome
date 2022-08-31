package server

import (
	"reflect"
	"testing"
	"time"
)

func TestSyncCache_Get(t *testing.T) {
	type fields struct {
		Cache Cache
		DB    DB
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Data
		want1  bool
	}{
		{"测试获取缓存",
			fields{
				map[string]*CacheData{
					"key": {
						&Data{Name: "wangzihao", Age: 18},
						time.Now().Unix() + time.Minute.Milliseconds(),
					},
				},
				map[string]*Data{},
			},
			args{key: "key"},
			&Data{Name: "wangzihao", Age: 18},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncCache{
				Cache: tt.fields.Cache,
				DB:    tt.fields.DB,
			}
			got, got1 := s.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
