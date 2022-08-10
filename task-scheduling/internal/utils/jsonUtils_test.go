package main

import (
	"reflect"
	"testing"
)

func TestStruct2Json(t *testing.T) {
	type args struct {
		stru any
	}
	tests := []struct {
		name string
		args args
		want string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Struct2Json(tt.args.stru); got != tt.want {
				t.Errorf("Struct2Json() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson2Struct(t *testing.T) {
	type args struct {
		j []byte
		v any
	}
	type User struct {
		UserName string
		Password string
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "j2s",
			args: args{
				j: []byte("{\"UserName\":\"congziqi\", \"Password\":\"123\"}"),
				v: new(User),
			},
			want: &User{
				UserName: "congziqi",
				Password: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Json2Struct(tt.args.j, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json2Struct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson2Map(t *testing.T) {
	type args struct {
		j []byte
		m *map[string]any
	}
	maps := make(map[string]any)
	maps2 := make(map[string]any)
	maps2["UserName"] = "congziqi"
	maps2["Password"] = "123"
	type User struct {
		UserName string
		Password string
	}
	maps2["user"] = User{UserName: "congziqi", Password: "123"}
	tests := []struct {
		name string
		args args
		want *map[string]any
	}{
		{
			name: "j2m",
			args: args{
				j: []byte("{\"UserName\":\"congziqi\", \"Password\":\"123\",\"user\":{\"UserName\":\"congziqi\",\"Password\":\"123\"}}"),
				m: &maps,
			},
			want: &maps2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Json2Map(tt.args.j, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json2Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap2Json(t *testing.T) {
	type args struct {
		m *map[string]any
	}
	maps2 := make(map[string]any)
	maps2["UserName"] = "congziqi"
	maps2["Password"] = "123"
	type User struct {
		UserName string
		Password string
	}
	maps2["user"] = User{"congziqi", "123"}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "m2j",
			args: args{
				&maps2,
			},
			want: "{\"UserName\":\"congziqi\", \"Password\":\"123\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map2Json(tt.args.m); got != tt.want {
				t.Errorf("Map2Json() = %v, want %v", got, tt.want)
			}
		})
	}
}
