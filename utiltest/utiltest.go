package utiltest

import (
	"os"

	"github.com/golang/mock/gomock"
	"github.com/onsi/ginkgo"
	"go.uber.org/fx"
)

func NewTestReporter() gomock.TestReporter {
	return ginkgo.GinkgoT()
}

type postgresQLRepositoryResult struct {
	fx.Out
	Production bool `name:"PRODUCTION"`
	// PostgresQL
	PostgresQLURL string `name:"postgresql_url"`
	Environment   string `name:"env"`
}

func NewPostgresQLTestVariable() (postgresQLRepositoryResult, error) {
	url := "host=127.0.0.1 port=5432 user=postgres dbname=gqlgen-clean-example-test sslmode=disable"
	if os.Getenv("POSTGRESQL_URL") != "" {
		url = os.Getenv("POSTGRESQL_URL")
	}
	return postgresQLRepositoryResult{
		Production:    false,
		PostgresQLURL: url,
		Environment:   "test",
	}, nil
}

type mongoRepositoryResult struct {
	fx.Out
	Production bool `name:"PRODUCTION"`
	// Mongodb variables
	MongoURL      string `name:"mongo_url"`
	MongoDatabase string `name:"mongo_database"`
	Environment   string `name:"env"`
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
	database := "gqlgen-clean-example-test"
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
