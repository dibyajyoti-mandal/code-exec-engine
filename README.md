# Scalable Code Execution Engine

A distributed system for executing user-submitted code in isolated environments with real-time feedback. Designed for competitive programming contests and online judge platforms.

[Design Documentation](https://github.com/dibyajyoti-mandal/code-exec-engine)

## System Architecture

The system consists of the following main components:

### Client
- Submits code using a Monaco-based code editor frontend
- Provides a modern, VS Code-like editing experience

### Gateway Server
- Receives all code execution jobs from clients
- Pushes jobs into the job queue for processing
- Acts as the entry point for all code submissions

### Job Queue
- Acts as a buffer between the API and Workers
- Ensures requests are not lost if workers are busy
- Provides reliable message delivery and retry mechanisms

### Worker Engine
- Subscribes to the job queue
- Polls for pending jobs
- Spawns worker nodes to execute code
- Manages concurrent execution with configurable limits

### Worker Nodes
- Execute code in isolated Docker containers
- Enforce time and memory limits
- Provide network isolation for security

### WebSocket
- Provides real-time updates for submitted jobs
- Enables live feedback to users during code execution

### Database
- Stores code execution results
- Persists code files for future reference

## Tech Stack

- **Frontend**: React with Monaco editor library
- **Backend**: Go (Worker, Gateway, and WebSocket services)
- **Message Queue**: Redis Streams
- **Execution Runtime**: Docker containers

## Key Design Choices

### Why Go?

Go is the ideal choice for this system for several reasons:

- **Network Performance**: Go is widely used in network services. Its `net/http` library and goroutine model allow it to handle tens of thousands of concurrent connections with very little memory overhead. This is crucial when thousands of users submit code simultaneously during a contest.

- **Strong Standard Library**: Reduces dependency on external packages, improving reliability and maintainability.

- **Performance**: Unlike Python/Node-based applications, Go compiles to machine code. This reduces the CPU cycle cost per request, providing low-latency service.

### Redis Streams as Message Queue

Redis Streams provides several advantages:

- **Speed**: Operations happen in memory over a persistent TCP connection. Pushing a job (`XADD`) and popping it (`XREADGROUP`) takes microseconds.

- **Scalability**: You can spin up multiple Worker instances all reading from the job stream (`REDIS_STREAM`). Redis ensures that each job is delivered to only one worker.

- **Reliability**: If a worker crashes while processing a job, it can retry the job once it restarts.

### Docker Container as Execution Runtime

Docker provides the perfect balance of isolation, performance, and ease of use:

- **Isolation**: Submitted code runs in an isolated Linux environment and cannot see other processes, access the host network, or reach internal APIs. Changes are destroyed on exit, ensuring a clean state for every submission.

- **Resource Control**: Easy to enforce memory and time limit control to judge MLE (Memory Limit Exceeded) or TLE (Time Limit Exceeded) responses.

- **Performance**: Fast boot time and low memory overhead allow high density (running hundreds of containers per server).

- **No Cold Start**: Compared to services like AWS Lambda (based on Firecracker microVM), Docker doesn't have cold-start issues and network restrictions are easier to implement.

**Example**: Script to run C++ code and Docker command that runs the script in the corresponding container. We use `timeout`, `--memory`, and `--network` parameters/flags to ensure network isolation and time/memory limits.

## Concurrency Handling

Concurrency is handled through a **Fan-Out / Fan-In pattern** using Goroutines and Channels, with a specific focus on limiting concurrency per worker.

### The Concurrency Pattern

**Fan-Out**: 
- We spawn a fixed number (`WORKER_COUNT`) of workers as separate goroutines working in the background.
- We use buffered channels (`jobQueue`, `workerLimiter`, and `resultQueue`) to pass data safely between Goroutines without locks.

**Channel Responsibilities**:
- `jobQueue`: Acts as the input buffer. It consumes jobs from the Redis stream.
- `workerLimiter`: Limits the number of jobs each worker can execute concurrently.
- `resultQueue`: Acts as the output buffer. Workers write results here, and the `resultBroadcaster` reads them.

## Future Improvements for Production

### Database Integration
- Store execution results for users
- Set up systems to check against given test cases
- The current version solely focuses on isolated and safe execution of submitted code

### Object Storage
- Source code and test cases cannot be stored in databases or passed as payload
- Use object storage like AWS S3 or MinIO for code files and test cases
- Suitable for scenarios of competitive programming contests

### Advanced Message Queue
- Migrate to RabbitMQ or Kafka for topic-based routing of jobs
- Route code jobs based on language/difficulty of question to workers with different configurations
- Redis Streams don't offer this feature

### Enhanced Security with gVisor
- Use Google gVisor as a secure sandbox for Docker containers
- Intercepts application system calls and runs them in a user-space kernel
- Provides a strong defense-in-depth layer for Docker containers without changing the entire infrastructure

### Auto-Scaling with Kubernetes
- Workers are currently a static set of Go processes that require manual scaling
- Containerize the Worker service itself and deploy it on Kubernetes (K8s)
- Use Horizontal Pod Autoscaler (HPA) for auto-scaling based on the number of pending jobs in the message queue

## Getting Started

```bash
# Run the entire system
./run.sh
```

## License

MIT
