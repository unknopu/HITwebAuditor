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
