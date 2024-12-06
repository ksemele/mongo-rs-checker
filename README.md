# mongo-rs-checker

The mongo-rs-checker is a Go app designed to interact with a MongoDB replica set. It performs several key operations:
* Database and Collection Generation: It automatically generates a random database and collection name if not specified.
* Document Insertion: It inserts a sample document into the generated collection.
* Document Retrieval: It retrieves the inserted document by its ID.
* Collection Querying: It queries all documents in the collection.
* Collection Dropping: It drops the collection after operations are completed.

## Environment Variables

The mongo-rs-checker expects the following environment variables:
* `MONGODB_URI`: The URI to connect to the MongoDB instance. It should include credentials and replica set information if applicable.
* `MONGODB_DATABASE_NAME`: (Optional) The name of the database to use. If not set, a random database name will be generated.

## How to Run

### Running Locally

To run the mongo-rs-checker locally using Docker, set the necessary environment variables and execute:

```shell
docker run -e MONGODB_URI="mongodb://<USER>:<PASSWORD>@localhost:27017" -e MONGODB_DATABASE_NAME="your_database_name" quay.io/greengrunt/checkers/mongo-rs-checker:0.0.1
```

### Deploying on Kubernetes

To deploy the mongo-rs-checker on Kubernetes, you can use the following Pod specification as an example in [./deploy](./deploy/mongo-pod-checker.yaml)

## Key Points
* The mongo-rs-checker will generate a random database name if `MONGODB_DATABASE_NAME` is not provided.
* It generates a new collection name for each run to avoid conflicts.
* After performing its operations, the mongo-rs-checker will drop the collection to clean up.

This setup is ideal for testing and development purposes where you need to interact with a MongoDB replica set to ensure connectivity and basic operations are functioning correctly.

## Additional info

Example of successfully ran Pod

```bash
2024/12/06 09:06:20 connectToMongoDB() successful
2024/12/06 09:06:20 Using database: test
2024/12/06 09:06:20 Using collection: mongo_rs_checker_26aecb2d
Inserted document with ID: ObjectID("6752be8cc2a84612baacfcc4")
2024/12/06 09:06:20 Found document: map[_id:ObjectID("6752be8cc2a84612baacfcc4") age:30 city:NY name:Bi Ba]
[{_id ObjectID("6752be8cc2a84612baacfcc4")} {name Bi Ba} {age 30} {city NY}]
2024/12/06 09:06:20 Collection mongo_rs_checker_26aecb2d dropped successfully
stream closed EOF for default/mongo-rs-checker (mongo-rs-checker)
```
