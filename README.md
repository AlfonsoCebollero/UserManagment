# UserManagment service based on gRPC and go

## Architecture

![image](https://user-images.githubusercontent.com/34543261/188322736-b154d6d3-430e-4e1f-8954-2666c651877a.png)

This service implements a grpc server which publishes an API (consult endpoints in postman file within this repo) to perform user actions against a database.
In the docker-compose, the above arquitecture is deployed, whhich consists of the server, the database and a client consuming a server-side streaming which simulates another
service receiving notifications.

When a user action is performed, a notification is sent through the stream and receive by the client, which logs a message containing a brief description of the action performed.


## Code organization
The service code is organized in a way that the implmented grpc server has within it a database client. This client implements an Interface: AdapterInterface whitin the database package.
Thanks to that, if the database used for the prject must change for whatever reason, only the database client code must be modified, maintaining the rest of the code unaltered.

## Proto files generation
From /proto directory execute the following commands:
```
>> protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  .\userManagement.proto //generates server/client code
>> protoc -I. --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true .\*.proto // generates api code

```

# Deploying the environment
