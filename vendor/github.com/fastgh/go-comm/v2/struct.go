package comm

import (
	"github.com/pkg/errors"
	"github.com/tiaotiao/mapstruct"
)

func StructToMap(src any) map[string]any {
	return mapstruct.Struct2Map(src)
}

func Map2StructP(src map[string]any, dest any) {
	err := Map2Struct(src, dest)
	if err != nil {
		panic(err)
	}
}

func Map2Struct(src map[string]any, dest any) error {
	if err := mapstruct.Map2Struct(src, dest); err != nil {
		return errors.Wrap(err, "convert map to struct")
	}
	return nil
}
