package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SyncStandings connects to MongoDB and replaces the standings collection if the fetched data
// is newer than what is stored. Returns true if a sync was performed, false if the data was
// already up to date.
func SyncStandings(config Config, entries []VRSEntry) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoUri))
	if err != nil {
		return false, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(config.DatabaseName).Collection(config.CollectionName)

	needed, err := syncNeeded(ctx, collection, entries[0].StandingsDate)
	if err != nil {
		return false, err
	}
	if !needed {
		return false, nil
	}

	if err := deleteStandings(ctx, collection); err != nil {
		return false, err
	}
	if err := insertStandings(ctx, collection, entries); err != nil {
		return false, err
	}
	return true, nil
}

// syncNeeded checks whether the collection needs to be replaced. It groups documents by
// standings_date and applies the following rules:
//   - 0 groups: collection is empty, sync needed
//   - 1 group: compare stored date against fetchedDate, sync only if different
//   - 2+ groups: collection is in a inconsistent state, sync to recover
func syncNeeded(ctx context.Context, col *mongo.Collection, fetchedDate string) (bool, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$standings_date"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return false, err
	}

	switch len(results) {
	case 0:
		return true, nil
	case 1:
		return results[0].ID != fetchedDate, nil
	default:
		return true, nil
	}
}

// deleteStandings removes all documents from the collection so that it is ready for a fresh insert
func deleteStandings(ctx context.Context, col *mongo.Collection) error {
	_, err := col.DeleteMany(ctx, bson.D{})
	return err
}

// insertStandings inserts the full VRS into the database
func insertStandings(ctx context.Context, col *mongo.Collection, entries []VRSEntry) error {
	now := time.Now()
	docs := make([]any, len(entries))
	for i, e := range entries {
		e.SyncedAt = now
		docs[i] = e
	}
	_, err := col.InsertMany(ctx, docs)
	return err
}
