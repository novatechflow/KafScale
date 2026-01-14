// Copyright 2025, 2026 Alexander Alten (novatechflow), NovaTechflow (novatechflow.com).
// This project is supported and financed by Scalytics, Inc. (www.scalytics.io).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package processor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kafscale/platform/addons/processors/sql-processor/internal/decoder"
	"github.com/kafscale/platform/addons/processors/sql-processor/internal/sink"
)

func TestMapRecords(t *testing.T) {
	records := []decoder.Record{
		{Topic: "orders", Partition: 1, Offset: 10, Value: []byte("a")},
		{Topic: "orders", Partition: 1, Offset: 11, Value: []byte("b")},
	}
	out := mapRecords(records)
	if len(out) != 2 {
		t.Fatalf("expected 2 records, got %d", len(out))
	}
	if out[0].Offset != 10 || string(out[0].Payload) != "a" {
		t.Fatalf("unexpected mapped record: %+v", out[0])
	}
}

func TestFilterRecords(t *testing.T) {
	records := []sink.Record{
		{Offset: 1},
		{Offset: 5},
		{Offset: 10},
	}
	out := filterRecords(records, 5)
	if len(out) != 1 || out[0].Offset != 10 {
		t.Fatalf("unexpected filtered records: %+v", out)
	}
}

func TestTopicLocker(t *testing.T) {
	locker := newTopicLocker()
	unlock := locker.Lock("orders")
	done := make(chan struct{})
	go func() {
		unlock()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatalf("lock release timed out")
	}
}

func TestRunContextCancelClosesSink(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock := &closeSink{}
	p := &Processor{
		sink:  mock,
		locks: newTopicLocker(),
	}
	if err := p.Run(ctx); err != nil {
		t.Fatalf("run: %v", err)
	}
	if !mock.closed() {
		t.Fatalf("expected sink to close")
	}
}

type closeSink struct {
	mu         sync.Mutex
	closedFlag bool
}

func (c *closeSink) Write(ctx context.Context, records []sink.Record) error {
	return nil
}

func (c *closeSink) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closedFlag = true
	return nil
}

func (c *closeSink) closed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closedFlag
}
