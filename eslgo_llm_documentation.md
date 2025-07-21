# eslgo Library: A Comprehensive Guide for LLMs

## 1. Introduction

`eslgo` is a powerful and idiomatic Go library for interacting with FreeSWITCH via the Event Socket Layer (ESL). It provides a robust and feature-rich interface for both inbound and outbound ESL connections, enabling developers to build sophisticated telephony applications.

This document provides a comprehensive guide to the `eslgo` library, specifically tailored for Large Language Models (LLMs). It covers the library's core concepts, API reference, and practical examples to enable you to understand and utilize `eslgo` effectively.

## 2. Installation

To use the `eslgo` library in your Go project, you can install it using the following command:

```bash
go get github.com/yosefkohan26/eslgo
```

## 3. Core Concepts

### 3.1. Connections

The `eslgo` library supports two types of ESL connections:

*   **Inbound Connections**: Your application connects to the FreeSWITCH ESL server. This is useful for controlling FreeSWITCH from an external application.
*   **Outbound Connections**: FreeSWITCH connects to your application, which acts as an ESL server. This is typically used for handling incoming calls and executing dialplan applications.

#### 3.1.1. Inbound Connections

To establish an inbound connection, you use the `eslgo.Dial` function:

```go
conn, err := eslgo.Dial("127.0.0.1:8021", "ClueCon", func() {
    fmt.Println("Inbound Connection Disconnected")
})
if err != nil {
    // Handle error
}
defer conn.ExitAndClose()
```

The `eslgo.Dial` function takes the FreeSWITCH ESL address, the ESL password, and an optional `onDisconnect` callback function.

#### 3.1.2. Outbound Connections

To create an outbound ESL server, you use the `eslgo.ListenAndServe` function:

```go
func handleConnection(ctx context.Context, conn *eslgo.Conn, response *eslgo.RawResponse) {
    // Handle the new connection
    // response contains the initial channel data
}

log.Fatalln(eslgo.ListenAndServe(":8084", handleConnection))
```

The `eslgo.ListenAndServe` function takes a listening address and a handler function that will be executed for each new outbound connection from FreeSWITCH.

**Important**: When using outbound sockets with a remote FreeSWITCH server, ensure the socket address in your dialplan uses the public IP address of your Go application server, not `localhost` or `127.0.0.1`.

### 3.2. Commands

You can send commands to FreeSWITCH using the `conn.SendCommand` method. The `eslgo` library provides a set of pre-defined command structs in the `command` package.

All command structs implement the `command.Command` interface, which has a single method: `BuildMessage() string`.

**Important**: All command structs are in the `github.com/yosefkohan26/eslgo/command` package, not the main `eslgo` package.

**Example: Sending an `api` command:**

```go
import "github.com/yosefkohan26/eslgo/command"

response, err := conn.SendCommand(ctx, command.API{
    Command:    "status",
    Background: false,
})
```

**Note**: For commands with arguments, use the `Arguments` field:

```go
response, err := conn.SendCommand(ctx, command.API{
    Command:    "sofia",
    Arguments:  "status profile internal",
    Background: false,
})
```

### 3.3. Events

FreeSWITCH generates events for various activities, such as call state changes, DTMF tones, and custom events. You can listen for these events using event listeners.

To register an event listener, you use the `conn.RegisterEventListener` method:

```go
listenerID := conn.RegisterEventListener(eslgo.EventListenAll, func(event *eslgo.Event) {
    fmt.Printf("Received event: %#v\n", event)
})
```

You can register listeners for all events (`eslgo.EventListenAll`) or for specific channels by providing the channel's UUID.

To enable events, you need to send an Event command:

```go
import "github.com/yosefkohan26/eslgo/command"

// Enable all events in plain format
_, err := conn.SendCommand(ctx, command.Event{
    Format: "plain",
    Listen: []string{"all"},
})

// Or enable JSON format events
_, err := conn.SendCommand(ctx, command.Event{
    Format: "json",
    Listen: []string{"all"},
})
```

**Note**: The helper method `conn.EnableEvents(ctx)` is a convenience function that enables all plain-text events.

### 3.4. Responses

When you send a command to FreeSWITCH, you receive a `RawResponse` object. This object contains the headers and body of the response.

You can check if a response was successful using the `IsOk()` method:

```go
if response.IsOk() {
    // Command was successful
}
```

You can also access response headers and variables using the provided helper methods.

## 4. Exhaustive API Reference

This section provides a complete reference to all public types, functions, and methods in the `eslgo` library and its sub-packages.

### 4.1. Core `eslgo` Package

#### Types & Structs

*   **`type EventListener func(event *Event)`**: The function signature for an event listener callback.
*   **`struct Conn`**: Represents an active ESL connection. This is the central object for interacting with FreeSWITCH. Key internal fields:
    *   `outbound bool`: Whether this is an outbound connection (FreeSWITCH connects to us)
    *   `logger Logger`: The logger instance for internal messages
    *   `exitTimeout time.Duration`: Timeout for the "exit" command when closing
    *   `closeDelay time.Duration`: Delay before closing (set by linger command)
*   **`struct Options`**: Generic options for an ESL connection (inbound or outbound).
    *   `Context context.Context`: The base running context for the connection.
    *   `Logger Logger`: The logger for internal library messages.
    *   `ExitTimeout time.Duration`: How long to wait for an "exit" command to complete.
*   **`struct Event`**: Represents an event received from FreeSWITCH.
    *   `Headers textproto.MIMEHeader`: The event headers.
    *   `Body []byte`: The event body.
*   **`struct RawResponse`**: Represents a raw response from FreeSWITCH to a command.
    *   `Headers textproto.MIMEHeader`: The response headers.
    *   `Body []byte`: The response body.
*   **`struct LogEntry`**: Represents a single log line received from FreeSWITCH.
    *   `Headers textproto.MIMEHeader`: The log entry headers.
    *   `Body []byte`: The log message.
*   **`struct Leg`**: Specifies an individual leg of a call for origination.
    *   `CallURL string`: The dial string for the leg.
    *   `LegVariables map[string]string`: Channel variables specific to this leg.
*   **`type Logger interface`**: An interface for logging.
    *   `Debug(format string, args ...interface{})`
    *   `Info(format string, args ...interface{})`
    *   `Warn(format string, args ...interface{})`
    *   `Error(format string, args ...interface{})`
*   **`struct NormalLogger`**: Default logger implementation that outputs to standard log with level prefixes.
*   **`struct NilLogger`**: No-op logger implementation that discards all log messages.

#### Constants & Variables

*   **`const EventListenAll = "ALL"`**: A special value to listen for all events.
*   **`const EndOfMessage = "\r\n\r\n"`**: The message delimiter used in the ESL protocol.
*   **`var DefaultOptions = Options{...}`**: The default options used for creating a connection (Context: background, Logger: NormalLogger, ExitTimeout: 5s).
*   **Response Type Constants**: 
    *   `TypeEventPlain = "text/event-plain"`: Plain text event format
    *   `TypeEventJSON = "text/event-json"`: JSON event format
    *   `TypeEventXML = "text/event-xml"`: XML event format
    *   `TypeReply = "command/reply"`: Command reply
    *   `TypeAPIResponse = "api/response"`: API command response
    *   `TypeAuthRequest = "auth/request"`: Authentication request
    *   `TypeDisconnect = "text/disconnect-notice"`: Disconnect notification
    *   `TypeLogData = "log/data"`: Log data

#### Functions

*   **`func BuildVars(format string, vars map[string]string) string`**: A helper that builds channel variable strings for commands.

#### `Conn` Methods

*   **`func (c *Conn) LogChannel() <-chan *LogEntry`**: Returns a read-only channel that receives FreeSWITCH log entries.
*   **`func (c *Conn) RegisterEventListener(channelUUID string, listener EventListener) string`**: Registers a new event listener. The channelUUID can be:
    *   `EventListenAll` to listen to all events
    *   A specific channel UUID (from `Unique-ID` header)
    *   An application UUID (from `Application-UUID` header)
    *   A job UUID (from `Job-UUID` header)
*   **`func (c *Conn) RemoveEventListener(channelUUID string, id string)`**: Removes a previously registered event listener.
*   **`func (c *Conn) SendCommand(ctx context.Context, cmd command.Command) (*RawResponse, error)`**: Sends a command to FreeSWITCH.
*   **`func (c *Conn) ExitAndClose()`**: Gracefully sends "exit" and closes the connection.
*   **`func (c *Conn) Close()`**: Immediately closes the connection without sending "exit".
*   **`func (c *Conn) EnableEvents(ctx context.Context) error`**: A helper to subscribe to all plain-text events.
*   **`func (c *Conn) EnableLogs(ctx context.Context, level ...int) (*RawResponse, error)`**: A helper to enable receiving log data.
*   **`func (c *Conn) DisableLogs(ctx context.Context) (*RawResponse, error)`**: A helper to disable receiving log data.
*   **`func (c *Conn) DebugEvents(w io.Writer) string`**: A helper to quickly log all events to a writer.
*   **`func (c *Conn) DebugOff(id string)`**: Disables the debug event listener.
*   **`func (c *Conn) Phrase(ctx context.Context, uuid, macro string, times int, wait bool) (*RawResponse, error)`**: Executes the `phrase` application.
*   **`func (c *Conn) PhraseWithArg(...)`**: Executes `phrase` with an argument.
*   **`func (c *Conn) Playback(...)`**: Executes the `playback` application.
*   **`func (c *Conn) Say(...)`**: Executes the `say` application.
*   **`func (c *Conn) Speak(...)`**: Executes the `speak` application.
*   **`func (c *Conn) WaitForDTMF(ctx context.Context, uuid string) (byte, error)`**: Blocks until a DTMF event is received for a specific UUID.
*   **`func (c *Conn) GetVar(ctx context.Context, varName string) (string, error)`**: Retrieves a channel variable in an outbound connection. This is only valid for outbound connections where FreeSWITCH connects to your application. Returns the variable value or an error.
*   **`func (c *Conn) Resume(ctx context.Context) (*RawResponse, error)`**: Resumes dialplan execution in an outbound connection. After calling this, control returns to the FreeSWITCH dialplan and the connection will close. Only valid for outbound connections.
*   **`func (c *Conn) OriginateCall(ctx context.Context, background bool, aLeg, bLeg Leg, vars map[string]string) (*RawResponse, error)`**: A helper to originate a new call. The `background` parameter determines if the command waits for completion.
*   **`func (c *Conn) EnterpriseOriginateCall(ctx context.Context, background bool, vars map[string]string, bLeg Leg, aLegs ...Leg) (*RawResponse, error)`**: A helper for enterprise origination (multiple A-legs using ":_:" separator).
*   **`func (c *Conn) BackgroundOriginateCall(ctx context.Context, background bool, aLeg, bLeg Leg, vars map[string]string) (*RawResponse, error)`**: A helper to originate a call using bgapi (always asynchronous).
*   **`func (c *Conn) HangupCall(ctx context.Context, uuid, cause string) error`**: A helper to hang up a call asynchronously.
*   **`func (c *Conn) AnswerCall(ctx context.Context, uuid string) error`**: A helper to answer a call synchronously.

#### Other Method Receivers

*   **`func (e Event) GetName() string`**: Returns the `Event-Name` header.
*   **`func (e Event) HasHeader(header string) bool`**: Checks if an event has a specific header.
*   **`func (e Event) GetHeader(header string) string`**: Gets an event header value (automatically URL-unescapes the value).
*   **`func (e Event) String() string`**: Implements Stringer interface for pretty printing.
*   **`func (e Event) GoString() string`**: Implements GoStringer interface for %#v formatting.
*   **`func (r RawResponse) IsOk() bool`**: Checks if a command response is `+OK`.
*   **`func (r RawResponse) GetReply() string`**: Gets the reply text from a response (uses Reply-Text header or body if header doesn't exist).
*   **`func (r RawResponse) ChannelUUID() string`**: Gets the `Unique-ID` header from a response.
*   **`func (r RawResponse) HasHeader(header string) bool`**: Checks if a response has a specific header.
*   **`func (r RawResponse) GetVariable(variable string) string`**: Gets a `Variable_*` header from a response.
*   **`func (r RawResponse) GetHeader(header string) string`**: Gets a response header value (automatically URL-unescapes the value).
*   **`func (r RawResponse) String() string`**: Implements Stringer interface for pretty printing.
*   **`func (r RawResponse) GoString() string`**: Implements GoStringer interface for %#v formatting.
*   **`func (l *LogEntry) GetHeader(name string) string`**: Gets a log entry header value.
*   **`func (l *LogEntry) Level() int`**: Gets the parsed log level.
*   **`func (l *LogEntry) Message() string`**: Gets the log message body.
*   **`func (l Leg) String() string`**: Formats the leg into a dial string.

### 4.2. Inbound Connections

*   **`struct InboundOptions`**: Options for dialing an inbound connection.
    *   `Options`: Embedded generic options
    *   `Network string`: Network type (usually "tcp", "tcp4", or "tcp6")
    *   `Password string`: Authentication password for FreeSWITCH
    *   `OnDisconnect func()`: Optional callback when connection is closed
    *   `AuthTimeout time.Duration`: Timeout for authentication (default 5s)
*   **`var DefaultInboundOptions`**: Default options for inbound connections (Network: "tcp", Password: "ClueCon", AuthTimeout: 5s).
*   **`func Dial(address, password string, onDisconnect func()) (*Conn, error)`**: Connects to FreeSWITCH with simplified parameters.
*   **`func (opts InboundOptions) Dial(address string) (*Conn, error)`**: Connects to FreeSWITCH with specific options.

### 4.3. Outbound Connections

*   **`type OutboundHandler func(ctx context.Context, conn *Conn, connectResponse *RawResponse)`**: The function signature for an outbound connection handler.
*   **`struct OutboundOptions`**: Options for listening for outbound connections.
    *   `Options`: Embedded generic options
    *   `Network string`: Network type to listen on (usually "tcp", "tcp4", or "tcp6")
    *   `ConnectTimeout time.Duration`: Timeout for the "connect" command (default 5s)
    *   `ConnectionDelay time.Duration`: Delay after connection before sending commands (default 25ms)
*   **`var DefaultOutboundOptions`**: Default options for outbound listeners (Network: "tcp", ConnectTimeout: 5s, ConnectionDelay: 25ms).
*   **`func ListenAndServe(address string, handler OutboundHandler) error`**: Listens for connections from FreeSWITCH with default options.
*   **`func (opts OutboundOptions) ListenAndServe(address string, handler OutboundHandler) error`**: Listens for connections with specific options.

### 4.4. `command` Package

This package contains structs that represent specific ESL commands. All implement the `Command` interface.

#### Functions

*   **`func FormatHeaderString(headers textproto.MIMEHeader) string`**: Formats headers for FreeSWITCH ESL protocol (converts \r\n to \n).

#### Types

*   **`interface Command`**: The interface all command structs must implement (`BuildMessage() string`).
*   **`struct API`**: For `api` and `bgapi` commands. Fields:
    *   `Command string`: The API command to execute
    *   `Arguments string`: Command arguments
    *   `Background bool`: If true, uses `bgapi` for async execution
*   **`struct Auth`**: For `auth` and `userauth` commands. Fields:
    *   `User string`: Username (if provided, uses `userauth` command)
    *   `Password string`: Authentication password
*   **`struct Connect`**: For the `connect` command in outbound mode.
*   **`struct Event`**: For the `event` and `nixevent` commands. Fields:
    *   `Ignore bool`: If true, uses `nixevent` to unsubscribe
    *   `Format string`: Event format ("plain", "json", or "xml")
    *   `Listen []string`: Event types to listen for
*   **`struct MyEvents`**: For the `myevents` command. Fields:
    *   `Format string`: Event format
    *   `UUID string`: Optional channel UUID to filter events
*   **`struct DisableEvents`**: For the `noevents` command (no fields).
*   **`struct DivertEvents`**: For the `divert_events` command. Fields:
    *   `Enabled bool`: Whether to enable event diversion
*   **`struct SendEvent`**: For the `sendevent` command. Fields:
    *   `Name string`: Event name
    *   `Headers textproto.MIMEHeader`: Event headers
    *   `Body string`: Event body
*   **`struct Exit`**: For the `exit` command.
*   **`struct Filter`**: For the `filter` command. Fields:
    *   `Delete bool`: Whether to delete filters
    *   `EventHeader string`: The header to filter on
    *   `FilterValue string`: The value to filter for
*   **`struct Linger`**: For the `linger` and `nolinger` commands. Fields:
    *   `Enabled bool`: Whether to enable lingering
    *   `Seconds time.Duration`: How long to linger (0 means indefinite)
*   **`struct Log`**: For the `log` and `nolog` commands. Fields:
    *   `Enabled bool`: Whether to enable logging
    *   `Level int`: Log level (0-7)
*   **`struct SendMessage`**: For the `sendmsg` command. Fields:
    *   `UUID string`: Channel UUID to send message to
    *   `Headers textproto.MIMEHeader`: Message headers
    *   `Body string`: Message body
    *   `Sync bool`: Wait for event to finish (event-lock)
    *   `SyncPri bool`: Priority event lock (event-lock-pri)
*   **`struct Resume`**: For the `resume` command.
*   **`struct GetVar`**: For the `getvar` command (outbound connections only). Fields:
    *   `VariableName string`: The name of the channel variable to retrieve.

### 4.5. `command/call` Package

This package contains structs for commands that are sent via `sendmsg` to control a specific call.

*   **`struct Execute`**: Executes a dialplan application. Fields:
    *   `UUID string`: Channel UUID
    *   `AppName string`: Application name to execute
    *   `AppArgs string`: Application arguments
    *   `AppUUID string`: Optional UUID to track application execution
    *   `Loops int`: Number of times to execute (default 1)
    *   `Sync bool`: Wait for completion
    *   `SyncPri bool`: Priority sync
    *   `ForceBody bool`: Force arguments into body
*   **`struct Set`**: A helper to execute the `set` application. Fields:
    *   `UUID string`: Channel UUID
    *   `Key string`: Variable name
    *   `Value string`: Variable value
    *   `Sync bool`: Wait for completion
    *   `SyncPri bool`: Priority sync
*   **`struct Export`**: A helper to execute the `export` application (same fields as Set).
*   **`struct Push`**: A helper to execute the `push` application (same fields as Set).
*   **`struct Hangup`**: Hangs up a call. Fields:
    *   `UUID string`: Channel UUID
    *   `Cause string`: Hangup cause
    *   `Sync bool`: Wait for completion
    *   `SyncPri bool`: Priority sync
*   **`struct NoMedia`**: Removes media from a call. Fields:
    *   `UUID string`: Channel UUID
    *   `NoMediaUUID string`: NoMedia UUID
    *   `Sync bool`: Wait for completion
    *   `SyncPri bool`: Priority sync
*   **`struct Transfer`**: Transfers a call to an application. Fields:
    *   `UUID string`: Channel UUID
    *   `Application string`: Application to transfer to
    *   `Sync bool`: Wait for completion
    *   `SyncPri bool`: Priority sync
*   **`struct Unicast`**: Bridges media to an external process (useful for mod_spandsp faxing). Fields:
    *   `UUID string`: Channel UUID
    *   `Local net.Addr`: Local address
    *   `Remote net.Addr`: Remote address
    *   `Flags string`: Unicast flags
    *   `Sync bool`: Wait for completion
    *   `SyncPri bool`: Priority sync

## 5. Examples

### 5.1. Inbound Connection: Originate a Call

```go
package main

import (
	"context"
	"fmt"
	"github.com/yosefkohan26/eslgo"
	"github.com/yosefkohan26/eslgo/command"
	"time"
)

func main() {
	conn, err := eslgo.Dial("127.0.0.1:8021", "ClueCon", nil)
	if err != nil {
		panic(err)
	}
	defer conn.ExitAndClose()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Using the helper method
	response, err := conn.OriginateCall(
		ctx,
		false,
		eslgo.Leg{CallURL: "user/1000"},
		eslgo.Leg{CallURL: "&echo"},
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Originate response: %#v\n", response)
	
	// Or using the API command directly
	response, err = conn.SendCommand(ctx, command.API{
		Command:    "originate",
		Arguments:  "user/1000 &echo",
		Background: false,
	})
}
```

### 5.2. Outbound Server: Answer and Hangup

```go
package main

import (
	"context"
	"fmt"
	"github.com/yosefkohan26/eslgo"
	"log"
	"time"
)

func handleConnection(ctx context.Context, conn *eslgo.Conn, response *eslgo.RawResponse) {
	fmt.Printf("New connection: %#v\n", response)

	uuid := response.ChannelUUID()
	err := conn.AnswerCall(ctx, uuid)
	if err != nil {
		log.Printf("Error answering call: %v", err)
		return
	}

	log.Println("Call answered, waiting 5 seconds...")
	time.Sleep(5 * time.Second)

	err = conn.HangupCall(ctx, uuid, "NORMAL_CLEARING")
	if err != nil {
		log.Printf("Error hanging up call: %v", err)
	}
}

func main() {
	log.Fatalln(eslgo.ListenAndServe(":8084", handleConnection))
}
```

### 5.3. Event Handling: Listening for DTMF

```go
package main

import (
	"context"
	"fmt"
	"github.com/yosefkohan26/eslgo"
	"time"
)

func main() {
	conn, err := eslgo.Dial("127.0.0.1:8021", "ClueCon", nil)
	if err != nil {
		panic(err)
	}
	defer conn.ExitAndClose()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = conn.EnableEvents(ctx)
	if err != nil {
		panic(err)
	}

	// Originate a call to a user and wait for DTMF
	response, err := conn.OriginateCall(
		ctx,
		true, // Background originate
		eslgo.Leg{CallURL: "user/1000"},
		eslgo.Leg{CallURL: "&playback(misc/demo_dtmf.wav)"},
		nil,
	)
	if err != nil {
		panic(err)
	}

	uuid := response.ChannelUUID()
	fmt.Printf("Call originated with UUID: %s\n", uuid)

	fmt.Println("Waiting for DTMF...")
	digit, err := conn.WaitForDTMF(ctx, uuid)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received DTMF digit: %c\n", digit)
}
```

## 6. Common Pitfalls and Important Notes

### 6.1. Import Paths
- Command structs (API, Event, etc.) are in `github.com/yosefkohan26/eslgo/command`
- Call-specific commands (Execute, Set, etc.) are in `github.com/yosefkohan26/eslgo/command/call`
- The main connection types and functions are in `github.com/yosefkohan26/eslgo`

### 6.2. Command Structure
- The `command.API` struct has separate `Command` and `Arguments` fields
- Don't combine them in the `Command` field - use `Arguments` for parameters

### 6.3. Outbound Socket Considerations
- `GetVar` and `Resume` commands are only valid for outbound connections
- When using outbound sockets with remote FreeSWITCH, use your server's public IP in the dialplan, not localhost
- After calling `Resume()`, the connection will close as control returns to the dialplan

### 6.4. Event Formats
- JSON events have the event data properly formatted as JSON in the `Body` field
- Use `command.Event{Format: "json", Listen: []string{"all"}}` to enable JSON events
- The helper method `conn.EnableEvents(ctx)` enables plain-text events

## 7. Known Limitations and Missing Features

A direct comparison of the `eslgo` library against FreeSWITCH's `mod_event_socket.c` source code reveals a few features of the Event Socket Layer protocol that are not currently implemented in `eslgo`. For a truly comprehensive understanding, it is important to be aware of these gaps:

*   **Incomplete XML Event Parsing**: While you can request events in `xml` format, `eslgo`'s `readXMLEvent` function is a stub and does not parse XML events.
*   **Limited `userauth` Response Handling**: The `userauth` command in FreeSWITCH can reply with headers that restrict the client's permissions (e.g., `Allowed-Events`, `Allowed-API`). The `eslgo` library does not parse these headers or enforce these restrictions on the client side.

This comprehensive guide should provide you with all the necessary information to effectively use the `eslgo` library for your FreeSWITCH integration needs, keeping in mind the limitations listed above.
