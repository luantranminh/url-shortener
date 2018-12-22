package models

// MyURL the properties in mongodb document
type MyURL struct {
	ID       string `bson:"_id" json:"id"`
	ShortURL string `bson:"short_url" json:"short_url"`
	LongURL  string `bson:"long_url" json:"long_url"`
}
