package models

import (
	"gopkg.in/mgo.v2/bson"
)

type job struct {
	Id      bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Title   string        `json:"title,omitempty"`
	Date    string        `json:"date,omitempty"`
	Author  string        `json:"author,omitempty"`
	Content string        `json:"content,omitempty"`
	View    int           `json:"view,omitempty"`
}

func JobList(page int) []job {
	result := []job{}
	Jobs.Find(bson.M{}).
		Select(bson.M{"content": 0, "_id": 0}).
		Limit(10).
		Sort("-date").
		Skip(page * 10).
		All(&result)
	return result
}

func JobDetail(title, date string) []job {
	db := Jobs.Find(bson.M{"title": title, "date": date}).Select(bson.M{"_id": 0})
	result := []job{}
	db.All(&result)
	return result
}

func JobUpView(title, date string) bool {
	err := Jobs.Update(bson.M{"title": title, "date": date},
		bson.M{"$inc": bson.M{"view": 1}})
	if err != nil {
		return false
	}
	return true
}
