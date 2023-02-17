package controllers

import (
	"context"
	"time"

	"github.com/kiprop-dave/2faAuth/pkg/config"
	"github.com/kiprop-dave/2faAuth/pkg/models"
)

var attemptsCollection = config.GetCollection(config.DB, "attempts")

func LogAttempt(name, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	attempt := models.Attempt{Name: name, UserId: id, Time: time.Now()}

	_, err := attemptsCollection.InsertOne(ctx, attempt)
	return err
}
