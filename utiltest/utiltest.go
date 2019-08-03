package utiltest

import (
	"os"

	"github.com/golang/mock/gomock"
	"github.com/onsi/ginkgo"
	"go.uber.org/fx"
	"gopkg.in/mgo.v2"
)

func NewTestReporter() gomock.TestReporter {
	return ginkgo.GinkgoT()
}

type mongoRepositoryResult struct {
	fx.Out
	Production bool `name:"PRODUCTION"`
	// Mongodb variables
	MongoURL      string `name:"mongo_url"`
	MongoDatabase string `name:"mongo_database"`
	Environment   string `name:"env"`
	Session       *mgo.Session
}

// type mongoRepositoryTarget struct {
// 	fx.In
// 	MongoDatabase string `name:"mongo_database"`
// }

func NewMongoTestVariable() (mongoRepositoryResult, error) {
	url := "mongodb://localhost:27017"
	if os.Getenv("MONGO_URL") != "" {
		url = os.Getenv("MONGO_URL")
	}
	database := "nithi-backend-test"
	if os.Getenv("MONGO_DATABASE") != "" {
		database = os.Getenv("MONGO_DATABASE")
	}
	return mongoRepositoryResult{
		Production:    false,
		MongoURL:      url,
		MongoDatabase: database,
		Environment:   "test",
	}, nil
}
