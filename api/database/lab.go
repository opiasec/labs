package database

import (
	"appseclabs/types"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *Database) GetLab(labSlug string) (*types.Lab, error) {
	var lab types.Lab

	labsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("labs")

	err := labsCollection.FindOne(context.Background(), bson.M{"slug": labSlug}).Decode(&lab)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("lab not found")
		}
		return nil, err
	}
	return &lab, nil
}

func (d *Database) SaveLabSession(labSession types.LabSession) error {
	labSession.CreatedAt = time.Now()
	labSession.UpdatedAt = time.Now()
	labSession.StartedAt = time.Now()
	labSessionsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("lab_sessions")
	_, err := labSessionsCollection.InsertOne(context.Background(), labSession)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetLabSession(namespace string) (*types.LabSession, error) {
	labSessionsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("lab_sessions")
	var labSession types.LabSession
	err := labSessionsCollection.FindOne(context.Background(), bson.M{"namespace": namespace}).Decode(&labSession)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("lab session not found")
		}
		return nil, err
	}
	return &labSession, nil
}

func (d *Database) UpdateLabSession(namespace string, labSession types.LabSession) error {
	labSession.UpdatedAt = time.Now()

	setFields := bson.M{
		"finish_result": labSession.FinishResult,
		"finished_at":   labSession.FinishedAt,
		"updated_at":    labSession.UpdatedAt,
	}
	labSessionsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("lab_sessions")
	_, err := labSessionsCollection.UpdateOne(context.Background(), bson.M{"namespace": namespace}, bson.M{"$set": setFields})
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetAllLabSessions(userID string) ([]types.LabSession, error) {
	labSessionsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("lab_sessions")
	var labSessions []types.LabSession
	cursor, err := labSessionsCollection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &labSessions)
	if err != nil {
		return nil, err
	}
	return labSessions, nil
}
