# faceit-user-service

A simple REST server providing CRUD operations on a User object. It stores the user records in a postgres database and publishes updates to User records to a Kafka topic so they can be consumed by any interested services.

## Running the service
The service can be run using its docker-compose file: `docker-compose up --build -d`

This will run:
- the service on `localhost:8080`
- a postgres container
- all necessary kafka containers

## Documentation
The documentation for the service is generated using swagger. Once the service has been run the documentation can be viewed at `http://localhost:8080/swagger/index.html#/` 

## Viewing the changelog
The messages published to kafka can be viewed using the kafka-ui at `http:localhost:9090`. They will be published to the `users-changelog` topic.

## Choices and assumptions
- Since the brief mentioned dockerised applications were preferred, I made the choice to make sure all technologies I used had to be started in the docker-compose file. This ruled out distributed options like MongoAB Atlas.
- As only one example of a user record was provided, I made the assumption that the fields for a user would be consistent for every user. This meant that using a relational database would be sufficient as I didn't need to store unstructured data. It would also be more extensible in the future if other tables relating to a user record needed to be added later on.
- I also made the choice to create a REST API instead of using gRPC as I didn't know what kind of applications would be accessing it, such as other microservices or a UI, so providing a REST API was the most flexible solution and provides the most compatibility compared to gRPC where the proto needs to be shared with the users of the API.

## Possible extensions and improvements
- One improvement that could be made to the service is I could use an ORM such as sqlc to query the database. This would make the service more maintainable as it would generate the code needed to query the database from the SQL queries you write.
- Obfuscating the users passwords when returning them in responses is another improvement I would make.
- For the changelog, I would also add additional fields showing the previous state of the user record and the new state of it.
- For the `GetUsers` endpoint, I would improve the filtering options by:
  - Allowing for users to be searched by `createdAt` or `updatedAt` fields, with less than, greater than, or date range options.
  - Provide ordering options, allowing the users to set the field to sort by and whether they are displayed in ascending or descending order.
- In the Dockerfile, using a scratch base image in the final stage for increased security.
