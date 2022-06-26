package main

import (
	"context"
	"fmt"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context,
		cancel context.CancelFunc){
			
	// CancelFunc to cancel to context
	defer cancel()
	
	// client provides a method to close
	// a mongoDB connection.
	defer func(){
	
		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil{
			panic(err)
		}
	}()
}

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.

func connect(uri string)(*mongo.Client, context.Context,
						context.CancelFunc, error) {
						
	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
									30 * time.Second)
	
	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func ping(client *mongo.Client, ctx context.Context) error{

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

// insertOne is a user defined method, used to insert
// documents into collection returns result of InsertOne
// and error if any.
func insertOne (client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{})(*mongo.InsertOneResult, error) {
 
    // select database and collection ith Client.Database method
    // and Database.Collection method
    collection := client.Database(dataBase).Collection(col)
     
    // InsertOne accept two argument of type Context
    // and of empty interface  
    result, err := collection.InsertOne(ctx, doc)
    return result, err
}

// query is user defined method used to query MongoDB,
// that accepts mongo.client,context, database name,
// collection name, a query and field.
 
//  database name and collection name is of type
// string. query is of type interface.
// field is of type interface, which limits
// the field being returned.
 
// query method returns a cursor and error.
func query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	 
		// select database and collection.
		collection := client.Database(dataBase).Collection(col)
		 
		// collection has an method Find,
		// that returns a mongo.cursor
		// based on query and field.
		result, err = collection.Find(ctx, query,
									  options.Find().SetProjection(field))
		return
	}

// insertMany is a user defined method, used to insert
// documents into collection returns result of
// InsertMany and error if any.
func insertMany (client *mongo.Client, ctx context.Context, dataBase, col string, docs []interface{}) (*mongo.InsertManyResult, error) {
 
    // select database and collection ith Client.Database
    // method and Database.Collection method
    collection := client.Database(dataBase).Collection(col)
     
    // InsertMany accept two argument of type Context
    // and of empty interface  
    result, err := collection.InsertMany(ctx, docs)
    return result, err
}

func main(){

	// Get Client, Context, CancelFunc and
	// err from connect method.
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil{
		panic(err)
	}
	
	// Release resource when the main
	// function is returned.
	defer close(client, ctx, cancel)
	
	// Ping mongoDB with Ping method
	ping(client, ctx)

	// Create  a object of type interface to  store
    // the bson values, that  we are inserting into database.
    var document interface{}
     
     
    document = bson.D{
        {"rollNo", 175},
        {"maths", 80},
        {"science", 90},
        {"computer", 95},
		{"me","Will be awesome at Dell, Look, I'm even learning MongoDB and GO!"},
	}
    // insertOne accepts client , context, database
    // name collection name and an interface that
    // will be inserted into the  collection.
    // insertOne returns an error and a result of
    // insert in a single document into the collection.
    insertOneResult, err := insertOne(client, ctx, "gfg",
                                      "marks", document)

	// handle the error
    if err != nil {
        panic(err)
    }

	/*
            List databases
    */
    databases, err := client.ListDatabaseNames(ctx, bson.M{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Databases:", databases)
     
    // print the insertion id of the document,
    // if it is inserted.
    fmt.Println("Result of InsertOne")
    fmt.Println(insertOneResult.InsertedID)

	// Now will be inserting multiple documents into
    // the collection. create  a object of type slice
    // of interface to store multiple  documents
    var documents []interface{}
     
    // Storing into interface list.
    documents = []interface{}{
        bson.D{
            {"rollNo", 153},
            {"maths", 65},
            {"science", 59},
            {"computer", 55},
        },
        bson.D{
            {"rollNo", 162},
            {"maths", 86},
            {"science", 80},
            {"computer", 69},
        },
    }
 
 
    // insertMany insert a list of documents into
    // the collection. insertMany accepts client,
    // context, database name collection name
    // and slice of interface. returns error
    // if any and result of multi document insertion.
    insertManyResult, err := insertMany(client, ctx, "gfg",
                                        "marks", documents)
	// handle the error
    if err != nil {
        panic(err)
    }
 
    fmt.Println("Result of InsertMany")
     
    // print the insertion ids of the multiple
    // documents, if they are inserted.
    for id := range insertManyResult.InsertedIDs {
        fmt.Println(id)
    }

	  // create a filter an option of type interface,
    // that stores bjson objects.
    var filter, option interface{}
     
    // filter  gets all document,
    // with maths field greater that 70
    filter = bson.D{
        {"maths", bson.D{{"$gt", 70}}},
    }
     
    //  option remove id field from all documents
    option = bson.D{{"_id", 0}}

	// call the query method with client, context,
    // database name, collection  name, filter and option
    // This method returns momngo.cursor and error if any.
    cursor, err := query(client, ctx, "gfg",
                         "marks", filter, option)
    // handle the errors.
    if err != nil {
        panic(err)
    }
 
    var results []bson.D
     
    // to get bson object  from cursor,
    // returns error if any.
    if err := cursor.All(ctx, &results); err != nil {
     
        // handle the error
        panic(err)
    }
     
    // printing the result of query.
    fmt.Println("Query Result")
    for _, doc := range results {
        fmt.Println(doc)
    }
}
