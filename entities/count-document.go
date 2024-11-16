package entities

// CountDocument count document
type CountDocument struct {
	NumberOfEntities int64 `json:"number_of_entities" bson:"number_of_entities"`
}
