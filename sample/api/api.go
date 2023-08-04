package api

import (
	"fmt"
	"log"
	"net"
	"net/http"

	auth "github.com/AccelByte/go-restful-plugins/v4/pkg/auth/iam"
	"github.com/AccelByte/iam-go-sdk"
	"github.com/AccelByte/observability-go-sdk/metrics"
	"github.com/emicklei/go-restful/v3"
)

func InitWebService(basePath string) *WebService {
	iamClient := iam.NewMockClient()

	authFilterOptions := auth.FilterInitializationOptionsFromEnv()
	authFilter := auth.NewFilterWithOptions(iamClient, authFilterOptions)

	serviceContainer := newServiceContainer(basePath, authFilter)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", "8080"))
	if err != nil {
		log.Fatalf("unable to listen on port 8080: %s", err.Error())
	}

	return &WebService{
		serviceContainer: serviceContainer,
		listener:         listener,
	}
}

func newServiceContainer(basePath string, authFilter *auth.Filter) *restful.Container {
	container := restful.NewContainer()

	// register metrics route
	container.Handle(basePath+"/metrics", metrics.PrometheusHandler())
	// register runtime debug route
	container.Add(metrics.NewRuntimeDebugRoute(basePath).WebService(authFilter))
	return container
}

type WebService struct {
	serviceContainer *restful.Container
	listener         net.Listener
}

func (w *WebService) Serve() {
	if err := http.Serve(w.listener, w.serviceContainer); err != nil {
		log.Fatal("unable to serve: ", err)
	}
}