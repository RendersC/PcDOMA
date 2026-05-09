package mongorepo

import (
	"context"
	"time"

	"github.com/pcdoma/notification-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepository struct {
	col *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{col: db.Collection("notifications")}
}

func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	n.ID = primitive.NewObjectID()
	n.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, n)
	return err
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]*domain.Notification, int64, error) {
	filter := bson.M{"user_id": userID}
	if unreadOnly {
		filter["is_read"] = false
	}

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

	var notifications []*domain.Notification
	cursor.All(ctx, &notifications)
	return notifications, total, nil
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"is_read": true}})
	return err
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.col.UpdateMany(ctx, bson.M{"user_id": userID, "is_read": false}, bson.M{"$set": bson.M{"is_read": true}})
	return err
}

func (r *NotificationRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
