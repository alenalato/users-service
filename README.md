# Users Service

This is a gRPC server.\
The service specification, along with brief documentation, is available in the [proto folder](./proto) as Protocol Buffers definitions.

## API Design Considerations

### Endpoints

The server has reflection enabled, allowing clients to query the service's endpoint paths.

#### User Creation
`users.v1.Users/CreateUser`\
This operation takes user data as input and returns the created user as output.\
The password is hashed before being stored in the database. The password hash is never read from storage nor returned to the client.\
> Further password management should be handled by dedicated API operations.

#### User Edit
`users.v1.Users/UpdateUser`\
This operation takes a user ID and user data as input and returns the updated user as output.\
The input payload accepts only fields that can be updated.\
A [field mask](https://protobuf.dev/reference/protobuf/google.protobuf/#field-mask) is used to specify which fields to consider during the update.  

#### User Deletion
`users.v1.Users/DeleteUser`\
This operation takes a user ID as input and returns an empty response.

#### User Listing
`users.v1.Users/ListUsers`\
This operation takes a filter, page size, and a page token as input and returns a page of users matching the filter criteria.\
Filtering is allowed for several fields at a time. All fields are optional and must match the corresponding value when provided.\
Pagination is cursor-based: a _next page token_ is provided in every response. It allows the client to move to the subsequent page when passed in a request.
On the last page of a listing, it will be empty.\
Sorting of results is fixed to the user's creation timestamp.

#### Health Checking
`grpc.health.v1.Health/Check` or `grpc.health.v1.Health/Watch`\
The gRPC server exposes a health service that indicates if the server is healthy and serving.  

### Data Structures
#### User
A user contains its details provided on creation, except for the password value, plus creation and last update timestamps.\
Email and nickname are considered unique and are checked during both user creation and editing.

#### UserFilter
A user filter contains fields available for filtering a list of users.\
Email and nickname are not part of this structure because, being unique for users, they would always return a list with only one user.
> A more idiomatic API design would define a GetUser operation to retrieve a user by ID, email, or nickname.
 
#### UserEvent
A UserEvent payload contains user data, an event type specifying what happened, an event timestamp, and
a field mask to indicate which fields changed in case of an update. 

## Architectural Considerations

### Storage
The adopted storage engine is MongoDB:
- A NoSQL database is well-suited for this kind of write-heavy operations and can handle a high volume of requests with fast write performance.
- A user typically holds a lot of data related to themselves but not to other users (e.g., addresses, personal information). A document-based database allows storing all the data related to a user in a single document, making it easier to retrieve all the data in a single query without joins.
- The documental structure allows easier extension of the user data model.

### Event Emitter
The adopted data bus is Kafka, which is well-suited for large-scale event streaming.\
The event emission implementation is basic: storage changes and event emission are not atomic operations, so event emission errors cannot be blocking.
> A better event emission implementation would use a [transactional outbox](https://microservices.io/patterns/data/transactional-outbox.html).

### High Throughput Traffic

- The gRPC server implementation, along with protobuf, is well-suited to scale thanks to goroutines and faster serialization times. The server application may need to scale vertically on resources or horizontally on infrastructure.  
- Password hashing is slow and could become a bottleneck during user creation: it could be moved to an asynchronous operation.
- Event emission is currently synchronous: it could be moved to an asynchronous operation.
- Even if the application handles traffic load, MongoDB may need to scale (vertically or horizontally) due to connection pooling limits or high CPU load from increased I/O.

## Codebase

### Project Structure
The application is mainly structured into the following directories:
- `cmd` contains the executable main packages.
  - `cmd/server` contains the main package that starts the gRPC server.
- `proto` contains the Protocol Buffers definitions of the Users service.
- `pkg/grpc` contains the gRPC server definition.
  - This code is **generated** by the [protoc generator](./scripts/generate-server-code.sh) from the Protocol Buffers files.
- `internal` contains the packages used internally by the application.
  - `grpc` contains the implementation of the gRPC handlers.
  - `businesslogic` contains the business logic abstraction for the application.
    - `user` contains the user management logic.
    - `password` contains the password manager implementation.
  - `events` contains the event emission handler.
  - `storage` contains the storage repository of the application.
    - `mongodb` contains the MongoDB storage implementation.
  - `logger` and `common` contain various utilities.

### Implementation Considerations
- The application's packages are designed to be decoupled using interfaces and dependency injection.
- Business logic is the center of the design, while gRPC input, event emitter output, and storage repository are on the edges,
inspired by [hexagonal architecture](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749).
- Future feature requests will adhere to abstractions but can also leverage the isolation of each package's responsibility to better understand the impact of new code and implement safer changes.
- How to expand the solution:
  - Add a GetUser endpoint to retrieve a user by unique values like ID, nickname, or email.
  - Add API authentication/authorization.
  - Refine health checks with dependency checks, e.g., pinging the database.
  - Refine validation of input formats, such as the email field or the UUID format for the ID.
  - Refine validation error reporting to clients by adopting structured protobuf error details or introducing dedicated response fields.
  - Allow clients to specify the sorting of ListUsers results.
  - Allow clients to specify multiple values for a field in the ListUsers filter.
  - Validate the next page token against the given filter for ListUsers after page 1.
  - Adopt a transactional outbox for event emission.
  - Provide deleted user data in the DeleteUser response and related UserEvent.
  - Add password management operations, e.g., password updates.

### Dependencies
Required tools are `docker` with the `docker compose plugin` and, optionally, a bash-compatible shell.

### Testing
The application is tested with **unit tests** for the business logic and the gRPC server using [generated mocked dependencies](https://github.com/uber-go/mock).\
The coverage of unit tests is intentionally not 100% because some parts were not worth testing for the scope of this development.

The MongoDB storage is tested with **integration tests** that use a [disposable container with a real MongoDB instance](https://github.com/ory/dockertest). These tests are slower than unit tests but are more realistic and can catch integration issues with the real MongoDB instance.

## Usage

**Note**: All the following commands are intended to be run from the root of the project.

Set up the local environment:
```bash
cp .env.example .env
```

Configurations can be customized. These are the default values for the provided [docker-compose.yaml](docker-compose.yaml) file:
```
# level for logging, info is default if not set
LOG_LEVEL=debug

# gRPC server configuration
GRPC_LISTEN_HOST=0.0.0.0
GRPC_LISTEN_PORT=9090

# MongoDB configuration
MONGODB_URI=mongodb://mongo:27017
MONGODB_DATABASE=users

# Kafka configuration
KAFKA_ADDRESSES=kafka:9092
KAFKA_EVENT_EMITTER_TOPIC_NAME=users
```

The Docker Compose project is available in the [docker-compose.yaml](docker-compose.yaml) file.

#### Build the Application
This will build the application executables:
```bash
docker compose run --rm server go build -o .build/server ./cmd/server/main.go
```

#### Test the Application
This will run unit and integration tests:
```bash
docker compose run --rm server go test ./...
```

#### Start Dependency Instances
```bash
docker compose --profile dependencies up -d
```

A Kafka UI pointing to the local Kafka instance will be available after a few minutes at http://localhost:9093.

#### Start the Application
This will start the server application.\
**Note**: The MongoDB instance must be running and reachable at `MONGODB_URI`, and the Kafka instance must be running and reachable at `KAFKA_ADDRESSES`.
```bash
docker compose run --rm --service-ports server .build/server
```

The server will be available on `localhost:9090` with the default configuration.

### Other Commands

#### Generate the pkg/grpc package

```bash
./scripts/generate-server-code.sh
```

#### Run the Application for Development

```bash
docker compose run --rm --service-ports server go run ./cmd/server/main.go
```
