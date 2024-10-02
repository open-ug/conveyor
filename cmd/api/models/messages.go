package models

import (
	"context"
	"fmt"

	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DriverMessageModel struct {
	Collection *mongo.Collection
}

func (m *DriverMessageModel) Insert(
	message craneTypes.DriverMessage) (*mongo.InsertOneResult, error) {
	// check if the message already exists

	insertResult, err := m.Collection.InsertOne(context.Background(), message)
	if err != nil {
		return nil, fmt.Errorf("error inserting message: %v", err)
	}
	return insertResult, nil
}

func (m *DriverMessageModel) FindOne(filter bson.M) (*craneTypes.DriverMessage, error) {
	var message craneTypes.DriverMessage
	err := m.Collection.FindOne(context.Background(), filter).Decode(&message)
	if err != nil {
		return nil, fmt.Errorf("error finding message: %v", err)
	}
	return &message, nil
}

func (m *DriverMessageModel) Find(filter bson.M) ([]craneTypes.DriverMessage, error) {
	var messages []craneTypes.DriverMessage
	cursor, err := m.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error finding messages: %v", err)
	}
	if err = cursor.All(context.Background(), &messages); err != nil {
		return nil, fmt.Errorf("error finding messages: %v", err)
	}
	return messages, nil
}

func (m *DriverMessageModel) UpdateOne(filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	fmt.Println(update)
	updateResult, err := m.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating message: %v", err)
	}
	return updateResult, nil
}

func (m *DriverMessageModel) DeleteOne(filter bson.M) (*mongo.DeleteResult, error) {
	deleteResult, err := m.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error deleting message: %v", err)
	}
	return deleteResult, nil
}
