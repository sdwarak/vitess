// Copyright 2015, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vtgateconntest provides the test methods to make sure a
// vtgateconn/vtgateservice pair over RPC works correctly.
package vtgateconntest

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	mproto "github.com/youtube/vitess/go/mysql/proto"
	"github.com/youtube/vitess/go/sqltypes"
	"github.com/youtube/vitess/go/vt/key"
	tproto "github.com/youtube/vitess/go/vt/tabletserver/proto"
	"github.com/youtube/vitess/go/vt/topo"
	"github.com/youtube/vitess/go/vt/vtgate/proto"
	"github.com/youtube/vitess/go/vt/vtgate/vtgateconn"
	"github.com/youtube/vitess/go/vt/vtgate/vtgateservice"
	"golang.org/x/net/context"
)

// fakeVTGateService has the server side of this fake
type fakeVTGateService struct {
	t      *testing.T
	panics bool
}

// Execute is part of the VTGateService interface
func (f *fakeVTGateService) Execute(ctx context.Context, query *proto.Query, reply *proto.QueryResult) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	execCase, ok := execMap[query.Sql]
	if !ok {
		return fmt.Errorf("no match for: %s", query.Sql)
	}
	if !reflect.DeepEqual(query, execCase.execQuery) {
		f.t.Errorf("Execute: %+v, want %+v", query, execCase.execQuery)
		return nil
	}
	*reply = *execCase.reply
	return nil
}

// ExecuteShard is part of the VTGateService interface
func (f *fakeVTGateService) ExecuteShard(ctx context.Context, query *proto.QueryShard, reply *proto.QueryResult) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	execCase, ok := execMap[query.Sql]
	if !ok {
		return fmt.Errorf("no match for: %s", query.Sql)
	}
	if !reflect.DeepEqual(query, execCase.shardQuery) {
		f.t.Errorf("Execute: %+v, want %+v", query, execCase.shardQuery)
		return nil
	}
	*reply = *execCase.reply
	return nil
}

// ExecuteKeyspaceIds is part of the VTGateService interface
func (f *fakeVTGateService) ExecuteKeyspaceIds(ctx context.Context, query *proto.KeyspaceIdQuery, reply *proto.QueryResult) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	return nil
}

// ExecuteKeyRanges is part of the VTGateService interface
func (f *fakeVTGateService) ExecuteKeyRanges(ctx context.Context, query *proto.KeyRangeQuery, reply *proto.QueryResult) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	return nil
}

// ExecuteEntityIds is part of the VTGateService interface
func (f *fakeVTGateService) ExecuteEntityIds(ctx context.Context, query *proto.EntityIdsQuery, reply *proto.QueryResult) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	return nil
}

// ExecuteBatchShard is part of the VTGateService interface
func (f *fakeVTGateService) ExecuteBatchShard(ctx context.Context, batchQuery *proto.BatchQueryShard, reply *proto.QueryResultList) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	return nil
}

// ExecuteBatchKeyspaceIds is part of the VTGateService interface
func (f *fakeVTGateService) ExecuteBatchKeyspaceIds(ctx context.Context, batchQuery *proto.KeyspaceIdBatchQuery, reply *proto.QueryResultList) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	return nil
}

// StreamExecute is part of the VTGateService interface
func (f *fakeVTGateService) StreamExecute(ctx context.Context, query *proto.Query, sendReply func(*proto.QueryResult) error) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	execCase, ok := execMap[query.Sql]
	if !ok {
		return fmt.Errorf("no match for: %s", query.Sql)
	}
	if !reflect.DeepEqual(query, execCase.execQuery) {
		f.t.Errorf("Execute: %+v, want %+v", query, execCase.execQuery)
		return nil
	}
	if execCase.reply.Result != nil {
		result := proto.QueryResult{Result: &mproto.QueryResult{}}
		result.Result.Fields = execCase.reply.Result.Fields
		if err := sendReply(&result); err != nil {
			return err
		}
		for _, row := range execCase.reply.Result.Rows {
			result := proto.QueryResult{Result: &mproto.QueryResult{}}
			result.Result.Rows = [][]sqltypes.Value{row}
			if err := sendReply(&result); err != nil {
				return err
			}
		}
	}
	if execCase.reply.Error != "" {
		return errors.New(execCase.reply.Error)
	}
	return nil
}

// StreamExecuteShard is part of the VTGateService interface
func (f *fakeVTGateService) StreamExecuteShard(ctx context.Context, query *proto.QueryShard, sendReply func(*proto.QueryResult) error) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	return nil
}

// StreamExecuteKeyRanges is part of the VTGateService interface
func (f *fakeVTGateService) StreamExecuteKeyRanges(ctx context.Context, query *proto.KeyRangeQuery, sendReply func(*proto.QueryResult) error) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	return nil
}

// StreamExecuteKeyspaceIds is part of the VTGateService interface
func (f *fakeVTGateService) StreamExecuteKeyspaceIds(ctx context.Context, query *proto.KeyspaceIdQuery, sendReply func(*proto.QueryResult) error) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	return nil
}

// Begin is part of the VTGateService interface
func (f *fakeVTGateService) Begin(ctx context.Context, outSession *proto.Session) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	*outSession = *session1
	return nil
}

// Commit is part of the VTGateService interface
func (f *fakeVTGateService) Commit(ctx context.Context, inSession *proto.Session) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	if !reflect.DeepEqual(inSession, session2) {
		return errors.New("commit: session mismatch")
	}
	return nil
}

// Rollback is part of the VTGateService interface
func (f *fakeVTGateService) Rollback(ctx context.Context, inSession *proto.Session) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	if !reflect.DeepEqual(inSession, session2) {
		return errors.New("rollback: session mismatch")
	}
	return nil
}

// SplitQuery is part of the VTGateService interface
func (f *fakeVTGateService) SplitQuery(ctx context.Context, req *proto.SplitQueryRequest, reply *proto.SplitQueryResult) error {
	if f.panics {
		panic(fmt.Errorf("test forced panic"))
	}
	if !reflect.DeepEqual(req, splitQueryRequest) {
		f.t.Errorf("SplitQuery has wrong input: got %#v wanted %#v", req, splitQueryRequest)
	}
	*reply = *splitQueryResult
	return nil
}

// CreateFakeServer returns the fake server for the tests
func CreateFakeServer(t *testing.T) vtgateservice.VTGateService {
	return &fakeVTGateService{
		t:      t,
		panics: false,
	}
}

// HandlePanic is part of the VTGateService interface
func (f *fakeVTGateService) HandlePanic(err *error) {
	if x := recover(); x != nil {
		*err = fmt.Errorf("uncaught panic: %v", x)
	}
}

// TestSuite runs all the tests
func TestSuite(t *testing.T, conn vtgateconn.VTGateConn, fakeServer vtgateservice.VTGateService) {
	testExecute(t, conn)
	testExecuteShard(t, conn)
	testStreamExecute(t, conn)
	testTxPass(t, conn)
	testTxFail(t, conn)
	testSplitQuery(t, conn)

	// force a panic at every call, then test that works
	fakeServer.(*fakeVTGateService).panics = true
	testExecutePanic(t, conn)
	testExecuteShardPanic(t, conn)
	testStreamExecutePanic(t, conn)
	testBeginPanic(t, conn)
	testSplitQueryPanic(t, conn)
}

func expectPanic(t *testing.T, err error) {
	expected1 := "test forced panic"
	expected2 := "uncaught panic"
	if err == nil || !strings.Contains(err.Error(), expected1) || !strings.Contains(err.Error(), expected2) {
		t.Fatalf("Expected a panic error with '%v' or '%v' but got: %v", expected1, expected2, err)
	}
}

func testExecute(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	execCase := execMap["request1"]
	qr, err := conn.Execute(ctx, execCase.execQuery.Sql, execCase.execQuery.BindVariables, execCase.execQuery.TabletType)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(qr, execCase.reply.Result) {
		t.Errorf("Unexpected result from Execute: got %+v want %+v", qr, execCase.reply.Result)
	}

	_, err = conn.Execute(ctx, "none", nil, "")
	want := "no match for: none"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("none request: %v, want %v", err, want)
	}

	_, err = conn.Execute(ctx, "errorRequst", nil, "")
	want = "app error"
	if err == nil || err.Error() != want {
		t.Errorf("errorRequst: %v, want %v", err, want)
	}
}

func testExecutePanic(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	execCase := execMap["request1"]
	_, err := conn.Execute(ctx, execCase.execQuery.Sql, execCase.execQuery.BindVariables, execCase.execQuery.TabletType)
	expectPanic(t, err)
}

func testExecuteShard(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	execCase := execMap["request1"]
	qr, err := conn.ExecuteShard(ctx, execCase.execQuery.Sql, "ks", []string{"1", "2"}, execCase.execQuery.BindVariables, execCase.execQuery.TabletType)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(qr, execCase.reply.Result) {
		t.Errorf("Unexpected result from Execute: got %+v want %+v", qr, execCase.reply.Result)
	}

	_, err = conn.ExecuteShard(ctx, "none", "", []string{}, nil, "")
	want := "no match for: none"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("none request: %v, want %v", err, want)
	}

	_, err = conn.ExecuteShard(ctx, "errorRequst", "", []string{}, nil, "")
	want = "app error"
	if err == nil || err.Error() != want {
		t.Errorf("errorRequst: %v, want %v", err, want)
	}
}

func testExecuteShardPanic(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	execCase := execMap["request1"]
	_, err := conn.ExecuteShard(ctx, execCase.execQuery.Sql, "ks", []string{"1", "2"}, execCase.execQuery.BindVariables, execCase.execQuery.TabletType)
	expectPanic(t, err)
}

func testStreamExecute(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	execCase := execMap["request1"]
	packets, errFunc := conn.StreamExecute(ctx, execCase.execQuery.Sql, execCase.execQuery.BindVariables, execCase.execQuery.TabletType)
	var qr mproto.QueryResult
	for packet := range packets {
		if len(packet.Fields) != 0 {
			qr.Fields = packet.Fields
		}
		if len(packet.Rows) != 0 {
			qr.Rows = append(qr.Rows, packet.Rows...)
		}
	}
	wantResult := *execCase.reply.Result
	wantResult.RowsAffected = 0
	wantResult.InsertId = 0
	if !reflect.DeepEqual(qr, wantResult) {
		t.Errorf("Unexpected result from Execute: got %+v want %+v", qr, wantResult)
	}
	err := errFunc()
	if err != nil {
		t.Error(err)
	}

	packets, errFunc = conn.StreamExecute(ctx, "none", nil, "")
	for packet := range packets {
		t.Errorf("packet: %+v, want none", packet)
	}
	err = errFunc()
	want := "no match for: none"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("none request: %v, want %v", err, want)
	}

	packets, errFunc = conn.StreamExecute(ctx, "errorRequst", nil, "")
	for packet := range packets {
		t.Errorf("packet: %+v, want none", packet)
	}
	err = errFunc()
	want = "app error"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("errorRequst: %v, want %v", err, want)
	}
}

func testStreamExecutePanic(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	execCase := execMap["request1"]
	packets, errFunc := conn.StreamExecute(ctx, execCase.execQuery.Sql, execCase.execQuery.BindVariables, execCase.execQuery.TabletType)
	if _, ok := <-packets; ok {
		t.Fatalf("Received packets instead of panic?")
	}
	err := errFunc()
	expectPanic(t, err)
}

func testTxPass(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Error(err)
	}
	execCase := execMap["txRequest"]
	_, err = tx.Execute(ctx, execCase.execQuery.Sql, execCase.execQuery.BindVariables, execCase.execQuery.TabletType)
	if err != nil {
		t.Error(err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		t.Error(err)
	}

	tx, err = conn.Begin(ctx)
	if err != nil {
		t.Error(err)
	}
	execCase = execMap["txRequest"]
	_, err = tx.ExecuteShard(ctx, execCase.shardQuery.Sql, execCase.shardQuery.Keyspace, execCase.shardQuery.Shards, execCase.shardQuery.BindVariables, execCase.shardQuery.TabletType)
	if err != nil {
		t.Error(err)
	}
	err = tx.Rollback(ctx)
	if err != nil {
		t.Error(err)
	}
}

func testBeginPanic(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	_, err := conn.Begin(ctx)
	expectPanic(t, err)
}

func testTxFail(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Error(err)
	}
	err = tx.Commit(ctx)
	want := "commit: session mismatch"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("Commit: %v, want %v", err, want)
	}

	_, err = tx.Execute(ctx, "", nil, "")
	want = "execute: not in transaction"
	if err == nil || err.Error() != want {
		t.Errorf("Execute: %v, want %v", err, want)
	}

	_, err = tx.ExecuteShard(ctx, "", "", nil, nil, "")
	want = "executeShard: not in transaction"
	if err == nil || err.Error() != want {
		t.Errorf("ExecuteShard: %v, want %v", err, want)
	}

	err = tx.Commit(ctx)
	want = "commit: not in transaction"
	if err == nil || err.Error() != want {
		t.Errorf("Commit: %v, want %v", err, want)
	}

	err = tx.Rollback(ctx)
	if err != nil {
		t.Error(err)
	}

	tx, err = conn.Begin(ctx)
	if err != nil {
		t.Error(err)
	}
	err = tx.Rollback(ctx)
	want = "rollback: session mismatch"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("Rollback: %v, want %v", err, want)
	}
}

func testSplitQuery(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	qsl, err := conn.SplitQuery(ctx, splitQueryRequest.Keyspace, splitQueryRequest.Query, splitQueryRequest.SplitCount)
	if err != nil {
		t.Fatalf("SplitQuery failed: %v", err)
	}
	if !reflect.DeepEqual(qsl, splitQueryResult.Splits) {
		t.Errorf("SplitQuery returned worng result: got %v wanted %v", qsl, splitQueryResult.Splits)
	}
}

func testSplitQueryPanic(t *testing.T, conn vtgateconn.VTGateConn) {
	ctx := context.Background()
	_, err := conn.SplitQuery(ctx, splitQueryRequest.Keyspace, splitQueryRequest.Query, splitQueryRequest.SplitCount)
	expectPanic(t, err)
}

var execMap = map[string]struct {
	execQuery  *proto.Query
	shardQuery *proto.QueryShard
	reply      *proto.QueryResult
	err        error
}{
	"request1": {
		execQuery: &proto.Query{
			Sql: "request1",
			BindVariables: map[string]interface{}{
				"bind1": int64(0),
			},
			TabletType: topo.TYPE_RDONLY,
			Session:    nil,
		},
		shardQuery: &proto.QueryShard{
			Sql: "request1",
			BindVariables: map[string]interface{}{
				"bind1": int64(0),
			},
			Keyspace:   "ks",
			Shards:     []string{"1", "2"},
			TabletType: topo.TYPE_RDONLY,
			Session:    nil,
		},
		reply: &proto.QueryResult{
			Result:  &result1,
			Session: nil,
			Error:   "",
		},
	},
	"errorRequst": {
		execQuery: &proto.Query{
			Sql:           "errorRequst",
			BindVariables: map[string]interface{}{},
			TabletType:    "",
			Session:       nil,
		},
		shardQuery: &proto.QueryShard{
			Sql:           "errorRequst",
			BindVariables: map[string]interface{}{},
			TabletType:    "",
			Keyspace:      "",
			Shards:        []string{},
			Session:       nil,
		},
		reply: &proto.QueryResult{
			Result:  nil,
			Session: nil,
			Error:   "app error",
		},
	},
	"txRequest": {
		execQuery: &proto.Query{
			Sql:           "txRequest",
			BindVariables: map[string]interface{}{},
			TabletType:    "",
			Session:       session1,
		},
		shardQuery: &proto.QueryShard{
			Sql:           "txRequest",
			BindVariables: map[string]interface{}{},
			TabletType:    "",
			Keyspace:      "",
			Shards:        []string{},
			Session:       session1,
		},
		reply: &proto.QueryResult{
			Result:  nil,
			Session: session2,
			Error:   "",
		},
	},
}

var result1 = mproto.QueryResult{
	Fields: []mproto.Field{
		mproto.Field{
			Name: "field1",
			Type: 42,
		},
		mproto.Field{
			Name: "field2",
			Type: 73,
		},
	},
	RowsAffected: 123,
	InsertId:     72,
	Rows: [][]sqltypes.Value{
		[]sqltypes.Value{
			sqltypes.MakeString([]byte("row1 value1")),
			sqltypes.MakeString([]byte("row1 value2")),
		},
		[]sqltypes.Value{
			sqltypes.MakeString([]byte("row2 value1")),
			sqltypes.MakeString([]byte("row2 value2")),
		},
	},
}

var session1 = &proto.Session{
	InTransaction: true,
	ShardSessions: []*proto.ShardSession{},
}

var session2 = &proto.Session{
	InTransaction: true,
	ShardSessions: []*proto.ShardSession{
		&proto.ShardSession{
			Keyspace:      "ks",
			Shard:         "1",
			TabletType:    topo.TYPE_MASTER,
			TransactionId: 1,
		},
	},
}

var splitQueryRequest = &proto.SplitQueryRequest{
	Keyspace: "ks",
	Query: tproto.BoundQuery{
		Sql: "in for SplitQuery",
		BindVariables: map[string]interface{}{
			"bind1": int64(43),
		},
	},
	SplitCount: 13,
}

var splitQueryResult = &proto.SplitQueryResult{
	Splits: []proto.SplitQueryPart{
		proto.SplitQueryPart{
			Query: &proto.KeyRangeQuery{
				Sql: "out for SplitQuery",
				BindVariables: map[string]interface{}{
					"bind1": int64(1114444),
				},
				Keyspace:   "ksout",
				KeyRanges:  []key.KeyRange{},
				TabletType: topo.TYPE_RDONLY,
			},
			Size: 12344,
		},
	},
}
