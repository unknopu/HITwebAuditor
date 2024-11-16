package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func IndexOfObjectID(element primitive.ObjectID, data []primitive.ObjectID) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func IsExisting(source string, payload []string) bool {
	for _, v := range payload {
		if source == v {
			return true
		}
	}
	return false
}
