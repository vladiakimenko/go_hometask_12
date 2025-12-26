package database

import (
	"fmt"
	"reflect"
	"strings"
)

// Дескрипторы таблиц в бд
type dbTable struct {
	name           string
	pk             string
	dataColumns    []string
	defaultColumns []string
}

var usersTable = dbTable{
	name: "users",
	pk:   "id",
	dataColumns: []string{
		"email",
		"username",
		"password_hash",
	},
	defaultColumns: []string{
		"created_at",
	},
}

// helper'ы для генерации безопасных запросов для любых таблиц
func safeInsertQuery(table dbTable) string {
	placeholders := make([]string, len(table.dataColumns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		table.name,
		strings.Join(table.dataColumns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(append(append([]string{table.pk}, table.dataColumns...), table.defaultColumns...), ", "),
	)
}

func safeLookupQuery(table dbTable) string {
	return fmt.Sprintf(
		"SELECT * FROM %s WHERE %s = $1",
		table.name,
		table.pk,
	)
}

func safeFilterQuery(table dbTable, instance any, existsMode bool) string {
	desiredObject := reflect.ValueOf(instance)
	if desiredObject.Kind() != reflect.Ptr || desiredObject.Elem().Kind() != reflect.Struct {
		panic("instance is not a pointer to a struct")
	}
	desiredObject = desiredObject.Elem()
	var conditions []string
	for i := 0; i < desiredObject.NumField(); i++ {
		fieldValue := desiredObject.Field(i)
		if fieldValue.IsZero() {
			continue
		}
		column := desiredObject.Type().Field(i).Name
		conditions = append(
			conditions,
			fmt.Sprintf("%s = $%d", column, len(conditions)+1),
		)
	}
	if len(conditions) == 0 {
		panic("no non-zero fields provided for filtering")
	}
	whereClause := strings.Join(conditions, " AND ")
	if existsMode {
		return fmt.Sprintf(
			"SELECT EXISTS (SELECT 1 FROM %s WHERE %s)",
			table.name,
			whereClause,
		)
	}
	return fmt.Sprintf(
		"SELECT * FROM %s WHERE %s",
		table.name,
		whereClause,
	)
}
