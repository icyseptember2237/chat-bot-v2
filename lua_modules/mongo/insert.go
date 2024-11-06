package mongo

import (
	"chatbot/logger"
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"context"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/bson/primitive"
	luar "layeh.com/gopher-luar"
)

func insertOne(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lDoc := state.Get(constant.Param2)
	if lDoc.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type err: doc must be of type table")
		return 0
	}
	doc := gluamapper.ToGoValue(lDoc, gluamapper.Option{NameFunc: func(s string) string {
		return s
	}})
	doc = luatool.ConvertLuaData(doc)

	bd, err := convertToBsonM(doc)
	if err != nil {
		state.ArgError(constant.Param2, "invalid data format")
		return 0
	}

	res, err := conn.col.InsertOne(context.Background(), bd)
	if err != nil {
		logger.Errorf(context.Background(), "mongo.insertOne %s error : %v", conn.name, err.Error())
		state.Push(lua.LFalse)
		return 1
	}

	state.Push(lua.LTrue)
	state.Push(luar.New(state, res.InsertedID))
	return 1
}

func insertMany(state *lua.LState) int {
	conn := checkConnection(state)
	if conn == nil {
		return 0
	}

	lDocs := state.Get(constant.Param2)
	if lDocs.Type() != lua.LTTable {
		state.ArgError(constant.Param2, "type err: docs must be of type table")
		return 0
	}

	var docs []interface{}
	lDocs.(*lua.LTable).ForEach(func(value lua.LValue, value2 lua.LValue) {
		data := gluamapper.ToGoValue(value2, gluamapper.Option{NameFunc: gluamapper.ToUpperCamelCase})
		data = luatool.ConvertLuaData(data)

		bd, err := convertToBsonM(data)
		if err != nil {
			return
		}
		docs = append(docs, bd)
	})

	rets, err := conn.col.InsertMany(context.Background(), docs)
	if err != nil {
		logger.Errorf(context.Background(), "mongo.insertOne %s error : %v", conn.name, err.Error())
		state.Push(lua.LFalse)
		return 1
	}

	lRets := lua.LTable{}
	for _, ret := range rets.InsertedIDs {
		lRets.Append(luar.New(state, ret.(primitive.ObjectID)))
	}

	state.Push(lua.LTrue)
	state.Push(&lRets)
	return 2
}
