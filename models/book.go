package models

import (
	"encoding/json"
	"time"
)

const ctLayout = "2006-01-02T00:00:00Z"

// TODO (ASK) - where to put this?
type CustomTime struct {
	time.Time
}

var _ json.Unmarshaler = &CustomTime{}

func (ct *CustomTime) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	time, err := time.ParseInLocation(ctLayout, s, time.UTC)
	if err != nil {
		return err
	}
	*ct = CustomTime{time}
	return nil
}

func (ct *CustomTime) marshalJSON() ([]byte, error) {
	s := ct.Time.Format(ctLayout)
	return []byte(s), nil

}

type Book struct {
	Title       string     `json:"title" binding:"required"`
	AuthorName  string     `json:"author_name" binding:"required"`
	Price       float64    `json:"price" binding:"required"`
	Ebook       bool       `json:"ebook, ebook_available"`
	PublishDate CustomTime `json:"publish_date" binding:"required"`
}

type UpdateBook struct {
	Id          string     `json:"id" binding:"required`
	Title       string     `json:"title" binding:"required"`
}
