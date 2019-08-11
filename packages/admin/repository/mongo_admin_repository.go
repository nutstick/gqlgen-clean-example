package repository

import (
	"log"
	"time"

	"github.com/nutstick/gqlgen-clean-example/database/mongodb"
	"github.com/nutstick/gqlgen-clean-example/model"
	"github.com/nutstick/gqlgen-clean-example/packages/admin"
	"github.com/nutstick/gqlgen-clean-example/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

const (
	// Collection name that store admins
	Collection = "admins"
)

// mongoRepository contains all the interactions
// with the admin collection stored in mongo.
type mongoRepository struct {
	db     string
	logger *zap.Logger
}

// MongoRepositoryTarget is `fx.In` struct for `fx` to get all dependency to create `AdminMongoRepository`
type MongoRepositoryTarget struct {
	fx.In
	MongoDatabase string `name:"mongo_database"`
	Logger        *zap.Logger
}

// NewMongoRepository is AdminRepository's constructor
func NewMongoRepository(target MongoRepositoryTarget) admin.Repository {
	return &mongoRepository{
		db:     target.MongoDatabase,
		logger: target.Logger,
	}
}

// Collection method extract mongo session from context and return mgo.Collection of this repository
func (m *mongoRepository) Collection(ctx context.Context) *mongo.Collection {
	return mongodb.ForContext(ctx).Database(m.db).Collection(Collection)
}

// GetAll returns all the admins stored in the database.
func (m *mongoRepository) GetAll(ctx context.Context) ([]*model.Admin, error) {
	var admins []*model.Admin
	cursor, err := m.Collection(ctx).Find(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var admin *model.Admin
		err := cursor.Decode(&admin)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}
	return admins, err
}

// GetByID returns one admin which is matched by input ID from the database.
func (m *mongoRepository) GetByID(ctx context.Context, id model.ID) (*model.Admin, error) {
	var admin model.Admin
	err := m.Collection(ctx).FindOne(ctx, bson.M{"_id": id}).Decode(&admin)
	return &admin, err
}

// GetByEmail returns one admin which is matched by email
func (m *mongoRepository) GetByEmail(ctx context.Context, email string) (*model.Admin, error) {
	var admin model.Admin
	err := m.Collection(ctx).FindOne(ctx, bson.M{"email": email}).Decode(&admin)
	return &admin, err
}

// Create will insert new admin into database
func (m *mongoRepository) Create(ctx context.Context, admin *model.Admin) (*model.Admin, error) {
	admin.ID = model.ID(primitive.NewObjectID().Hex())
	hashedPassword, err := hashPassword(admin.Password)
	if err != nil {
		return nil, err
	}
	admin.Password = hashedPassword
	admin.Roles = []string{}
	admin.CreateAt = time.Now()
	admin.UpdateAt = time.Now()
	_, err = m.Collection(ctx).InsertOne(ctx, admin)
	return admin, err
}

// Update will update admin by id
func (m *mongoRepository) Update(ctx context.Context, id model.ID, update *model.Admin) (*model.Admin, error) {
	var admin model.Admin
	if err := m.Collection(ctx).FindOne(ctx, bson.M{"_id": id}).Decode(&admin); err != nil {
		return nil, err
	}
	admin.UpdateAt = time.Now()
	if err := utils.Merge(&admin, *update); err != nil {
		return nil, err
	}
	_, err := m.Collection(ctx).UpdateOne(ctx, bson.M{"_id": id}, admin)
	return &admin, err
}

// Delete will remove all admins
func (m *mongoRepository) Delete(ctx context.Context) error {
	_, err := m.Collection(ctx).DeleteMany(ctx, bson.M{})
	return err
}

// DeleteByID will remove admin by id from database
func (m *mongoRepository) DeleteByID(ctx context.Context, id model.ID) error {
	_, err := m.Collection(ctx).DeleteOne(ctx, bson.M{"_id": id})
	return err
}
