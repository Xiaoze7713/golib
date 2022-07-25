package redis_instance_tool

import (
	"context"
	"errors"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

type Table struct {
	RedisToolBase
}

func NewTable(businessKey, sep string, client redis.Cmdable) (table *Table, err error) {
	if businessKey == "" {
		//err = errors.New("null business key")
		xlog.Warn("table null business key")
	}
	if client == nil {
		err = errors.New("redis pool is nil")
	}
	table = &Table{
		RedisToolBase{
			client:  client,
			selfKey: businessKey,
			sep:     sep,
		},
	}
	return
}

func (m *Table) SelfKey(s string) (key string) {
	return m.selfKey + m.sep + s
}

func (m *Table) SelfSep() (sep string) {
	return m.sep
}

func (m *Table) GetColumn(ctx context.Context, tableKey, column string) (value interface{}, err error) {
	valueMap, err := m.GetColumns(ctx, tableKey, column)
	if err != nil {
		return nil, err
	}
	value, ok := valueMap[column]
	if !ok {
		return nil, errors.New("value not found")
	}
	return value, nil
}

func (m *Table) SetColumn(ctx context.Context, tableKey, column string, value interface{}) (err error) {
	columnsToValue := map[string]interface{}{column: value}
	err = m.SetColumns(ctx, tableKey, columnsToValue)
	return err
}

func (m *Table) DelColumn(ctx context.Context, tableKey, column string) (err error) {
	err = m.DelColumns(ctx, tableKey, column)
	return err
}

func (m *Table) GetAllColumns(ctx context.Context, tableKey string) (columnsToValue map[string]string, err error) {
	columnsToValue = map[string]string{}
	columnsToValue, err = m.client.HGetAll(ctx, m.SelfKey(tableKey)).Result()
	if err != nil {
		return nil, err
	}
	return columnsToValue, nil
}

func (m *Table) GetColumns(ctx context.Context, tableKey string, columns ...string) (columnsToValue map[string]string, err error) {
	values, err := m.client.HMGet(ctx, m.SelfKey(tableKey), columns...).Result()
	if err != nil {
		return nil, err
	}
	columnsToValue = map[string]string{}
	for idx, value := range values {
		columnsToValue[columns[idx]] = value.(string)
	}
	return columnsToValue, nil
}

func (m *Table) SetColumns(ctx context.Context, tableKey string, columnsToValue map[string]interface{}) (err error) {
	if len(columnsToValue) == 0 {
		return nil
	}
	var cmdArgs []interface{}
	for k, valueIf := range columnsToValue {
		var value = ""
		value, err = jsoniter.MarshalToString(valueIf)
		if err != nil {
			xlog.Errorf("%v", err)
			return err
		}
		cmdArgs = append(cmdArgs, k, value)
	}
	_, err = m.client.HMSet(ctx, m.SelfKey(tableKey), cmdArgs...).Result()
	return
}

func (m *Table) DelColumns(ctx context.Context, tableKey string, columns ...string) (err error) {
	var cmdArgs []interface{}
	for _, column := range columns {
		cmdArgs = append(cmdArgs, column)
	}
	var valueFields = []interface{}{m.SelfKey(tableKey)}
	valueFields = append(valueFields, cmdArgs...)
	xlog.Debugf("redis input cmd %v", valueFields)
	_, err = m.client.HDel(ctx, m.SelfKey(tableKey), columns...).Result()
	if err != nil {
		return err
	}
	return nil
}

func (m *Table) Del(ctx context.Context, tableKey string) (err error) {
	_, err = m.client.Del(ctx, m.SelfKey(tableKey)).Result()
	return
}
