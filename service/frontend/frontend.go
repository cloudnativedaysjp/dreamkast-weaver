package frontend

import (
	"dreamkast-weaver/service/cfpsvc"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/ServiceWeaver/weaver"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	cookieMaxAge = 60 * 60 * 48

	cookiePrefix    = "shop_"
	cookieSessionID = cookiePrefix + "session-id"
	cookieCurrency  = cookiePrefix + "currency"
)

var (
	validEnvs = []string{"local", "gcp"}
)

type platformDetails struct {
	css      string
	provider string
}

func (plat *platformDetails) setPlatformDetails(env string) {
	if env == "gcp" {
		plat.provider = "Google Cloud"
		plat.css = "gcp-platform"
	} else {
		plat.provider = "local"
		plat.css = "local"
	}
}

// Server is the application frontend.
type Server struct {
	handler  http.Handler
	root     weaver.Instance
	platform platformDetails
	hostname string

	cfpSvc cfpsvc.T
}

// NewServer returns the new application frontend.
func NewServer(root weaver.Instance) (*Server, error) {
	// Setup the services.
	cfpSvc, err := weaver.Get[cfpsvc.T](root)
	if err != nil {
		return nil, err
	}

	// Find out where we're running.
	// Set ENV_PLATFORM (default to local if not set; use env var if set;
	// otherwise detect GCP, which overrides env).
	var env = os.Getenv("ENV_PLATFORM")
	// Only override from env variable if set + valid env
	if env == "" || !stringinSlice(validEnvs, env) {
		fmt.Println("env platform is either empty or invalid")
		env = "local"
	}
	// Autodetect GCP
	addrs, err := net.LookupHost("metadata.google.internal.")
	if err == nil && len(addrs) >= 0 {
		root.Logger().Debug("Detected Google metadata server, setting ENV_PLATFORM to GCP.", "address", addrs)
		env = "gcp"
	}
	root.Logger().Debug("ENV_PLATFORM", "platform", env)
	platform := platformDetails{}
	platform.setPlatformDetails(strings.ToLower(env))
	hostname, err := os.Hostname()
	if err != nil {
		root.Logger().Debug(`cannot get hostname for frontend: using "unknown"`)
		hostname = "unknown"
	}

	// Create the server.
	s := &Server{
		root:     root,
		platform: platform,
		hostname: hostname,
		cfpSvc:   cfpSvc,
	}

	r := http.NewServeMux()

	// Helper that adds a handler with HTTP metric instrumentation.
	instrument := func(label string, fn func(http.ResponseWriter, *http.Request), methods []string) http.Handler {
		allowed := map[string]struct{}{}
		for _, method := range methods {
			allowed[method] = struct{}{}
		}
		handler := func(w http.ResponseWriter, r *http.Request) {
			if _, ok := allowed[r.Method]; len(allowed) > 0 && !ok {
				msg := fmt.Sprintf("method %q not allowed", r.Method)
				http.Error(w, msg, http.StatusMethodNotAllowed)
				return
			}
			fn(w, r)
		}
		return weaver.InstrumentHandlerFunc(label, handler)
	}

	const get = http.MethodGet
	const post = http.MethodPost
	const head = http.MethodHead

	r.Handle("/vote", instrument("cfp", s.voteHandler, []string{get, head}))
	r.Handle("/robots.txt", instrument("robots", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") }, nil))

	// No instrumentation of /healthz
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })

	// Set handler and return.
	var handler http.Handler = r
	handler = otelhttp.NewHandler(handler, "http") // add tracing
	s.handler = handler

	return s, nil
}

func (s *Server) Run(localAddr string) error {
	lis, err := s.root.Listener("boutique", weaver.ListenerOptions{LocalAddress: localAddr})
	if err != nil {
		return err
	}
	s.root.Logger().Debug("Frontend available", "addr", lis)
	return http.Serve(lis, s.handler)
}

func stringinSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
