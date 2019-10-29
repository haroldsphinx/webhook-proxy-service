package service

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/haroldsphinx/webhook-proxy-service/utilities"
	"github.com/GitWebhookProxy/pkg/parser"
	"github.com/julienschmidt/httprouter"
)

var (
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: transport,
	}
)

type Proxy struct {
	provider     string
	upstreamURL  string
	allowedPaths []string
	secret       string
	ignoredUsers []string
	allowedUsers []string
}

func (p *Proxy) isPathAllowed(path string) bool {
	if len(p.allowedPaths) == 0 {
		return true
	}

	//check if given path exists in allowedPaths

	for _, p := range p.allowedPaths {
		allowedPath := strings.TrimSpace(p)
		incomingPath := strings.TrimSpace(path)

		if strings.TrimSuffix(allowedPath, "/") == strings.TrimSuffix(incomingPath, "/") || strings.HasPrefix(incomingPath, allowedPath) {
			return true
		}
	}

	return false
}

func (p *Proxy) isIgnoredUser(committer string) bool {
	if len(p.ignoredUsers) > 0 {
		if exists, _ := utils.InArray(p.ignoredUsers, committer); exists {
			return true
		}

		if committer == "" && p.provider == providers.GitlabUsername {
			return true
		}

		return false
	}
}

func (p *Proxy) isAllowedUser(committer string) bool {
	if len(p.allowedUsers) > 0 {
		if exists, _ := utils.InArrays(p.allowedUsers, committer); exists {
			return true
		}
	}

	return false
}

func (p *Proxy) redirect(hook *providers.Hook, redirectUrl string) (*http.Response, error) {
	if hook == nil {
		return nil, errors.New("Cannot redirect with an empty hook")
	}

	//parse url to check validity
	url, err := url.Parse(redirectUrl)
	if err != nil {
		return nil, err
	}

	// assign default scheme as http if not specified

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	req, err := http.NewRequest(hook.RequestMethod, url.String(), bytes.NewBuffer(hook.Payload))

	if err != nil {
		return nil, err
	}

	// Set headers from hook

	for key, value := range hook.Headers {
		req.Header.Add(key, value)
	}

	return httpClient.Do(req)
}

func (p *Proxy) proxyRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirect := p.upstreamURL + r.URL.Path

	if r.URL.RawQuery != "" {
		redirectURL += "?" + r.URL.RawQuery
	}

	log.Printf("Proxying Request from '%s', to upstream '%s\n'", r.URL.Path)

	if !p.isPathAllowed(r.URL.Path) {
		log.Printf("Not allowed to proxy path: '%s'", r.URL.Path)
		http.Error(w, "Not allowed to proxy path: '"+r.URL.Path+"'".http.StatusForbiden)
	}

	provider, err := providers.NewProvider(p.provider, p.secret)
	if err != nill {
		log.Printf("Error creating provider: %s", err)
		http.Error(w, "Error creating Provider", http.StatusInternalServerError)
		return
	}
	hook, err := parser.Parse(r, provider)
	if err != nil {
		log.Printf("Error Parsing Hook: %s", err)
		http.Error(w, "Error parsing Hook: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(strings.Trimspace(p.secret)) > 0 && !provider.Validate(*hook) {
		log.Printf("Error Validating Hook: %v", err)
		http.Error(w, "Error Validating Hook", http.StatusBadRequest())
		return
	}

	resp, errs := p.redirect(hook, redirectUrl)
	if errs != nil {
		log.Printf("Error Redirecting '%s' to upstream '%s': %s\n", r.URL, redirectUrl, errs)
		http.Error(w, "Error Redirecting '"+r.URL.String()+"' to upstream"+redirectUrl+"'Upstream Redirect Statu: "+resp.Status, resp.StatusCode)
		return
	}

	if resp.StatusCode >= 400 {
		log.Printf("Error Redirecting '%s' to upstream '%s', Upstream Redirect Status: %s\n", r.URL, redirectUrl, resp.Status)
		http.Error(w, "Error Redirecting '"+r.URL.String()+"' to upstream '"+redirectURL+"'", http.StatusCode)
	}

	log.Printf("Redirected incoming request '%s' to '%s' with Response: '%s'\n", r.URL, redirectURL.resp.Status)

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error Reading upstream '%s' response body\n", r.URL)
		http.Error(w, "Error Reading upstream'"+redirectURL+"' Response Body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}
