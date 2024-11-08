package mongo

import (
	"chatbot/utils/constant"
	"chatbot/utils/luatool"
	"context"
	lua "github.com/yuin/gopher-lua"
	"go.mongodb.org/mongo-driver/mongo"
)

const cursorMetaName = "cursor{meta}"

var cursorExports = map[string]lua.LGFunction{
	"next":         next,
	"all":          cursorAll,
	"close":        cursorClose,
	"setBatchSize": setBatchSize,
}

func next(state *lua.LState) int {
	cursor := checkCursor(state)
	if cursor == nil {
		state.Push(lua.LNil)
		return 0
	}

	if cursor.Next(context.Background()) {
		data := cursor.Current

		row, err := convertFromBsonRaw(data)
		if err != nil {
			state.Push(lua.LNil)
			return 0
		}

		lRow := luatool.ConvertToTable(state, row)
		state.Push(lRow)
		return 1
	}

	state.Push(lua.LNil)
	return 1
}

func cursorAll(state *lua.LState) int {
	cursor := checkCursor(state)
	if cursor == nil {
		state.Push(lua.LNil)
		return 0
	}

	lResults := state.NewTable()
	for cursor.Next(context.Background()) {
		doc := cursor.Current
		result, err := convertFromBsonRaw(doc)
		if err != nil {
			state.Push(lua.LFalse)
			cursor.Close(context.Background())
			return 1
		}
		dst := luatool.ConvertToTable(state, result)
		lResults.Append(dst)
	}
	state.Push(lResults)
	cursor.Close(context.Background())
	return 1
}

func cursorClose(state *lua.LState) int {
	cursor := checkCursor(state)
	if cursor == nil {
		return 0
	}

	cursor.Close(context.Background())
	return 1
}

func setBatchSize(state *lua.LState) int {
	cursor := checkCursor(state)
	if cursor == nil {
		state.Push(lua.LFalse)
		return 0
	}

	lBatchSize := state.Get(constant.Param2)
	if lBatchSize.Type() != lua.LTNumber {
		state.ArgError(constant.Param2, "type error: batch size must be number")
		return 0
	}
	batchSize := int32(lBatchSize.(lua.LNumber))

	cursor.SetBatchSize(batchSize)

	state.Push(lua.LTrue)
	return 1
}

func checkCursor(state *lua.LState) *mongo.Cursor {
	lv := state.Get(constant.Param1)
	if lv.Type() != lua.LTUserData {
		state.ArgError(constant.Param1, "type error: cursor must be userdata")
		return nil
	}

	if cursor, ok := lv.(*lua.LUserData).Value.(*mongo.Cursor); ok {
		if cursor == nil {
			state.ArgError(constant.Param1, "cursor is nil")
		}
		return cursor
	}
	state.ArgError(constant.Param1, "type error: cursor expected")
	return nil
}
