package mongorepo

import (
	"context"
	"strings"
	"time"

	"github.com/pcdoma/catalog-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PCRepository struct {
	col *mongo.Collection
}

func NewPCRepository(db *mongo.Database) *PCRepository {
	return &PCRepository{col: db.Collection("pcs")}
}

func (r *PCRepository) Create(ctx context.Context, pc *domain.PC) error {
	pc.ID = primitive.NewObjectID()
	pc.Slug = strings.ToLower(strings.ReplaceAll(pc.Name, " ", "-"))
	pc.Currency = "KZT"
	pc.Status = domain.StatusAvailable
	pc.CreatedAt = time.Now()
	pc.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, pc)
	return err
}

func (r *PCRepository) GetByID(ctx context.Context, id string) (*domain.PC, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var pc domain.PC
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&pc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &pc, nil
}

func (r *PCRepository) List(ctx context.Context, filter domain.PCFilter) ([]*domain.PC, int64, error) {
	query := bson.M{}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if filter.Location != "" {
		query["location.district"] = bson.M{"$regex": filter.Location, "$options": "i"}
	}
	if filter.MinPrice > 0 {
		query["price_per_hour"] = bson.M{"$gte": filter.MinPrice}
	}
	if filter.MaxPrice > 0 {
		if existing, ok := query["price_per_hour"].(bson.M); ok {
			existing["$lte"] = filter.MaxPrice
		} else {
			query["price_per_hour"] = bson.M{"$lte": filter.MaxPrice}
		}
	}

	total, _ := r.col.CountDocuments(ctx, query)

	opts := options.Find().
		SetLimit(int64(filter.Limit)).
		SetSkip(int64(filter.Offset))

	switch filter.Sort {
	case "price_asc":
		opts.SetSort(bson.D{{Key: "price_per_hour", Value: 1}})
	case "price_desc":
		opts.SetSort(bson.D{{Key: "price_per_hour", Value: -1}})
	case "rating":
		opts.SetSort(bson.D{{Key: "rating_avg", Value: -1}})
	default:
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	cursor, err := r.col.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var pcs []*domain.PC
	if err := cursor.All(ctx, &pcs); err != nil {
		return nil, 0, err
	}
	return pcs, total, nil
}

func (r *PCRepository) Update(ctx context.Context, id string, update *domain.UpdatePCRequest) (*domain.PC, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	set := bson.M{"updated_at": time.Now()}
	if update.Name != "" {
		set["name"] = update.Name
	}
	if update.Specs != nil {
		set["specs"] = update.Specs
	}
	if update.PricePerHour > 0 {
		set["price_per_hour"] = update.PricePerHour
	}
	if update.PricePerDay > 0 {
		set["price_per_day"] = update.PricePerDay
	}
	if update.Location != nil {
		set["location"] = update.Location
	}
	if update.Images != nil {
		set["images"] = update.Images
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var pc domain.PC
	err = r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": set}, opts).Decode(&pc)
	if err != nil {
		return nil, err
	}
	return &pc, nil
}

func (r *PCRepository) UpdateStatus(ctx context.Context, id string, status domain.PCStatus) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"status": status, "updated_at": time.Now()},
	})
	return err
}

func (r *PCRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *PCRepository) UpdateRating(ctx context.Context, id string, avg float64, count int) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"rating_avg": avg, "rating_count": count, "updated_at": time.Now()},
	})
	return err
}

func (r *PCRepository) GetLocations(ctx context.Context) ([]string, error) {
	result, err := r.col.Distinct(ctx, "location.district", bson.M{})
	if err != nil {
		return nil, err
	}
	locations := make([]string, 0, len(result))
	for _, v := range result {
		if s, ok := v.(string); ok {
			locations = append(locations, s)
		}
	}
	return locations, nil
}
