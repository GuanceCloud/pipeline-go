// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

// Package stats used to record pl metrics
//go:build !(windows && 386)

package refertable

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestConcurrentSQLiteLifecycle(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	p := &PlReferTablesSqlite{db: db}
	defer p.Close()
	tables := []referTable{{TableName: "review", ColumnName: []string{"id"}, ColumnType: []string{columnTypeInt}, RowData: [][]any{{int64(1)}}}}
	if err := p.updateAll(tables); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				p.Query("review", nil, nil, []string{"id"})
				p.Stats()
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 30; i++ {
			if err := p.updateAll(tables); err != nil {
				t.Error(err)
				return
			}
		}
	}()
	wg.Wait()
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	if _, ok := p.Query("review", nil, nil, nil); ok {
		t.Fatal("query after close succeeded")
	}
}

func TestSQLiteCanceledUpdatePreservesMemoryData(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	p := &PlReferTablesSqlite{db: db}
	defer p.Close()
	old := []referTable{{TableName: "old", ColumnName: []string{"id"}, ColumnType: []string{columnTypeInt}, RowData: [][]any{{int64(7)}}}}
	if err := p.updateAll(old); err != nil {
		t.Fatal(err)
	}
	rows := make([][]any, 100000)
	for i := range rows {
		rows[i] = []any{int64(i)}
	}
	next := []referTable{{TableName: "next", ColumnName: []string{"id"}, ColumnType: []string{columnTypeInt}, RowData: rows}}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err = p.updateAllContext(ctx, next)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("update: %v", err)
	}
	got, ok := p.Query("old", []string{"id"}, []any{int64(7)}, []string{"id"})
	if !ok || got["id"] != int64(7) {
		t.Fatalf("old data lost: %v %v", got, ok)
	}
	if err := p.updateAll(old); err != nil {
		t.Fatalf("connection unusable: %v", err)
	}
}

func TestCanceledInMemoryUpdateKeepsSnapshot(t *testing.T) {
	p := &PlReferTablesInMemory{}
	tables := []referTable{{TableName: "old", ColumnName: []string{"id"}, ColumnType: []string{columnTypeInt}, RowData: [][]any{{int64(7)}}}}
	if err := p.updateAll(tables); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := p.updateAllContext(ctx, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("update: %v", err)
	}
	if _, ok := p.Query("old", []string{"id"}, []any{int64(7)}, []string{"id"}); !ok {
		t.Fatal("old snapshot lost")
	}
}
