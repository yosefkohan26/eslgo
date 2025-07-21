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

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetVar_BuildMessage(t *testing.T) {
	getVar := GetVar{VariableName: "test_var"}
	assert.Equal(t, "getvar test_var", getVar.BuildMessage())
}
