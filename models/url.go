package models

// MyURL the properties in mongodb document
type MyURL struct {
	ID  string `bson:"_id" json:"id"`
	URL string `bson:"url" json:"url"`
}
