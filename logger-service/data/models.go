package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main(mongo *mongo.Client) Models{
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}

}

type Models struct{
	LogEntry LogEntry
}

type LogEntry struct{
	ID string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

}

func (l *LogEntry) Insert(entry LogEntry) error{
	collection:= client.Database("logs").Collection("logs")
	_, err:= collection.InsertOne(context.TODO(), LogEntry{
		ID: entry.ID,
		Name: entry.Name,
		Data: entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err!=nil{
		log.Print(err)
		return err
	}

	return nil
}
func (l *LogEntry) All()([]*LogEntry, error){
	ctx, cancel:= context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection:= client.Database("logs").Collection("logs")

	opts:=options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err:= collection.Find(context.TODO(), bson.D{}, opts)

	if err!=nil{
		log.Print("finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*LogEntry
	for cursor.Next(ctx){
		var logEntry LogEntry
		if err:= cursor.Decode(&logEntry); err!=nil{
			return nil, err
		}

		results = append(results, &logEntry)
	}

	return results, nil
}
func (l *LogEntry) GetOne(id string) (*LogEntry, error){
	ctx, cancel:= context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	collection:=client.Database("logs").Collection("logs")

	docID, err:=primitive.ObjectIDFromHex(id)
	if err!=nil{
		log.Print("converting id to object id error:", err)
	}

	var entry LogEntry
	err= collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)

	if err!= nil{
		return nil , err
	}
	return &entry, nil


}

func (l *LogEntry)  DropCollection() error{
	ctx, cancel:= context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection:= client.Database("logs").Collection("logs")
	if err:=collection.Drop(ctx); err!=nil{
		return err
	}
	return nil 
}
func (l *LogEntry) Update() (*mongo.UpdateResult, error){

	ctx, cancel:= context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection:= client.Database("logs").Collection("logs")



	docID, err:=primitive.ObjectIDFromHex(l.ID)
	if err!=nil{
		return nil, err
	}

	result, err:= collection.UpdateOne(
		ctx, 
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	return result, nil

}