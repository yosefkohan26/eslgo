# `eslgo` Future Development Roadmap

## 1. Introduction

This document outlines a roadmap for the future development of the `eslgo` library. It is based on a detailed comparison between the library's current feature set and the full capabilities of the FreeSWITCH Event Socket Layer, as defined in the `mod_event_socket.c` source code. The goal is to identify all missing features, from major architectural components to minor command options, to guide `eslgo` toward becoming a fully comprehensive ESL client implementation.

---

## 2. Major Missing Features

These are high-level capabilities of the ESL protocol that are completely absent in the current version of `eslgo`. Implementing these would significantly expand the library's power and utility.

### 2.1. Full Event Format Parsing (XML)

-   **Current State**: `eslgo` allows a user to request events in `xml` format, but the `readXMLEvent` function is a stub.
-   **Missing Functionality**: The library cannot currently understand or process events delivered in XML format.
-   **Proposed Action**: Implement robust parsing for XML events. This would involve using the standard Go `encoding/xml` library to deserialize the event body into the `eslgo.Event` struct, populating the `Headers` map from the XML data.

### 2.2. Event Sink API Client

-   **Current State**: `eslgo` has no knowledge of the "Event Sink" functionality.
-   **Missing Functionality**: `mod_event_socket` provides a complete HTTP-based API (`event_sink`) for creating, managing, and polling temporary, stateful event listeners. This is a powerful mechanism for web applications to interact with FreeSWITCH without maintaining a persistent socket connection.
-   **Proposed Action**: Create a new, separate client within the `eslgo` package (e.g., in a sub-package like `eslgo/eventsink`) designed to interact with this HTTP API. This client would have methods corresponding to the `event_sink` commands:
    -   `CreateListener(events []string, format string, loglevel string) (*StatefulListener, error)`
    -   `DestroyListener(listenerID uint32) error`
    -   `CheckListener(listenerID uint32) ([]*Event, []*LogEntry, error)`
    -   Methods for adding/deleting filters and managing log levels on a stateful listener.

---

## 3. Command and Feature Enhancements

These are existing features in `eslgo` that are incomplete or could be enhanced to fully match the capabilities of `mod_event_socket`.

### 4.1. `auth`/`userauth` - Full Permissions Handling

-   **Current State**: The `doAuth` method only checks if the authentication reply is `+OK`.
-   **Missing Functionality**: The `userauth` command can return `Allowed-Events`, `Allowed-API`, and `Allowed-LOG` headers that define the permissions for the authenticated session. `eslgo` does not parse, store, or act on these permissions.
-   **Proposed Action**:
    1.  Extend the `Conn` struct to include fields for storing allowed events and APIs (e.g., `allowedEvents map[string]struct{}`, `allowedAPIs map[string]struct{}`).
    2.  Modify the `doAuth` function to parse these headers from the authentication reply and populate the new fields on the `Conn` struct.
    3.  (Advanced) Implement client-side checks in `SendCommand` and `RegisterEventListener` to return an error if the user attempts an action they are not permissioned for.

### 4.2. `api`/`bgapi` - `console_execute` Support

-   **Missing Functionality**: The C code shows that the `api` command can accept a `console_execute: true` header, which changes how the command is executed internally by FreeSWITCH.
-   **Proposed Action**: Add a `ConsoleExecute bool` field to the `command.API` struct. When true, the `BuildMessage` method should prepend the `console_execute: true` header to the command. This is a minor but important feature for full parity.

---

## 5. Outbound Connection Enhancements

### 5.1. Advanced Dialplan Helpers

-   **Missing Functionality**: The `socket` dialplan application in FreeSWITCH has more capabilities than `eslgo`'s helpers suggest. Specifically, it can try multiple destination hosts (`<ip1>:<port1>|<ip2>:<port2>`) and its behavior can be modified by the `socket_resume` channel variable.
-   **Proposed Action**:
    1.  Create a new helper function, e.g., `BuildSocketDialplanString(hosts []string, mode string) string`, that generates the correct string for use in a dialplan.
    2.  Add documentation explaining how to use the `socket_resume` variable in conjunction with an outbound `eslgo` server.

---

## 6. Conclusion

While `eslgo` is a powerful and well-written library, this analysis reveals significant opportunities for growth. By tackling the major missing features like full event parsing and log handling, and then filling in the smaller gaps in command support and feature enhancements, `eslgo` can evolve into a truly complete and comprehensive tool for any Go developer working with FreeSWITCH. This document should serve as a clear and actionable guide for that evolution.
