package gqlapolloenginetracing

import (
	"bytes"
	"compress/gzip"
	"container/list"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/99designs/gqlgen/graphql"
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
)

var _ graphql.Tracer = (*tracerImpl)(nil)

func NewTracer(Logger *zap.Logger) graphql.Tracer {
	return &tracerImpl{
		logger: Logger,
	}
}

var timeNowFunc = time.Now

var ctxTracingKey = &struct{ tmp string }{}
var ctxExecSpanKey = &struct{ tmp string }{}

type tracerImpl struct {
	logger *zap.Logger
}

func getTracingData(ctx context.Context) *tracingData {
	return ctx.Value(ctxTracingKey).(*tracingData)
}

func getExecutionSpan(ctx context.Context) *executionSpan {
	return ctx.Value(ctxExecSpanKey).(*executionSpan)
}

func (t *tracerImpl) StartOperationParsing(ctx context.Context) context.Context {
	now := timeNowFunc()
	td := &tracingData{
		StartTime: now,
		Parsing: &startOffset{
			StartTime: now,
		},
	}
	ctx = context.WithValue(ctx, ctxTracingKey, td)
	return ctx
}

func (t *tracerImpl) EndOperationParsing(ctx context.Context) {
	td := getTracingData(ctx)
	td.Parsing.EndTime = timeNowFunc()
}

func (t *tracerImpl) StartOperationValidation(ctx context.Context) context.Context {
	td := getTracingData(ctx)
	td.Validation = &startOffset{}
	td.Validation.StartTime = timeNowFunc()
	return ctx
}

func (t *tracerImpl) EndOperationValidation(ctx context.Context) {
	td := getTracingData(ctx)
	td.Validation.EndTime = timeNowFunc()
}

func (t *tracerImpl) StartOperationExecution(ctx context.Context) context.Context {
	return ctx
}

func (t *tracerImpl) StartFieldExecution(ctx context.Context, field graphql.CollectedField) context.Context {
	td := getTracingData(ctx)
	es := &executionSpan{
		startOffset: startOffset{
			StartTime: timeNowFunc(),
		},
		ParentType: field.ObjectDefinition.Name,
		FieldName:  field.Name,
		ReturnType: field.Definition.Type.String(),
	}
	ctx = context.WithValue(ctx, ctxExecSpanKey, es)
	td.mu.Lock()
	defer td.mu.Unlock()
	if td.Execution == nil {
		td.Execution = &execution{}
	}
	td.Execution.Resolvers = append(td.Execution.Resolvers, es)

	return ctx
}

func (t *tracerImpl) StartFieldResolverExecution(ctx context.Context, rc *graphql.ResolverContext) context.Context {
	es := getExecutionSpan(ctx)
	es.Path = rc.Path()

	return ctx
}

func (t *tracerImpl) StartFieldChildExecution(ctx context.Context) context.Context {
	return ctx
}

func (t *tracerImpl) EndFieldExecution(ctx context.Context) {
	es := getExecutionSpan(ctx)
	es.EndTime = timeNowFunc()
}

func (t *tracerImpl) EndOperationExecution(ctx context.Context) {
	td := getTracingData(ctx)
	td.EndTime = timeNowFunc()
	td.prepare()

	rootTree := t.buildTree(td)
	rootTrace := t.buildTraceTree(&Trace_Node{}, rootTree)

	requestContext := graphql.GetRequestContext(ctx)
	reducedWithSpace := printWithReducedWhitespace(requestContext)

	traces := map[string]*Traces{}
	startTime, err := ptypes.TimestampProto(td.StartTime)
	if err != nil {
		return
	}
	endTime, err := ptypes.TimestampProto(td.EndTime)
	if err != nil {
		return
	}

	traces["# "+requestContext.OperationName+"\n"+reducedWithSpace] = &Traces{
		Trace: []*Trace{
			&Trace{
				StartTime:  startTime,
				EndTime:    endTime,
				DurationNs: uint64(td.Duration.Nanoseconds()),
				Root:       rootTrace,
				Http: &Trace_HTTP{
					Method: Trace_HTTP_POST,
				},
			},
		},
	}

	fullTracesReport := &FullTracesReport{
		Header: &ReportHeader{
			Hostname:       "Nuttapats-MacBook-Pro.local",
			AgentVersion:   "gqlgen-apollo-engine-tracer@0.1.0",
			RuntimeVersion: "golang v" + runtime.Version(),
			Uname:          runtime.GOOS + ", " + runtime.GOARCH,
			SchemaHash:     "e68feced787086fd9242d217ab84f303ae271cebbb5dffb7fbbc0deeb4b2be137a28bbb102d96b70f8cc18369cde94db29e31c0b7a3f0d27b4d7e3178253c9ad",
		},
		TracesPerQuery: traces,
	}

	t.sendReport(ctx, fullTracesReport)
}

func (t *tracerImpl) sendReport(ctx context.Context, report *FullTracesReport) (*http.Response, error) {
	var err error

	data, err := proto.Marshal(report)
	if err != nil {
		t.logger.Error("marshaling error:", zap.Error(err))
		return nil, err
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err = g.Write(data); err != nil {
		t.logger.Error("gzip write error:", zap.Error(err))
		return nil, err
	}
	if err = g.Close(); err != nil {
		t.logger.Error("gzip close error:", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://engine-report.apollodata.com/api/ingress/traces", &buf)
	req = req.WithContext(ctx)
	req.Header.Set("user-agent", "gqlgen-apollo-engine-reporting")
	req.Header.Set("content-encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.logger.Error("http request err:", zap.Error(err))
		return nil, err
	}
	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.logger.Error("body read err:", zap.Any("status", resp.Status), zap.Error(err))
			return nil, err
		}
		t.logger.Error("http request err:", zap.Any("status", resp.Status), zap.ByteString("body", body))
		return nil, errors.New("HTTP status " + resp.Status)
	}
	return resp, nil
}

func (t *tracerImpl) buildTree(td *tracingData) *treeNode {
	rootTree := newTreeNode()
	for _, resolve := range td.Execution.Resolvers {
		currentNode := rootTree

		for _, path := range resolve.Path {
			_, ok := currentNode.children[path]
			if !ok {
				currentNode.children[path] = newTreeNode()
			}

			currentNode = currentNode.children[path]
		}
		currentNode.node = resolve
	}
	return rootTree
}

func (t *tracerImpl) printTree(path interface{}, node *treeNode) {
	t.logger.Debug("node", zap.Any("path", path), zap.Int("node children len", len(node.children)), zap.Any("node", node.node))
	for key, child := range node.children {
		t.printTree(key, child)
	}
}

// Reconstruct tree using BFS
func (t *tracerImpl) buildTraceTree(rootTrace *Trace_Node, rootTree *treeNode) *Trace_Node {
	type QueueEntry struct {
		trace *Trace_Node
		node  *treeNode
	}

	queue := list.New()
	queue.PushBack(QueueEntry{rootTrace, rootTree})
	for queue.Len() > 0 {
		t := queue.Front()
		queue.Remove(t)

		entry := t.Value.(QueueEntry)
		trace := entry.trace
		node := entry.node

		if node != rootTree && node.node != nil {
			trace.OriginalFieldName = node.node.FieldName
			trace.Type = node.node.ReturnType
			trace.StartTime = uint64(node.node.StartTime.Nanosecond())
			trace.EndTime = uint64(node.node.EndTime.Nanosecond())
			trace.ParentType = node.node.ParentType
		}
		children := make([]*Trace_Node, 0, len(node.children))
		for key, child := range node.children {
			if child == nil {
				continue
			}
			var id isTrace_Node_Id
			switch keyT := key.(type) {
			case int:
				id = &Trace_Node_Index{
					Index: uint32(keyT),
				}
				break
			case string:
				id = &Trace_Node_ResponseName{
					ResponseName: keyT,
				}
			}
			childTrace := &Trace_Node{
				Id: id,
			}
			children = append(children, childTrace)
			queue.PushBack(QueueEntry{childTrace, child})
		}
		trace.Child = children
	}
	return rootTrace
}
