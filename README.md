# Overview
The application is built to run in Docker with a local stack setup.
The local stack consists of the image for the app itself and an image for the DB, in this case we are using `MySQL`
### How to run the application

To build, run, stop and delete the application from Docker, use the following `Makefile` commands:
- `make build` - builds the local stack and seeds the DB
- `make up` - runs the local stack in Docker
- `make down` - stops the application
- `make nuke` - deletes the application from docker (will require a make build again)

## Architecture Overview

### Endpoints
The application consists of 4 endpoints
- POST http://localhost:8080/login
- POST http://localhost:8080/user/create
- GET  http://localhost:8080/discover
	- sort by distance http://localhost:8080/user/discover?sort_by=distance
	- sort by age http://localhost:8080/user/discover?sort_by=age
    - sort by age http://localhost:8080/user/discover?sort_by=gender
	- sort by rank (attractiveness) http://localhost:8080/user/discover?sort_by=rank
		- To increase the rank of someone one, use the same `target_user_id`, but change the `user_id` when calling `/swipe?...`
- POST http://localhost:8080/swipe?user_id=1&target_user_id=2&match=YES (or NO)
### Authentication and Authorisation

A middleware is implemented to check if a Bearer token has been passed and it also checks the integrity of the token with a Public key. That way we can determine if the is valid or tampered with and allow/block further access.

The Bearer token is asymmetrically encrypted JWT, meaning we sign the JWT with a private, securely stored key on our end and we can check the integrity of the key in our middleware once it is passed from upstream callers. This ensures the token is valid and not tampered.

#### Trade offs with asymmetric encryption vs symmetric
Asymmetric encryption is a slower operation but more secure.
Symmetric encryption is faster but less secure.

For this application I have used asymmetric encryption to gain the security benefit.

### Database
For this implementation, I am using MySQL to store the data in a relational way.
Since we have users, interactions and potential matches, a relational database was an easy choice for this current implementation.

### Functional and non-functional requirements considered
- data validation from the request query parameters could be better
- currently there is one DB instance - could split the DB int a Leader and Follower
	- write only into the leader
	- read only from the follower
- requested data can be cached with Redis or other caching mechanisms, few strategies considered:
	- write through - add data into the cache first, then store data in the DB
	- write around - write data into the DB first, then update the data
	- write back - write the data into the cache, and at some point later update the db
	- caching strategy may differ for each endpoint
- queues and message brokers
	- if we need to serialise and process a large number of operations, we can queue them
	- message brokers can be used to handle the communication between different components of the app, such as notification, active feeds, scalability and load balancing
- gRPC - some the functionality can be decoupled to be server-to-server only, which will gRPC a good candidate if we are looking for performance gains
- API Gateway - currently each individual endpoint used a middleware for authorisation. We could implement this on a API Gateway level with an authorizer function that checks all incoming requests, this will simplify our endpoints
