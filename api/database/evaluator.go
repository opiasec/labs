package database

import (
	"appseclabs/types"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *Database) GetEvaluationBySlug(slug string) (*types.Evaluation, error) {
	var evaluation types.Evaluation

	evaluationsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("evaluations")

	err := evaluationsCollection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&evaluation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("evaluation not found")
		}
		return nil, err
	}
	return &evaluation, nil
}

func (d *Database) GetEvaluators() ([]types.Evaluation, error) {
	evaluatorsCollection := d.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("evaluations")
	var evaluators []types.Evaluation

	cursor, err := evaluatorsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &evaluators)
	if err != nil {
		return nil, err
	}
	return evaluators, nil
}
