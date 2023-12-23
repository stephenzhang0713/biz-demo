package rpc

import (
	"context"
	"fmt"
	"github.com/baiyutang/gomall/app/frontend/infra/mtl"
	"github.com/baiyutang/gomall/app/frontend/kitex_gen/cart/cartservice"
	"github.com/baiyutang/gomall/app/frontend/kitex_gen/product"
	"github.com/baiyutang/gomall/app/frontend/kitex_gen/product/productcatalogservice"
	"github.com/baiyutang/gomall/app/frontend/kitex_gen/user/userservice"
	frontendutils "github.com/baiyutang/gomall/app/frontend/utils"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"os"
	"sync"
)

var (
	ProductClient productcatalogservice.Client
	UserClient    userservice.Client
	CartClient    cartservice.Client
	once          sync.Once
	err           error
)

func InitClient() {
	once.Do(func() {
		initProductClient()
		initUserClient()
		initCartClient()
	})
}

func initProductClient() {
	var opts []client.Option
	if os.Getenv("REGISTRY_ENABLE") == "true" {
		r, err := consul.NewConsulResolver(os.Getenv("REGISTRY_ADDR"))
		frontendutils.MustHandleError(err)
		opts = append(opts, client.WithResolver(r))
	} else {
		opts = append(opts, client.WithHostPorts("localhost:8881"))
	}
	p := provider.NewOpenTelemetryProvider(provider.WithSdkTracerProvider(mtl.TracerProvider), provider.WithEnableMetrics(false))
	defer p.Shutdown(context.Background())
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: frontendutils.ServiceName}), client.WithSuite(tracing.NewClientSuite()))

	cbs := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string {
		return circuitbreak.RPCInfo2Key(ri)
	})
	cbs.UpdateServiceCBConfig("shop-frontend/product/GetProduct", circuitbreak.CBConfig{Enable: true, ErrRate: 0.5, MinSample: 2})

	opts = append(opts, client.WithCircuitBreaker(cbs), client.WithFallback(fallback.NewFallbackPolicy(fallback.UnwrapHelper(func(ctx context.Context, req, resp interface{}, err error) (fbResp interface{}, fbErr error) {
		methodName := rpcinfo.GetRPCInfo(ctx).To().Method()
		if err == nil {
			return resp, err
		}
		if methodName != "ListProducts" {
			return resp, err
		}
		return &product.ListProductsResponse{
			Products: []*product.Product{{
				Price:       6.6,
				Id:          3,
				Picture:     "/static/image/t-shirt.jpeg",
				Name:        "T-Shirt",
				Description: "CloudWeGo T-Shirt"},
			},
		}, nil

	}))))

	ProductClient, err = productcatalogservice.NewClient("product", opts...)
	frontendutils.MustHandleError(err)
}

func initUserClient() {
	var opts []client.Option
	if os.Getenv("REGISTRY_ENABLE") == "true" {
		r, err := consul.NewConsulResolver(os.Getenv("REGISTRY_ADDR"))
		frontendutils.MustHandleError(err)
		opts = append(opts, client.WithResolver(r))
	} else {
		opts = append(opts, client.WithHostPorts("localhost:8882"))
	}

	UserClient, err = userservice.NewClient("user", opts...)
	frontendutils.MustHandleError(err)
}

func initCartClient() {
	var opts []client.Option
	//if os.Getenv("REGISTRY_ENABLE") == "true" {
	//	r, err := consul.NewConsulResolver(os.Getenv("REGISTRY_ADDR"))
	//	frontendutils.MustHandleError(err)
	//	opts = append(opts, client.WithResolver(r))
	//} else {
	opts = append(opts, client.WithHostPorts("localhost:8883"))
	//}
	p := provider.NewOpenTelemetryProvider(provider.WithSdkTracerProvider(mtl.TracerProvider), provider.WithEnableMetrics(false))
	defer p.Shutdown(context.Background())
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: frontendutils.ServiceName}), client.WithSuite(tracing.NewClientSuite()))

	CartClient, err = cartservice.NewClient("cart", opts...)
	fmt.Println(err, "cart client error")
	frontendutils.MustHandleError(err)
}
