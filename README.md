# UserManagment service based on gRPC and go

## Architecture

![image](https://user-images.githubusercontent.com/34543261/188322736-b154d6d3-430e-4e1f-8954-2666c651877a.png)

This service implements a grpc server which publishes an API to perform user management actions against a database.
In the docker-compose, the above arquitecture is deployed, whhich consists of the server, the database and a client consuming a server-side streaming which simulates another service receiving notifications.

When a user action is performed, a notification is sent through the stream and received by the client, which logs a message containing a brief description of the action performed.

## Code structure
The service code is organized in a way that the implemented grpc server contains a database client. This client implements an Interface: AdapterInterface from the database package.
Thanks to that, if the database used to store the users must change for whatever reason, only the database client code must be modified, maintaining the rest of the service unaltered.

## Proto files generation
From /proto directory execute the following commands:
```
>> protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  .\userManagement.proto //generates server/client code
>> protoc -I. --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true .\*.proto // generates api code
>> protoc -I.  --openapiv2_out . --openapiv2_opt logtostderr=true .\*.proto // generates swagger JSON

```

## Deploying the environment

```
>> docker-compose up
```

![image](https://user-images.githubusercontent.com/34543261/188349517-5c314e2e-dc07-4fbe-92fd-33a3665bbd04.png)

## Swagger
To complement the deployment and have a better understanding of how the different API endpoints from the service work, there is a swagger UI available on localhost:8081 after deploying the service.

![image](https://user-images.githubusercontent.com/34543261/188349748-0bdba5c5-f46f-4b9e-99cd-084f1578a5fd.png)

There is also an exported file from postman, that also grants useful information about the endpoints.

## Calling the API endpoints
There are a couple of things that should be explained before start using the service.
When calling the endpoints and filtering to find a specific user, like in an update operation or in a get user request, this search can be done both by database user_id or by email:

![image](https://user-images.githubusercontent.com/34543261/188350112-d04cac53-011f-4087-94df-3e4469fe05d7.png)

![image](https://user-images.githubusercontent.com/34543261/188350172-69f63558-07d7-43ad-b77f-f40fb6cf0d87.png)

In the above images it can be observed that both email and id returns the same.

### List users endpoint
When calling ListUsers endpoint, different filters can be applied. These filters are the some of the fields which conform the user, in other words:

- firstName 
- lastName
- email
- country

All of them can be received as query params:

![image](https://user-images.githubusercontent.com/34543261/188350628-fe960efd-b087-4f6e-a893-37889fa50df3.png)

As it can be observed, the field to be used as query param must be preceded by "filter.". So, for example, in case a filter by country is wanted to be applied, the next query param should be added to the request "?filter.country=UK". This filtering is yet to be added to the swagger UI, but an example can be found within the postman exported collection.

## User actions notifications
As it can be seen in the diagram at the beginning of this Readme, there is a notifications receiver which logs a brief description of the different actions that are performed when calling the API:

![image](https://user-images.githubusercontent.com/34543261/188351029-0e4b8105-c8ee-4bc8-8d02-c876787746bd.png)

This feature should be improved to be able to stream the notifications to more than one client, or maybe make this notification receiver kind of a middleware that receives notifications, process them, and later on sends them where it is needed.

## About the tests
Inside the tests folder two files can be found. One for the grpc server and client methods and the other for mongodb client operations. The first file's tests are prepared to be run in any environment due to the fact that all the external needed resources are mocked. On the other hand, the mongo client tests require of a mongodb instance running on port 27017, which can be easily accomplished using docker:

```
>> docker run -d -p 27017:27017 --name test-mongo mongo:latest
```





