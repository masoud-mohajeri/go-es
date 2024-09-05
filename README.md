# Event Sourcing System in Go

This repository contains a highly modularized event sourcing system built in Go, including components for queues, cron jobs, projections, storage, and event dispatching. Event sourcing is a pattern where state changes are stored as a sequence of events, allowing full reconstruction of the state and enabling advanced features like auditing, history tracking, and replaying events.
This is an example use case of the event sourcing system built using Go. The system tracks the lifecycle of an order through various states (e.g., `ORDER_PLACED`, `ORDER_IN_DELIVERY`, etc.) and stores the state history in a file-based storage system. The example demonstrates how to use the event dispatcher, queue, projection, storage, and cron modules to simulate an event-driven architecture.

## Overview

In this use case, we simulate the lifecycle of an order in an e-commerce system using event sourcing. Events are dispatched to represent different states of the order (e.g., `ORDER_PLACED`, `ORDER_IN_TRANSIT`), and each state is persisted in storage. A projection updates the order's current state based on the dispatched events, and a cron job periodically generates snapshots from the event queue.

## Project Components

1. **Queue**: Handles incoming order state events.
2. **Event Dispatcher**: Dispatches order state changes to the event queue.
3. **Projection**: Applies the events from the queue to update and store the current state.
4. **Storage**: Saves the current and previous states of the order.
5. **Cron**: Periodically processes the event queue and generates state snapshots.

## Features

- **Event Queues**: FIFO queue implementation for processing events.
- **Cron Jobs**: Trigger periodic tasks such as snapshot generation.
- **Projections**: Apply transformations or computations on events to build different models (e.g., for read views).
- **Event Dispatcher**: Dispatch events to their respective handlers.
- **Storage**: Event and state persistence for durable and recoverable systems.
- **Concurrency Safe**: Leveraging Go's goroutines for concurrent event processing and projections.

## Modules

### 1. Queue
The `queue` package provides a simple FIFO queue that allows pushing events and processing them one at a time.

- **Push(Event)**: Adds an event to the queue.
- **Pop()**: Retrieves and removes the next event from the queue.

### 2. Cron
The `cron` package allows you to schedule periodic tasks such as snapshot generation or event dispatching at regular intervals. It uses the `gocron` library to handle job scheduling and gracefully handles system signals for shutdown.

#### Methods:
- **NewCron(interval time.Duration, action Action, ctx context.Context)**: Creates a new cron scheduler that runs the provided `action` at the specified `interval`. The scheduler listens for system signals (like `SIGTERM`) and can also be controlled via the provided `context` for graceful shutdown.


### 3. Projection
The `projection` package generates read models from events. Handlers are registered for specific event types to apply transformations or updates.

- **Register(eventType, handlerFunc)**: Register a handler for a specific event type.
- **GenerateSnapshot()**: Generates snapshots by reading from the queue and applying event handlers.

### 4. Storage
The `storage` package provides a simple file-based storage mechanism for persisting event data and snapshots. It uses JSON serialization to save data to text files and supports retrieving the latest and previous values for a given key.

#### Methods:
- **Save(key string, value []byte)**: Saves a new value to the file system. If a previous value exists, it is stored in the `PreviousValues` array, and the new value is set as `LatestValue`. Data is stored in JSON format.
- **Get(key string)**: Retrieves the file contents for the given key. If the file doesn't exist, it returns an empty JSON object.

### 5. Event Dispatcher
The `dispatcher` package is responsible for dispatching events to the event queue for later processing. It allows decoupling of event generation from event handling by inserting events into a shared queue, which can be processed by various parts of the system asynchronously.

#### Methods:
- **Dispatch(payload string, kind string)**: Dispatches an event by adding it to the queue. The `payload` is the event data, and `kind` is the event type (e.g., `userCreated`, `orderPlaced`).
- **NewDispatcher(queue queue.FifoQueue)**: Initializes a new `EventDispatcher` using the provided FIFO queue reference. The same queue instance can be shared between multiple dispatchers.


