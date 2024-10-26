# Data Services

A high-performance, distributed data access layer implementing request coalescing and hash-based routing to reduce database load and prevent hot partitions. 

## Overview

Data Services is a middleware layer that sits between API servers and Cassandra clusters, providing request coalescing. It's designed to handle high-traffic scenarios efficiently by reducing duplicate database queries and preventing database overload.
It is inspired by Discord's architecture explained in their blog post: [HOW DISCORD STORES TRILLIONS OF MESSAGES](https://discord.com/blog/how-discord-stores-trillions-of-messages).

An example of usecase from Discord is when a big announcement is sent on a large server (Discord group) that notifies @everyone: users are going to open the app and read the same message, sending tons of traffic to the database. This is where request coalescing comes in handy, as it can combine all the requests for the same data into a single database query, reducing the load on the database and preventing hot partitions.

A simpler way of understanding it is: caching with the duration equal to the time spent running the query. No client has to be aware of the coalescing because the max amount of staleness is the same as if each client had run the query themselves. It also doesn't require extra memory, because the query result falls out of scope as soon as it is sent to all waiters.

### Key Features

- **Request Coalescing**: Automatically combines duplicate requests for the same data into a single database query
- **Consistent Hash-based Routing**: Routes related requests to the same service instance for optimal coalescing
- **Distributed Architecture**: Multiple service instances working in parallel
- **High Availability**: Data service nodes are stateless and can be scaled horizontally
- **Monitoring**: Built-in metrics for tracking requests and queries counts

## Setup

```bash
$ docker-compose up --build
```
Wait for the services to start up. The Cassandra cluster will be initialized with the required keyspace. You will see something like this in the logs that shows that the data service instances are ready to accept requests:
```
data-service1-1   | 2024/10/26 16:18:45 Connected to cassandra
data-service1-1   | 2024/10/26 16:18:45 Starting server on port 50051
data-service2-1   | 2024/10/26 16:18:45 Connected to cassandra
data-service2-1   | 2024/10/26 16:18:45 Starting server on port 50052
```

Run the client CLI to send test requests to the data service:
```bash
$ go run ./client -h
  -channels int
        Number of unique channels to distribute requests across (number of unique requests) (default 20)
  -requests int
        Total number of requests to send (default 10000)
```


Example usage:
```
$ go run ./client
2024/10/26 17:08:11 Unique requests: 20, Total requests: 10000, Total queries executed: 184
2024/10/26 17:08:11 Average queries per request: 0.0184
2024/10/26 17:08:11 Saved queries by coalescing: 9816
2024/10/26 17:08:11 Total time taken: 816.727364ms
```

## Architecture

### Components

1. **Data Service Nodes**: gRPC servers that handle incoming requests and manage database connections
2. **Cassandra Cluster**: A 3-node Cassandra cluster for data storage
3. **Client**: gRPC Test client for simulating high-traffic scenarios

```mermaid
%%{init: {'theme': 'base', 'themeVariables': { 'fontSize': '16px', 'fontFamily': 'arial', 'nodeTextSize': '16px', 'labelTextSize': '16px', 'titleTextSize': '20px' }}}%%

flowchart TB
    subgraph Users ["Multiple Users"]
        U1[User 1]
        U2[User 2]
        U3[User 3]
    end

    subgraph API ["Hash-based Routing API"]
        A1[API Server 1]
        A2[API Server 2]
    end

    subgraph DS ["Data Services Layer"]
        direction TB
        subgraph DS1 ["Data Service Instance 1"]
            direction TB
            subgraph Coalescing1 ["Request Coalescing"]
                R1["Request (Channel 1)"]
                R2["Request (Channel 1)"]
                R3["Request (Channel 1)"]
                RC1[Request Coalescer]
                DQ1["Single DB Query"]
                R1 & R2 & R3 --> RC1
                RC1 --> DQ1
            end
        end
        subgraph DS2 ["Data Service Instance 2"]
            direction TB
            subgraph Coalescing2 ["Request Coalescing"]
                R4["Request (Channel 2)"]
                R5["Request (Channel 2)"]
                R6["Request (Channel 2)"]
                RC2[Request Coalescer]
                DQ2["Single DB Query"]
                R4 & R5 & R6 --> RC2
                RC2 --> DQ2
            end
        end
    end

    subgraph DB ["Cassandra Cluster"]
        C1[Node 1] <--> C2[Node 2] <--> C3[Node 3] <--> C1
    end

    %% Connect users to API
    U1 --> A1
    U2 --> API
    U3 --> A2

    %% Connect API to Data Services
    A1 --> DS1
    A1 --> DS2
    A2 --> DS1
    A2 --> DS2

    %% Connect Data Services to Cassandra
    DS1 --> DB
    DS2 --> DB

    classDef users fill:#B3E5FC,stroke:#0277BD,color:#000000,font-size:16px
    classDef api fill:#FFB74D,stroke:#E65100,color:#000000,font-size:16px
    classDef dataservice fill:#CE93D8,stroke:#6A1B9A,color:#000000,font-size:16px
    classDef cassandra fill:#81C784,stroke:#2E7D32,color:#000000,font-size:16px
    classDef component fill:#E0E0E0,stroke:#424242,color:#000000,font-size:16px

    class U1,U2,U3 users
    class A1,A2 api
    class DS1,DS2,R1,R2,R3,R4,R5,R6,RC1,RC2,DQ1,DQ2 dataservice
    class C1,C2,C3 cassandra
```

## Technical Learnings

### Go Concurrency Patterns
1. **Channels**
   - Used for async communication between goroutines
   - Each request is in its own goroutine with a channel for response from the query executer goroutine

2. **Mutex Operations**
   - Implemented thread-safe access to shared resources

3. **Atomic Operations**
   - Used lock-free atomic counters for metrics tracking

4. **WaitGroups**
   - Used for waiting on multiple goroutines to complete in the CLI client

5. **Context Management**
   - Used context for request timeouts and cancellation


### gRPC Implementation
- Defined service interfaces using Protocol Buffers
- Managed timeout handling using context

### Docker and Container Orchestration
- Implemented health checks for service readiness
- Managed container dependencies and startup order
- Configured networking between services
- Implemented volume management for data persistence

## Future Improvements

1. **Monitoring & Observability**
   - Add distributed tracing

2. **Scalability**
   - Implement dynamic service discovery (e.g. Consul or etcd)

3. **Resilience**
   - Add circuit breakers
   - Implement retry policies
   - Add rate limiting

## License

This project is licensed under the MIT License.
