package route

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stonear/api_gateway/model"
)

func New(log *logrus.Logger, routes []model.Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, r := range routes {
		target, err := url.JoinPath(r.TargetService, r.TargetPath)
		if err != nil {
			log.Errorln(err)
		}
		log.Infoln("register", r.Method, r.RequestPath, "->", target)

		mux.HandleFunc(r.RequestPath, func(writer http.ResponseWriter, request *http.Request) {
			url, err := url.Parse(target)
			if err != nil {
				log.Errorln(err)
				return
			}

			// method check here
			if request.Method != r.Method {
				writer.WriteHeader(http.StatusNotFound)
				return
			}

			// middleware check here
			if r.NeedAuth {
				// do something
				// TODO
				request.Header.Add("X-User-Id", "<user id>")
			}

			// add header here
			// TODO
			request.Header.Add("User-Agent", "<user agent>")

			proxy := httputil.NewSingleHostReverseProxy(url)
			proxy.Director = func(req *http.Request) {
				req.Header = request.Header
				req.Host = url.Host
				req.Method = request.Method
				req.URL.Host = url.Host
				req.URL.Path = url.Path
				req.URL.Scheme = url.Scheme
			}
			proxy.Transport = &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(network, addr, 5*time.Second)
					if err != nil {
						return conn, err
					}
					return conn, err
				},
				DisableKeepAlives:   false,
				MaxIdleConns:        300,
				MaxIdleConnsPerHost: 300,
			}
			proxy.ServeHTTP(writer, request)
		})
	}
	return mux
}
