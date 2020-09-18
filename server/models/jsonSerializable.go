package models

type JSONSerializable interface {
	ToJson() map[string]interface{}
	FromJson(map[string]interface{}) interface{}
}
