package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationSpec struct {
	AppName   string               `json:"app-name" bson:"app-name"`
	Image     string               `json:"image" bson:"image"`
	Volumes   []ApplicationVolume  `json:"volumes" bson:"volumes"`
	Ports     []ApplicationPortMap `json:"ports" bson:"ports"`
	Resources ApplicationResource  `json:"resources" bson:"resources"`
	Env       []ApplicationEnvVar  `json:"envFrom" bson:"envFrom"`
}
type ApplicationEnvVar struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type ApplicationResource struct {
	Storage int    `json:"storage" bson:"storage"`
	Memory  string `json:"memory" bson:"memory"`
	CPU     string `json:"cpu" bson:"cpu"`
}

type ApplicationVolume struct {
	VolumeName string `json:"volume-name" bson:"volume-name"`
	Path       string `json:"path" bson:"path"`
}

type ApplicationPortMap struct {
	Internal int    `json:"internal" bson:"internal"`
	External int    `json:"external" bson:"external"`
	Domain   string `json:"domain" bson:"domain"`
	SSL      bool   `json:"SSL" bson:"SSL"`
}

type Application struct {
	Name string          `json:"name" bson:"name"`
	Spec ApplicationSpec `json:"spec" bson:"spec"`
}

type ApplicationModel struct {
	Collection *mongo.Collection
}

func (m *ApplicationModel) Insert(app Application) (*mongo.InsertOneResult, error) {
	insertResult, err := m.Collection.InsertOne(context.Background(), app)
	if err != nil {
		return nil, fmt.Errorf("error inserting application: %v", err)
	}
	return insertResult, nil
}

func (m *ApplicationModel) FindOne(filter bson.M) (*Application, error) {
	var app Application
	err := m.Collection.FindOne(context.Background(), filter).Decode(&app)
	if err != nil {
		return nil, fmt.Errorf("error finding application: %v", err)
	}
	return &app, nil
}

func (m *ApplicationModel) Find(filter bson.M) ([]Application, error) {
	var apps []Application
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
