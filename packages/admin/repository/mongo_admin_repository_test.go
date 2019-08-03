package repository_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"

	"github.com/nutstick/nithi-backend/database/mongodb"
	"github.com/nutstick/nithi-backend/logging"
	"github.com/nutstick/nithi-backend/model"
	"github.com/nutstick/nithi-backend/packages/admin"
	. "github.com/nutstick/nithi-backend/packages/admin/repository"
	"github.com/nutstick/nithi-backend/utiltest"
)

var _ = Describe("admin.MongoDBRepository", func() {
	// Setup test dependencies
	var (
		app      *fxtest.App
		mockCtrl *gomock.Controller
		repo     admin.Repository
		conn     *mongodb.Connection
		ctx      context.Context
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			fx.Provide(utiltest.NewMongoTestVariable),
			fx.Provide(logging.New),
			fx.Provide(gomock.NewController),
			fx.Provide(utiltest.NewTestReporter),
			fx.Provide(mongodb.New),
			fx.Provide(NewMongoRepository),
			fx.Populate(&mockCtrl),
			fx.Populate(&conn),
			fx.Populate(&repo),
		)
		ctx = conn.WithContext(context.Background())
		app.RequireStart()
		repo.Delete(ctx)
	})

	AfterEach(func() {
		defer app.RequireStop()
		defer mockCtrl.Finish()
	})

	Describe(".Create", func() {
		Context("with empty model", func() {
			It("should be successfully created", func() {
				admin, err := repo.Create(ctx, &model.Admin{})
				Ω(err).To(BeNil())
				Ω(admin).ToNot(BeNil())
				bson.ObjectIdHex(string(admin.ID))
			})
			It("should be successfully created", func() {
				admin, err := repo.Create(ctx, &model.Admin{
					Email:    "test@nithi.io",
					Password: "abc",
					Name:     "test",
				})
				Ω(err).To(BeNil())
				Ω(admin).ToNot(BeNil())
				bson.ObjectIdHex(string(admin.ID))
			})
		})
	})
})
