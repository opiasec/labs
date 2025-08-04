package database

import (
	"appseclabs/types"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *Database) GetAllActiveLabsDefinitions() ([]types.Lab, error) {
	var labs []types.Lab

	labsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("labs")

	projection := bson.M{"name": 1,
		"slug":       1,
		"labSpec":    1,
		"created_at": 1,
		"updated_at": 1,
	}

	cursor, err := labsCollection.Find(context.Background(), bson.M{"status": "active"}, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &labs)
	if err != nil {
		return nil, err
	}
	return labs, nil
}

func (d *Database) GetLabDefinitionBySlug(slug string) (*types.Lab, error) {
	var lab types.Lab

	labsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("labs")

	err := labsCollection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&lab)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("lab definition not found")
		}
		return nil, err
	}
	return &lab, nil
}

func (d *Database) UpdateLabDefinition(slug string, lab *types.Lab) error {
	lab.UpdatedAt = time.Now()

	labsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("labs")
	_, err := labsCollection.UpdateOne(context.Background(), bson.M{"slug": slug}, bson.M{"$set": lab})
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) CreateLabDefinition(lab *types.Lab) error {
	lab.CreatedAt = time.Now()
	lab.UpdatedAt = time.Now()

	labsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("labs")
	_, err := labsCollection.InsertOne(context.Background(), lab)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) DeleteLabDefinition(slug string) error {
	labsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("labs")
	_, err := labsCollection.DeleteOne(context.Background(), bson.M{"slug": slug})
	if err != nil {
		return err
	}
	return nil
}
