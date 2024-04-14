# Overview

The application is built to run in Docker with a local stack setup.
The local stack consists of the image for the app itself and an image for the DB, in this case we are using `MySQL`.\
Few tables and example data are seeded when the application is built for the first time.

Secrets and passwords have no meaning outside of this excercise so they are hardcoded or stored in local files. Supporting comments are provided in the code.

## Get Started

To build, run, stop and delete the application from Docker, use the following `Makefile` commands:

- `make build` - builds the local stack and seeds the DB with three users and login credentials
- `make up` - runs the local stack in Docker
- `make down` - stops the application
- `make nuke` - deletes the application from docker (will require a `make build` again)

## Seeded Users & Credentials

Available users ready to be used to test different flows of the applications

```json
{
    "email": "billieeilish@example.com",
    "password": "password-user-two"
}

{
    "email": "jonwick@example.com",
    "password": "password-user-one"
}

{
    "email": "tonystark@example.com",
    "password": "password-user-two"
}
```

## Architecture Overview

### Endpoints

The application consists of 4 endpoints. \
Each endpoints apart from `/user/login` is protected by an auth middleware

- POST http://localhost:8080/user/login
- POST http://localhost:8080/user/create
- GET http://localhost:8080/user/discover
  - sort by distance http://localhost:8080/user/discover?sort_by=distance
  - sort by age http://localhost:8080/user/discover?sort_by=age
    - sort by gender http://localhost:8080/user/discover?sort_by=gender
  - sort by rank (attractiveness) http://localhost:8080/user/discover?sort_by=rank
    - To increase the rank of someone one, use the same `target_user_id`, but change the `user_id` when calling `/swipe?user_id=1&target_user_id=2&match=YES`
  - To decrease the rank of a given user pass `&match=NO`
- POST http://localhost:8080/user/swipe?user_id=1&target_user_id=2&match=YES (or NO)

### Authentication and Authorisation

A middleware is implemented to check if a Bearer token has been passed and it also checks the integrity of the token with a Public key. That way we can determine if the is valid or tampered with and allow/block further access.

The Bearer token is asymmetrically encrypted JWT, meaning we sign the JWT with a private, securely stored key on our end and we can check the integrity of the key in our middleware once it is passed from upstream callers. This ensures the token is valid and not tampered.

### Encryption Tradeoffs

Asymmetric encryption is a slower operation but more secure.\
Symmetric encryption is faster but less secure.

For this application I have used **asymmetric** encryption to simulate a stronger secuity pattern.

### Database

For this implementation I am using MySQL as a primary data store.\
Since we have users, interactions and potential matches, a relational database was an easy choice for this current implementation.

## Areas for Improvement

- Query param validation
- Currently there is one DB instance, Leader & Follower could be a good choice
  - **write** to the leader
  - **read** from the follower
- Caching and caching strategies:
  - **Write Through** - add data into the cache first, then store data in the DB
  - **Write Around** - write data into the DB first, then update the data
  - **Write Back** - write the data i nto the cache, and at some point later update the db
  - caching strategy may differ for each endpoint
- Queues and message brokers
  - if we need to serialise and process a large number of operations, we can queue them
  - message brokers can be used to handle the communication between different components of the app, such as notification, active feeds, scalability and load balancing
- gRPC - some the functionality can be decoupled to be server-to-server only, which will gRPC a good candidate if we are looking for performance gains
- API Gateway and single authoriser function behind all routes
- Implement routing, example `github.com/gorilla/mux` to cleanup the code and make it easier to add new routes or subroutes
- Further edge cases considerations
- Comprehensive set of tests
