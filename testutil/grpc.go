package testutil

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/k1LoW/grpcstub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Cacert = func() []byte {
	b, err := os.ReadFile(filepath.Join(Testdata(), "cacert.pem"))
	if err != nil {
		panic(err)
	}
	return b
}()

var Cert = func() []byte {
	b, err := os.ReadFile(filepath.Join(Testdata(), "cert.pem"))
	if err != nil {
		panic(err)
	}
	return b
}()

var Key = func() []byte {
	b, err := os.ReadFile(filepath.Join(Testdata(), "key.pem"))
	if err != nil {
		panic(err)
	}
	return b
}()

func GRPCServer(t *testing.T, useTLS bool, disableReflection bool) *grpcstub.Server {
	pf := filepath.Join(Testdata(), "grpctest.proto")
	opts := []grpcstub.Option{
		grpcstub.EnableHealthCheck(),
		grpcstub.BufLock(filepath.Join(Testdata(), "buf.lock")),
	}
	if disableReflection {
		opts = append(opts, grpcstub.DisableReflection())
	}
	if useTLS {
		opts = append(opts, grpcstub.UseTLS(Cacert, Cert, Key))
	}
	ts := grpcstub.NewServer(t, pf, opts...)
	t.Cleanup(func() {
		ts.Close()
	})

	// slow responses
	ts.Method("grpctest.GrpcTestService/Hello").Match(func(r *grpcstub.Request) bool {
		h := r.Headers.Get("slow")
		if len(h) > 0 {
			time.Sleep(10 * time.Second)
			return true
		}
		return false
	}).Response(func() map[string]any {
		return map[string]any{}
	}())
	ts.Method("grpctest.GrpcTestService/ListHello").Match(func(r *grpcstub.Request) bool {
		h := r.Headers.Get("slow")
		if len(h) > 0 {
			time.Sleep(10 * time.Second)
			return true
		}
		return false
	}).Response(func() map[string]any {
		return map[string]any{}
	}())
	ts.Method("grpctest.GrpcTestService/MultiHello").Match(func(r *grpcstub.Request) bool {
		h := r.Headers.Get("slow")
		if len(h) > 0 {
			time.Sleep(10 * time.Second)
			return true
		}
		return false
	}).Response(func() map[string]any {
		return map[string]any{}
	}())
	ts.Method("grpctest.GrpcTestService/HelloChat").Match(func(r *grpcstub.Request) bool {
		h := r.Headers.Get("slow")
		if len(h) > 0 {
			time.Sleep(10 * time.Second)
			return true
		}
		return false
	}).Response(func() map[string]any {
		return map[string]any{}
	}())

	// error responses
	ts.Method("grpctest.GrpcTestService/Hello").Match(func(r *grpcstub.Request) bool {
		h := r.Headers.Get("error")
		return len(h) > 0
	}).Status(status.New(codes.Canceled, "request canceled"))

	// default responses
	ts.Method("grpctest.GrpcTestService/Hello").
		Header("hello", "header").Trailer("hello", "trailer").
		ResponseString(`{"message":"hello", "num":32, "create_time":"2022-06-25T05:24:43.861872Z"}`)
	ts.Method("grpctest.GrpcTestService/ListHello").
		Header("listhello", "header").Trailer("listhello", "trailer").
		ResponseString(`{"message":"hello", "num":33, "create_time":"2022-06-25T05:24:43.861872Z"}`).
		ResponseString(`{"message":"hello", "num":34, "create_time":"2022-06-25T05:24:44.382783Z"}`)
	ts.Method("grpctest.GrpcTestService/MultiHello").
		Header("multihello", "header").Trailer("multihello", "trailer").
		ResponseString(`{"message":"hello", "num":35, "create_time":"2022-06-25T05:24:45.382783Z"}`)
	ts.Method("grpctest.GrpcTestService/HelloChat").Match(func(r *grpcstub.Request) bool {
		n, ok := r.Message["name"]
		if !ok {
			return false
		}
		ns, ok := n.(string)
		if !ok {
			return false
		}
		return ns == "alice"
	}).Header("hellochat", "header").Trailer("hellochat", "trailer").
		ResponseString(`{"message":"hello", "num":34, "create_time":"2022-06-25T05:24:46.382783Z"}`)
	ts.Method("grpctest.GrpcTestService/HelloChat").Match(func(r *grpcstub.Request) bool {
		n, ok := r.Message["name"]
		if !ok {
			return false
		}
		ns, ok := n.(string)
		if !ok {
			return false
		}
		return ns == "bob"
	}).Header("hellochat-second", "header").Trailer("hellochat-second", "trailer").
		ResponseString(`{"message":"hello", "num":35, "create_time":"2022-06-25T05:24:47.382783Z"}`)
	ts.Method("grpctest.GrpcTestService/HelloChat").Match(func(r *grpcstub.Request) bool {
		n, ok := r.Message["name"]
		if !ok {
			return false
		}
		ns, ok := n.(string)
		if !ok {
			return false
		}
		return ns == "charlie"
	}).Header("hellochat-third", "header").Trailer("hellochat-second", "trailer").
		ResponseString(`{"message":"hello", "num":36, "create_time":"2022-06-25T05:24:48.382783Z"}`)
	ts.Method("grpctest.GrpcTestService/HelloFields").Match(func(r *grpcstub.Request) bool { return true }).
		Handler(func(r *grpcstub.Request) *grpcstub.Response {
			return &grpcstub.Response{
				Messages: []grpcstub.Message{
					{
						"field_bytes": r.Message["field_bytes"],
					},
				},
			}
		})

	return ts
}
