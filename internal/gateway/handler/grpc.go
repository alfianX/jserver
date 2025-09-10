package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/structpb"
)

// Get or create grpc.ClientConn and grpcreflect.Client per target address
func (p *service) getReflectClient(addr string) (*grpc.ClientConn, *grpcreflect.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	rClient := grpcreflect.NewClientAuto(context.Background(), conn)

	return conn, rClient, nil
}

// ProxyRequest menerima:
// - grpcAddr: address grpc service
// - fullMethod: "ServiceName/MethodName" (case sensitive)
// - ctx: gin.Context or http.Request context
// - jsonBody: JSON request body sebagai []byte
//
// Return JSON response atau error
func (p *service) ProxyRequestDynamic(ctx context.Context, grpcAddr, fullMethod string, jsonBody map[string]interface{}) ([]byte, error) {
	conn, rClient, err := p.getReflectClient(grpcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed connect grpc: %w", err)
	}

	// Parse full method: e.g. "package.Service/Method"
	parts := strings.Split(fullMethod, "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid method format, expected Service/Method")
	}
	svcName, methodName := parts[0], parts[1]

	svcDesc, err := rClient.ResolveService(svcName)
	if err != nil {
		return nil, fmt.Errorf("service %s not found: %w", svcName, err)
	}

	method := svcDesc.FindMethodByName(methodName)
	if method == nil {
		return nil, fmt.Errorf("method %s not found in %s", methodName, svcName)
	}

	// inputType := method.GetInputType()
	// fileProto := inputType.GetFile().AsFileDescriptorProto()

	// fileDesc, err := protodesc.NewFile(fileProto, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to convert file descriptor: %w", err)
	// }

	// msgDesc := fileDesc.Messages().ByName(protoreflect.Name(inputType.GetName()))
	// if msgDesc == nil {
	// 	return nil, fmt.Errorf("message descriptor not found: %s", inputType.GetName())
	// }

	// var m map[string]interface{}
	// if err := json.Unmarshal(jsonBody, &m), err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	// }

	structMsg, err := structpb.NewStruct(jsonBody)
	if err != nil {
		return nil, fmt.Errorf("failed to convert map to structpb.Struct: %w", err)
	}

	// Prepare grpc client and invoke
	methodFullName := fmt.Sprintf("/%s/%s", svcName, methodName)

	var header, trailer metadata.MD
	out := &structpb.Struct{}

	// Dynamic invoke via grpc.Invoke
	err = grpc.Invoke(ctx, methodFullName, structMsg, out, conn, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return nil, err
	}

	// Check if out contains content_type: text/html
	if val, ok := out.AsMap()["content_type"]; ok {
		if ctStr, ok := val.(string); ok && strings.HasPrefix(ctStr, "text/html") {
			htmlStr := ""
			if body, ok := out.AsMap()["body"].(string); ok {
				htmlStr = body
			}
			return fmt.Appendf(nil, `{"__html__": %q, "__content_type__": %q}`, htmlStr, ctStr), nil
		}
	}

	// Marshal response to JSON
	respJSON, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("marshal response error: %w", err)
	}

	return respJSON, nil
}

func (p *service) ProxyRequest(ctx context.Context, grpcAddr, fullMethod string, jsonBody map[string]interface{}) ([]byte, error) {
	conn, rClient, err := p.getReflectClient(grpcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed connect grpc: %w", err)
	}

	// Parse full method: e.g. "package.Service/Method"
	parts := strings.Split(fullMethod, "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid method format, expected Service/Method")
	}
	svcName, methodName := parts[0], parts[1]

	svcDesc, err := rClient.ResolveService(svcName)
	if err != nil {
		return nil, fmt.Errorf("service %s not found: %w", svcName, err)
	}

	method := svcDesc.FindMethodByName(methodName)
	if method == nil {
		return nil, fmt.Errorf("method %s not found in %s", methodName, svcName)
	}

	inputMsg := dynamicpb.NewMessage(method.GetInputType().UnwrapMessage())

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, fmt.Errorf("marshal json input: %w", err)
	}

	err = protojson.Unmarshal(jsonData, inputMsg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to dynamicpb: %w", err)
	}

	// Prepare grpc client and invoke
	methodFullName := fmt.Sprintf("/%s/%s", svcName, methodName)

	var header, trailer metadata.MD
	out := dynamicpb.NewMessage(method.GetOutputType().UnwrapMessage())

	// Dynamic invoke via grpc.Invoke
	err = grpc.Invoke(ctx, methodFullName, inputMsg, out, conn, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return nil, err
	}

	// Marshal response to JSON
	respJSON, err := protojson.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("marshal response error: %w", err)
	}

	return respJSON, nil
}
