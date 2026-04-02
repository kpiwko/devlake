/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*cleanupNullCiSourcePredictions)(nil)

type cleanupNullCiSourcePredictions struct{}

// Up deletes orphaned prediction rows that have ci_failure_source = NULL.
//
// These rows were written before migration 20260331_add_ci_failure_source introduced
// the ci_failure_source column. At that time generatePredictionId hashed (prId, aiTool, ""),
// producing IDs that differ from the IDs produced by subsequent task runs that include
// the ci_failure_source value in the hash. As a result, new task runs always INSERT fresh
// rows rather than updating the old NULL rows, leaving them as permanently stale orphans.
//
// Deleting them here is safe: they will be re-created with correct data and ci_failure_source
// set on the next calculateFailurePredictions task run.
func (script *cleanupNullCiSourcePredictions) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.Exec("DELETE FROM _tool_aireview_failure_predictions WHERE ci_failure_source IS NULL"); err != nil {
		return errors.Default.Wrap(err, "failed to delete orphaned null-ci_failure_source predictions")
	}
	return nil
}

func (script *cleanupNullCiSourcePredictions) Version() uint64 {
	return 20260402000001
}

func (script *cleanupNullCiSourcePredictions) Name() string {
	return "aireview delete orphaned failure predictions with null ci_failure_source"
}
