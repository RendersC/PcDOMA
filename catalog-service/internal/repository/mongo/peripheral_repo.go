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

type PeripheralRepository struct {
	col *mongo.Collection
}

func NewPeripheralRepository(db *mongo.Database) *PeripheralRepository {
	return &PeripheralRepository{col: db.Collection("peripherals")}
}

func (r *PeripheralRepository) Create(ctx context.Context, p *domain.Peripheral) error {
	p.ID = primitive.NewObjectID()
	p.Currency = "KZT"
	p.Status = domain.StatusAvailable
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, p)
	return err
}

func (r *PeripheralRepository) GetByID(ctx context.Context, id string) (*domain.Peripheral, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var p domain.Peripheral
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&p); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *PeripheralRepository) List(ctx context.Context, peripheralType string, locationID string, limit, offset int) ([]*domain.Peripheral, int64, error) {
	query := bson.M{}
	if peripheralType != "" {
		query["type"] = peripheralType
	}
	if locationID != "" {
		query["location_id"] = locationID
	}

	total, _ := r.col.CountDocuments(ctx, query)

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.col.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var peripherals []*domain.Peripheral
	if err := cursor.All(ctx, &peripherals); err != nil {
		return nil, 0, err
	}
	return peripherals, total, nil
}

func (r *PeripheralRepository) UpdateStatus(ctx context.Context, id string, status domain.PCStatus) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"status": status, "updated_at": time.Now()},
	})
	return err
}

func (r *PeripheralRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
