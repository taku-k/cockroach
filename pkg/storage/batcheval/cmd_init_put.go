// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package batcheval

import (
	"golang.org/x/net/context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage/batcheval/result"
	"github.com/cockroachdb/cockroach/pkg/storage/engine"
)

// InitPut sets the value for a specified key only if it doesn't exist. It
// returns a ConditionFailedError if the key exists with an existing value that
// is different from the value provided. If FailOnTombstone is set to true,
// tombstones count as mismatched values and will cause a ConditionFailedError.
func InitPut(
	ctx context.Context, batch engine.ReadWriter, cArgs CommandArgs, resp roachpb.Response,
) (result.Result, error) {
	args := cArgs.Args.(*roachpb.InitPutRequest)
	h := cArgs.Header

	if h.DistinctSpans {
		if b, ok := batch.(engine.Batch); ok {
			// Use the distinct batch for both blind and normal ops so that we don't
			// accidentally flush mutations to make them visible to the distinct
			// batch.
			batch = b.Distinct()
			defer batch.Close()
		}
	}
	if args.Blind {
		return result.Result{}, engine.MVCCBlindInitPut(ctx, batch, cArgs.Stats, args.Key, h.Timestamp, args.Value, args.FailOnTombstones, h.Txn)
	}
	return result.Result{}, engine.MVCCInitPut(ctx, batch, cArgs.Stats, args.Key, h.Timestamp, args.Value, args.FailOnTombstones, h.Txn)
}
