package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/fancar/tmp_xm/internal/api/auth"
	"github.com/fancar/tmp_xm/internal/config"
	"github.com/fancar/tmp_xm/internal/storage"
	"github.com/fancar/tmp_xm/static"
	log "github.com/sirupsen/logrus"
)

var (
	bind            string
	tlsCert         string
	tlsKey          string
	jwtSecret       string
	corsAllowOrigin string
)

// Setup configures the API endpoints.
func Setup(ctx context.Context, conf config.Config) error {
	if conf.ExternalAPI.JWTSecret == "" {
		return fmt.Errorf("jwt_secret must be set")
	}

	bind = conf.ExternalAPI.Bind
	tlsCert = conf.ExternalAPI.TLSCert
	tlsKey = conf.ExternalAPI.TLSKey
	jwtSecret = conf.ExternalAPI.JWTSecret
	corsAllowOrigin = conf.ExternalAPI.CORSAllowOrigin

	// init grpc server and register it
	validator := auth.NewJWTValidator(storage.DB(), "HS256", jwtSecret)
	// ctx := context.Background()
	grpcServer := grpc.NewServer() // getgRPCServerOptions()...

	// RegisterInternalServiceServer(grpcServer, NewMainAPI()) // temp no validator
	RegisterCompanyServiceServer(grpcServer, NewCompanyAPI(validator))

	return startHTTPServer(ctx, conf, grpcServer)
}

// startHTTPServer init http1/http2 servers
// setup the client http interface variable
// we need to start the gRPC service first, as it is used by the
// grpc-gateway
func startHTTPServer(ctx context.Context,
	conf config.Config, grpcServer *grpc.Server) error {

	if grpcServer == nil {
		return fmt.Errorf("grpcServer is nil")
	}
	var err error
	var clientHTTPHandler http.Handler
	// switch between gRPC and "plain" http handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 &&
			strings.Contains(
				r.Header.Get("Content-Type"),
				"application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			if clientHTTPHandler == nil {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}

			if corsAllowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
				w.Header().Set("Access-Control-Allow-Methods",
					"POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers",
					"Accept, Content-Type, Content-Length, Accept-Encoding, Grpc-Metadata-Authorization")

				if r.Method == "OPTIONS" {
					return
				}
			}

			clientHTTPHandler.ServeHTTP(w, r)
		}
	})

	// start the API server
	go func() {
		log.WithFields(log.Fields{
			"bind":     bind,
			"tls-cert": tlsCert,
			"tls-key":  tlsKey,
		}).Info("api/external: starting api server ...")

		if tlsCert == "" || tlsKey == "" {
			log.Fatal(http.ListenAndServe(bind, h2c.NewHandler(handler, &http2.Server{})))
		} else {
			log.Fatal(http.ListenAndServeTLS(
				bind,
				tlsCert,
				tlsKey,
				h2c.NewHandler(handler, &http2.Server{}),
			))
		}
	}()

	// setup the HTTP handler
	clientHTTPHandler, err = setupHTTPAPI(conf)
	if err != nil {
		return err
	}

	return nil
}

func setupHTTPAPI(conf config.Config) (http.Handler, error) {
	r := mux.NewRouter()

	// setup json api handler
	jsonHandler, err := getJSONGateway(context.Background())
	if err != nil {
		return nil, err
	}

	log.WithField("path", "/api").Info("api/external: registering /api endpoint")
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		data, err := static.FS.ReadFile("swagger/index.html")
		if err != nil {
			log.WithError(err).Error("get swagger template error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}).Methods("get")
	r.PathPrefix("/api").Handler(jsonHandler)

	// setup static file server
	r.PathPrefix("/").Handler(http.FileServer(http.FS(static.FS)))

	return r, nil
}

func getJSONGateway(ctx context.Context) (http.Handler, error) {
	// dial options for the grpc-gateway
	var grpcDialOpts []grpc.DialOption

	if tlsCert == "" || tlsKey == "" {
		grpcDialOpts = append(grpcDialOpts, grpc.WithInsecure())
	} else {
		b, err := ioutil.ReadFile(tlsCert)
		if err != nil {
			return nil, errors.Wrap(err, "read external api tls cert error")
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, errors.Wrap(err, "failed to append certificate")
		}
		grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
				RootCAs:            cp,
			})))
	}

	bindParts := strings.SplitN(bind, ":", 2)
	if len(bindParts) != 2 {
		log.Fatal("get port from bind failed")
	}
	apiEndpoint := fmt.Sprintf("localhost:%s", bindParts[1])

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			EnumsAsInts:  false,
			EmitDefaults: true,
		},
	))

	if err := RegisterCompanyServiceHandlerFromEndpoint(
		ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register application handler error")
	}

	return mux, nil
}
