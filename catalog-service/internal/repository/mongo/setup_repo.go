package mongorepo

import (
	"context"
	"time"

	"github.com/pcdoma/catalog-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SetupRepository struct {
	col *mongo.Collection
}

func NewSetupRepository(db *mongo.Database) *SetupRepository {
	return &SetupRepository{col: db.Collection("setups")}
}

func (r *SetupRepository) Create(ctx context.Context, s *domain.Setup) error {
	s.ID = primitive.NewObjectID()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, s)
	return err
}

func (r *SetupRepository) GetByID(ctx context.Context, id string) (*domain.Setup, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var s domain.Setup
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&s); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SetupRepository) List(ctx context.Context, category string, limit, offset int) ([]*domain.Setup, int64, error) {
	query := bson.M{}
	if category != "" {
		query["category"] = category
	}

	total, _ := r.col.CountDocuments(ctx, query)

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "rating_avg", Value: -1}})

	cursor, err := r.col.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var setups []*domain.Setup
	if err := cursor.All(ctx, &setups); err != nil {
		return nil, 0, err
	}
	return setups, total, nil
}

func (r *SetupRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
