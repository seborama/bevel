# bevel

[![Build Status](https://travis-ci.org/seborama/bevel.svg?branch=master)](https://travis-ci.org/seborama/bevel)

An angle on Business Event Logger in Golang.

**bevel** is a simple and extensible module that provides a generalised approach to logging business events in a larger application. Events can be vehicled to any number of writers (file, Kafka topic, etc).

Examples are supplied to get started in minutes.

An application only needs to create a message for the business event and post it to the **bevel** bus. The bus will relay the messages to the registered writers. That's all it takes!

By design, **bevel** can be used as an event loop dispatcher: the listeners are the **bevel writers** and the dispatcher is the **bevel Manager** (where **events** are sent to for disptaching).

## Installation

To install, use the standard Go installation procedure:

```bash
go get github.com/seborama/bevel
```

You can pick a specific major release for compatibility. For example, to use a v1.x release, use this command:

```bash
go get gopkg.in/seborama/bevel.v1
```

And your source code would use this import:

```go
import "gopkg.in/seborama/bevel.v1"
```

## Documentation

The code documentation can be found on [godoc](http://godoc.org/github.com/seborama/bevel).

## Project contents

- Bevel:
  - Listener
  - WriterPool
- Writers:
  - ConsoleBEWriter - An example Console writer.
  - KafkaBEWriter - An example Kafka writer.

## High level architecture

Messages are posted to the Manager's listener loop.

The listener passes messages to each of the registered Writers.

Writers are then free to process and persist messages as they please.

                               Manager                             Writers
                                                             ___________________
                                            Write(Message)  |                   |
                                         />>>>>>>>>>>>>>>>>>|  KafkaBEWriter    |
                                       //                   |___________________|
                             __________||                    ___________________
             Post(Message)  |            |  Write(Message)  |                   |
    Message  >>>>>>>>>>>>>  |  Listener  |>>>>>>>>>>>>>>>>>>|  ConsoleBEWriter  |
                            |____________|                  |___________________|
                                       ||                    ___________________
                                       \\   Write(Message)  |                   |
                                         \>>>>>>>>>>>>>>>>>>|  Other BE Writer  |
                                                            |___________________|


## Usage

For a simple example of usage, please see [main_test.go](https://github.com/seborama/bevel/blob/0.1/main_test.go).

The example defines a `CounterMsg` that acts as a business event.

### Step 1 - Create a Message structure

To get started, we create a message structure to hold our information about the Business Event.

Our message must embed the `StandardMessage` structure as demonstrated below.

`StandardMessage` implements  `Message`, an interface that is consumed by `Writer` implementations.

```go
type CounterMsg struct {
    bevel.StandardMessage
    Counter int
}
```

In this example, our `CounterMsg` simply holds a `Counter`.

### Step 2 - Ignite the event listener

We now need to create a listener to receive our `CounterMsg` (which is a `Message` implementor).

```go
    bem := bevel.StartNewListener(&bevel.ConsoleBEWriter{})
    defer func() {
        bem.Done()
    }()
```

This function is at the heart of `bevel`, the Business Events Logger and performs these actions:

1. It registers the supplied Writer (in this instance a simple Console Writer called) in the `WriterPool`.
1. It creates a `Manager` and starts the `Manager`'s listener.
1. Finally, it returns the Manager for our use.

In addition to running the main Business Event listener, a `Manager` offer convenient services:

- To `Post()` messages to the listener.
- To instruct the `Manager` to terminate gracefully - via `Done()`.
- To add more `Writer`'s to the `Manager`'s `WriterPool` - via `AddWriter()`.

### Step 3 - Optionally register more `Writer`s

We can optionally add all the `Writer`'s we wish our listener to write messages to.

Writers are flexible and may:

- use different persistence: log files, databases, message queues (RabbitMQ, Kafka, etc).
- define the data format: CSV, plain text, table columns, queues, etc.
- filter what messages they wish to persist, from a criteria based on the contents of the message: importance, source/type, etc.

### Step 4 - Start generating messages!

We're now ready to send a few `Message`'s to our listener:

```go
    for i := 1; i <= 5; i++ {
        m := CounterMsg{
            StandardMessage: bevel.StandardMessage{
                EventName:         "test_event",
                CreatedTSUnixNano: time.Now().UnixNano(),
            },
            Counter: i,
        }

        bem.Post(m)
    }
```

The key method to note here is `Post()`.

#### Step 5 - When we're done, let the listener know via the `Manager`

This was already prepared via `defer` above, at step 2.
