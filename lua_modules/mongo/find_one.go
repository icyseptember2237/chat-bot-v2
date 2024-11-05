package mongo

import (
	"chatbot/logger"
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func findOne(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lFilter := state.Get(constant.Param2)
	if lFilter.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type error: filter must be table")
		return 0
	}

	filter := gluamapper.ToGoValue(lFilter, gluamapper.Option{NameFunc: func(s string) string { return s }})
	filter = luatool.ConvertLuaData(filter)

	bsm, err := convertToBsonM(filter)
	if err != nil {
		fmt.Println(err)
		state.ArgError(constant.Param2, "invalid filter format")
		return 0
	}

	var res bson.Raw
	err = conn.col.FindOne(context.Background(), bsm).Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) == false {
			logger.Errorf(context.Background(), "col.FindOne error %s", err.Error())
		}
		state.Push(lua.LNil)
		return 1
	}

	result, err := convertFromBsonRaw(res)
	if err != nil {
		logger.Errorf(context.Background(), "convertFromBsonRaw error %s", err.Error())
		state.Push(lua.LNil)
		return 1
	}
	dst := luatool.ConvertToTable(state, result)
	state.Push(dst)
	return 1
}
