package repository_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"golang.org/x/net/context"

	"github.com/nutstick/nithi-backend/database/postgresql"
	"github.com/nutstick/nithi-backend/logging"
	"github.com/nutstick/nithi-backend/model"
	"github.com/nutstick/nithi-backend/packages/admin"
	. "github.com/nutstick/nithi-backend/packages/admin/repository"
	"github.com/nutstick/nithi-backend/utiltest"
)

var _ = Describe("admin.SQLRepository", func() {
	// Setup test dependencies
	var (
		app      *fxtest.App
		mockCtrl *gomock.Controller
		repo     admin.Repository
		conn     *postgresql.Connection
		ctx      context.Context
	)

	BeforeEach(func() {
		app = fxtest.New(GinkgoT(),
			fx.Provide(utiltest.NewPostgresQLTestVariable),
			fx.Provide(logging.New),
			fx.Provide(gomock.NewController),
			fx.Provide(utiltest.NewTestReporter),
			fx.Provide(postgresql.New),
			fx.Provide(NewSQLRepository),
			fx.Populate(&mockCtrl),
			fx.Populate(&conn),
			fx.Populate(&repo),
		)
		ctx = conn.WithContext(context.Background())
		app.RequireStart()
		repo.Delete(ctx)
	})

	AfterEach(func() {
		if err := conn.Client().DropTable(&model.Admin{}).Error; err != nil {
			panic(err)
		}
		defer app.RequireStop()
		defer mockCtrl.Finish()
	})

	Describe(".Create", func() {
		Context("with empty model", func() {
			It("should be successfully created", func() {
				admin, err := repo.Create(ctx, &model.Admin{})
				Ω(err).To(BeNil())
				Ω(admin).ToNot(BeNil())
				Ω(string(admin.ID)).To(Equal("1"))
				Ω(admin.Roles).To(Equal(model.StringArray([]string{})))
			})
			It("should be successfully created", func() {
				admin, err := repo.Create(ctx, &model.Admin{
					Email:    "test@nithi.io",
					Password: "abc",
					Name:     "test",
				})
				Ω(err).To(BeNil())
				Ω(admin).ToNot(BeNil())
				Ω(string(admin.ID)).To(Equal("1"))
				Ω(admin.Email).To(Equal("test@nithi.io"))
				Ω(admin.Password).NotTo(Equal("abc"))
				Ω(admin.Name).To(Equal("test"))
				Ω(admin.Roles).To(Equal(model.StringArray([]string{})))

				admin, err = repo.Create(ctx, &model.Admin{
					Email:    "test2@nithi.io",
					Password: "def",
					Name:     "test2",
				})
				Ω(err).To(BeNil())
				Ω(admin).ToNot(BeNil())
				Ω(string(admin.ID)).To(Equal("2"))
				Ω(admin.Email).To(Equal("test2@nithi.io"))
				Ω(admin.Password).NotTo(Equal("def"))
				Ω(admin.Name).To(Equal("test2"))
			})
		})
	})
})
