package mongo

import "time"

type audit struct {
	ProjectName  string    `bson:"project_name"`
	Time         time.Time `bson:"time"`
	ResourceType string    `bson:"resource_type"`
	Method       string    `bson:"method"`
	Path         string    `bson:"path"`
	IsSuccess    bool      `bson:"is_success"`
	Message      string    `bson:"message"`
	UnixTime     int64     `bson:"unixtime"`
}
