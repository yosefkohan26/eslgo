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
package command

import "fmt"

// GetVar is used in outbound socket mode to retrieve a channel variable.
type GetVar struct {
	VariableName string
}

// BuildMessage builds the ESL command string for the GetVar command.
func (g GetVar) BuildMessage() string {
	return fmt.Sprintf("getvar %s", g.VariableName)
}
