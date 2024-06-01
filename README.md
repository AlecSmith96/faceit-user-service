# faceit-user-service

A simple REST server providing CRUD operations on a User object. It stores the user records in a postgres database and publishes updates to User records to a Kafka topic so they can be consumed by any interested services.

## Running the service
The service can be run using its docker-compose file: `docker-compose up --build -d` 
- Note: please use the `--build` command when running as otherwise there can be kafka issues. 

This will run:
- the service on `localhost:8080`
- a postgres container
- all necessary kafka containers

## Documentation
The documentation for the service is generated using swagger. Once the service has been run the documentation can be viewed at `http://localhost:8080/swagger/index.html#/` 

## Running the tests
The tests can be run using the following make command `make test`.

## Viewing the changelog
The messages published to kafka can be viewed using the kafka-ui at `http://localhost:9090`. They will be published to the `users-changelog` topic.

## Choices and assumptions
- I chose to implement the service using Clean Architecture as it is a design principle that aims to make code more readable and maintainable. It decouples the services business logic from its application code by separating code into layers, making it easier to tell what the service does rather than what it's built with. The four layers are:
  - `drivers`: This layer is for specific framework or application code, the only code in this layer is the gin router.
  - `usecases`: This layer holds the main business logic for each endpoint, separated into individual files for readability.
  - `adapters`: This layer contains the interfaces needed for the application, in this case it contains the code to interact with postgres and kafka.
  - `entities`: This layer has all of the internal structs for the service.
- I also made the choice that the communication of user updates to interested services should be asynchronous. This is because there could be multiple services interested in receiving updates, so writing to a message queue would allow multiple services to consume the updates at the same time. Using a message queue also means that the endpoints wouldn't have to wait for the services to consume the update before returning the response. 
- Since the brief mentioned dockerised applications were preferred, I made the choice to make sure all technologies I used had to be started in the docker-compose file. 
- i also made the assumption that passwords didn't need to be stored securely as it wasn't mentioned in the brief. This is something I would do to improve the service if I had more time.
- As only one example of a user record was provided, I made the assumption that the fields for a user would be consistent for every user. This meant that using a relational database would be sufficient as I didn't need to store unstructured data. It would also be more extensible in the future if other tables relating to a user record needed to be added later on.
- I also made the choice to create a REST API instead of using gRPC as I didn't know what kind of applications would be accessing it, such as other microservices or a UI, so providing a REST API was the most flexible solution and provides the most compatibility compared to gRPC where the proto needs to be shared with the users of the API.

## Possible extensions and improvements
- One improvement that could be made to the service is I could use an ORM such as sqlc to query the database. This would make the service more maintainable as it would generate the code needed to query the database from the SQL queries you write.
- Obfuscating the users passwords when returning them in responses is another improvement I would make.
- For the changelog, I would also add additional fields showing the previous state of the user record and the new state of it.
- For the `GetUsers` endpoint, I would improve the filtering options by:
  - Allowing for users to be searched by `createdAt` or `updatedAt` fields, with less than, greater than, or date range options.
  - Provide ordering options, allowing the users to set the field to sort by and whether they are displayed in ascending or descending order.
  - Make sure that the final page of results for a query doesn't include a pageToken to an empty page.
- In the Dockerfile, using a scratch base image in the final stage for increased security.
- Implement tracing at the usecase and adapter layers to identify any potential performance optimisations.
