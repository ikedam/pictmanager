package model

type Tag struct {
	ID           string `json:"id"`
	NormalizedTo string `json:"normalizedTo"`
	Count        int    `json:"count"`
}
