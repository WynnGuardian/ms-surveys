package util

import (
	"database/sql"
	"strings"
	"victo/wynnguardian/internal/domain/response"
)

type Pair[T any, U any] struct {
	First  T
	Second U
}

func MustVal[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func KeySlice[T comparable, U any](m map[T]U) []T {
	keys := make([]T, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func HighestLength(items []string) int {
	leng := 0
	for _, v := range items {
		if len(v) > leng {
			leng = len(v)
		}
	}
	return leng
}

func PadText(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(" ", length-len(str))
}

func NotFoundOrInternalErr(err error, notFound response.WGResponse) response.WGResponse {
	if err == sql.ErrNoRows {
		return notFound
	}
	return response.ErrInternalServerErr(err)
}

func WrapInt(v int) *int {
	var wrap *int
	wrap = new(int)
	wrap = &v
	return wrap
}
