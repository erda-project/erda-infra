// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

import (
	"context"
	"fmt"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// checkDebugOptions of protocol.
func checkDebugOptions(ctx context.Context, debugOptions *cptype.ComponentProtocolDebugOptions) error {
	if debugOptions == nil {
		return nil
	}
	if debugOptions.ComponentKey == "" {
		return fmt.Errorf(i18n(ctx, "debugoptions.missing.componentkey"))
	}
	return nil
}
