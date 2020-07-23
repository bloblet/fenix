package models

type JsonSerializable interface {
	ToJson() map[string]interface{}
	FromJson(map[string]interface{} ) interface{}
}