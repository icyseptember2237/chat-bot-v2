package mongo

import (
	"chatbot/utils/constant"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/mongo"
)

const metaName = "mongo{meta}"

type Connection struct {
	name       string
	database   string
	collection string
	col        *mongo.Collection
}

func checkConnection(state *lua.LState) *Connection {
	ud := state.Get(constant.Param1)
	if ud.Type() != lua.LTUserData {
		state.ArgError(constant.Param1, "connection expected")
		return nil
	}

	if v, ok := ud.(*lua.LUserData).Value.(*Connection); ok {
		return v
	}
	state.ArgError(constant.Param1, "connection expected")
	return nil
}

var clientExports = map[string]lua.LGFunction{
	"insertOne":  insertOne,
	"insertMany": insertMany,
	"deleteOne":  deleteOne,
	"deleteMany": deleteMany,
	"updateOne":  updateOne,
	"updateMany": updateMany,
	"findOne":    findOne,
	"find":       find,
	"count":      count,
}
