/*
 * Copyright (c) 2023 yosefkohan26
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Contributor(s):
 * Andrew Querol <aquerol@yosefkohan26.com>
 */
package command

// Resume is used in outbound socket mode to return control of the call to the FreeSWITCH dialplan.
type Resume struct{}

// BuildMessage builds the ESL command string for the Resume command.
func (Resume) BuildMessage() string {
	return "resume"
}
