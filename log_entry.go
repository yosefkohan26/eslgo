/*
 * Copyright (c) 2023 Percipia
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Contributor(s):
 * Andrew Querol <aquerol@percipia.com>
 */
package eslgo

import (
	"net/textproto"
	"strconv"
)

// LogEntry represents a single log line received from FreeSWITCH.
type LogEntry struct {
	Headers textproto.MIMEHeader
	Body    []byte
}

// GetHeader is a helper function to retrieve a header value.
func (l *LogEntry) GetHeader(name string) string {
	return l.Headers.Get(name)
}

// Level returns the parsed log level of the entry.
func (l *LogEntry) Level() int {
	levelStr := l.GetHeader("Log-Level")
	if level, err := strconv.Atoi(levelStr); err == nil {
		return level
	}
	return 0 // Default to 0 if parsing fails
}

// Message returns the log message body as a string.
func (l *LogEntry) Message() string {
	return string(l.Body)
}
