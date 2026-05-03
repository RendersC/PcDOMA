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

type ReviewRepository struct {
	col *mongo.Collection
}

func NewReviewRepository(db *mongo.Database) *ReviewRepository {
	return &ReviewRepository{col: db.Collection("reviews")}
}

func (r *ReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	review.ID = primitive.NewObjectID()
	review.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, review)
	return err
}

func (r *ReviewRepository) ListByPCID(ctx context.Context, pcID string, limit, offset int) ([]*domain.Review, int64, error) {
	filter := bson.M{"pc_id": pcID}
	total, _ := r.col.CountDocuments(ctx, filter)

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reviews []*domain.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, 0, err
	}
	return reviews, total, nil
}

func (r *ReviewRepository) GetRatingStats(ctx context.Context, pcID string) (avg float64, count int, err error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"pc_id": pcID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"avg":   bson.M{"$avg": "$rating"},
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Avg   float64 `bson:"avg"`
		Count int     `bson:"count"`
	}
	cursor.All(ctx, &result)
	if len(result) > 0 {
		return result[0].Avg, result[0].Count, nil
	}
	return 0, 0, nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
