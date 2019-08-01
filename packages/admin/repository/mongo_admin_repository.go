package repository

import (
	"log"
	"time"

	"github.com/nutstick/nithi-backend/database/mongodb"
	"github.com/nutstick/nithi-backend/model"
	"github.com/nutstick/nithi-backend/packages/admin"
	"github.com/nutstick/nithi-backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
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
	return mongodb.ForContext(ctx).Database(m.db).Collection("admins")
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
func (m *mongoRepository) GetByID(ctx context.Context, id string) (*model.Admin, error) {
	var admin model.Admin
	err := m.Collection(ctx).FindOne(ctx, bson.M{"_id": bson.ObjectIdHex(id)}).Decode(&admin)
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
	admin.ID = bson.NewObjectId().Hex()
	admin.CreateAt = time.Now()
	admin.UpdateAt = time.Now()
	_, err := m.Collection(ctx).InsertOne(ctx, admin)
	return admin, err
}

// Update will update admin by id
func (m *mongoRepository) Update(ctx context.Context, id string, update *model.Admin) (*model.Admin, error) {
	var admin model.Admin
	if err := m.Collection(ctx).FindOne(ctx, bson.M{"_id": bson.ObjectIdHex(id)}).Decode(&admin); err != nil {
		return nil, err
	}
	admin.UpdateAt = time.Now()
	if err := utils.Merge(&admin, *update); err != nil {
		return nil, err
	}
	_, err := m.Collection(ctx).UpdateOne(ctx, bson.M{"_id": bson.ObjectIdHex(id)}, admin)
	return &admin, err
}

// Delete will remove admin by id from database
func (m *mongoRepository) Delete(ctx context.Context, id string) error {
	_, err := m.Collection(ctx).DeleteOne(ctx, bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}
