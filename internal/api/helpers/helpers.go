package helpers

import (
	// "bufio"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"

	"github.com/fancar/tmp_xm/internal/config"
)

// IPaddrCheker checks if the client's ip is allowed
// via remote ip-api
func IPaddrCheker(ctx context.Context) error {
	cfg := config.C.CountryCheck
	if !cfg.Enabled {
		return nil
	}

	// getting ip from http.Request
	p, _ := peer.FromContext(ctx)
	pp := strings.Split(p.Addr.String(), ":")
	if len(pp) != 2 {
		return grpc.Errorf(codes.Internal, "unable to get ip from context")
	}

	// add ip to url
	url, err := buildURL(cfg.UrlTmpl, pp[0])
	if err != nil {
		return grpc.Errorf(codes.Internal, "unable to buil url")
	}

	// getting country from remote ip->location service ...
	resp, err := http.Get(url)
	if err != nil {
		return grpc.Errorf(codes.Internal, "unable to check ip")
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{
		"url":      url,
		"response": resp.Status,
	}).Debug("IPaddrCheker")

	if resp.StatusCode != http.StatusOK {
		return grpc.Errorf(codes.Internal,
			fmt.Sprintf("can't check your ip: %s", resp.Status))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return grpc.Errorf(codes.Internal, err.Error())
	}

	country := string(b)

	if strings.ToLower(cfg.CountryAllowed) != strings.ToLower(country) {
		return grpc.Errorf(codes.PermissionDenied,
			fmt.Sprintf("denied for requests from %s", country))
	}

	return nil
}

func buildURL(tmpl, ip string) (string, error) {

	t, err := template.New("url").Parse(tmpl)
	if err != nil {
		return "", err
	}

	data := struct {
		IPaddress string
	}{ip}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
