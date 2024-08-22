/*
Copyright © 2024 Cranom Technologies Limited info@cranom.tech
*/
package models

import (
	"context"
	"fmt"

	craneTypes "crane.cloud.cranom.tech/cmd/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationModel struct {
	Collection *mongo.Collection
}

func (m *ApplicationModel) Insert(app craneTypes.Application) (*mongo.InsertOneResult, error) {
	// check if the application already exists
	_, err := m.FindOne(bson.M{"name": app.Name})
	if err == nil {
		return nil, fmt.Errorf("application with name %s already exists", app.Name)
	}

	insertResult, err := m.Collection.InsertOne(context.Background(), app)
	if err != nil {
		return nil, fmt.Errorf("error inserting application: %v", err)
	}
	return insertResult, nil
}

func (m *ApplicationModel) FindOne(filter bson.M) (*craneTypes.Application, error) {
	var app craneTypes.Application
	err := m.Collection.FindOne(context.Background(), filter).Decode(&app)
	if err != nil {
		return nil, fmt.Errorf("error finding application: %v", err)
	}
	return &app, nil
}

func (m *ApplicationModel) Find(filter bson.M) ([]craneTypes.Application, error) {
	var apps []craneTypes.Application
	cursor, err := m.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error finding applications: %v", err)
	}
	if err = cursor.All(context.Background(), &apps); err != nil {
		return nil, fmt.Errorf("error finding applications: %v", err)
	}
	return apps, nil
}

func (m *ApplicationModel) UpdateOne(filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	fmt.Println(update)
	updateResult, err := m.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating application: %v", err)
	}
	return updateResult, nil
}

func (m *ApplicationModel) DeleteOne(filter bson.M) (*mongo.DeleteResult, error) {
	deleteResult, err := m.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error deleting application: %v", err)
	}
	return deleteResult, nil
}

func NewApplicationModel(db *mongo.Database) *ApplicationModel {
	return &ApplicationModel{
		Collection: db.Collection("applications"),
	}
}
