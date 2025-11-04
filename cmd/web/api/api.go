// API http Handlers
package api

import (
	"advent2024/pkg/solver"
	"advent2024/web/config"
	"advent2024/web/middleware"
	"advent2024/web/weberrors"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// API Request
type SolveRequest struct {
	Input string `json:"input" format:"base64" example:"MyAgIDQKNCAgIDMKMiAgIDUKMSAgIDMKMyAgIDkKMyAgIDMK"`
} //@name Request

// API Response
type SolveResult struct {
	Output string `json:"output"`
} //@name Response

type CodeExchangeRequest struct {
	Provider string `json:"provider"`
	Code     string `json:"code"`
} //@name Code Exchange Request

type JWTToken struct {
	Token string `json:"access_token"`
} //@name Token Response

type InfoResponse struct {
	Version        string `json:"version"`
	Authentication string `json:"authentication"`
} //@name Info Response

// Solve godoc
//
//	@Summary		Solver
//	@Description	Provides solution for the day and part based on input
//	@Tags			solver
//	@Accepts		json
//	@Produces		json
//	@Security
//	@Param		day						path		string				true	"Day, format d[0-9]*"	example(d1)
//	@Param		part					path		int					true	"Problem part"			example(1)
//	@Param		input					body		SolveRequest		true	"Solve Base64 encoded input"
//	@Param		Authorization			header		string				true	"Bearer format, prefix with Bearer"
//	@Success	200						{object}	SolveResult			"Result"
//	@Failure	400						{object}	weberrors.AoCError	"Bad Request"
//	@Failure	401						{object}	weberrors.AoCError	"Unathorized"
//	@Failure	404						{object}	weberrors.AoCError	"Solver for the day not found"
//	@Failure	429						{object}	weberrors.AoCError	"Request was Rate limited"
//	@Failure	500						{object}	weberrors.AoCError	"Internal Server Error"
//	@Failure	504						{object}	weberrors.AoCError	"Request took too long to compute"
//	@Router		/solvers/{day}/{part}	[post]
//	@Security	OAuth2AccessCode [read]
//
// Handles solve API endpoint
func Solve(w http.ResponseWriter, r *http.Request) {

	var rc int
	var errMsg string

	// get logger and config
	logger := middleware.GetLogger(r)
	cfg, ok := middleware.GetConfig(r)

	// unable to get config
	rc = http.StatusInternalServerError
	errMsg = "configuration error: index: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	// prepare response headers, always JSON
	w.Header().Set("Content-Type", "application/json")

	// get part and day from request URL
	day := r.PathValue("day")
	part := r.PathValue("part")
	part_converted, err := strconv.Atoi(part)

	rc = http.StatusBadRequest
	errMsg = fmt.Sprintf("Solver for day %s part %s not implemented: part is not numerical", day, part)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	// read request
	// limit the size of read response
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	rc = http.StatusBadRequest
	errMsg = "unable to read body"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	// unmarshall request body
	var p SolveRequest
	err = json.Unmarshal(body, &p)

	rc = http.StatusBadRequest
	errMsg = "unable to read body: Invalid JSON"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	// decode the base64 encoded request
	decoded_body, err := base64.StdEncoding.DecodeString(string(p.Input))

	rc = http.StatusBadRequest
	errMsg = "unable to read body: Invalid Base64 encoding"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	// get solver
	slvr, ok := solver.NewWithCtx(day)

	rc = http.StatusNotFound
	errMsg = fmt.Sprintf("Solver for day %s part %s not implemented: day not implemented", day, part)
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	// cancel request after deadline
	ctx, cancel := context.WithTimeout(r.Context(), cfg.SolverTimeout)
	defer cancel()

	// init
	err = slvr.InitCtx(ctx, strings.NewReader(string(decoded_body)))

	// initialization took too long or input error?
	if errors.Is(err, solver.ErrTimeout) {
		rc = http.StatusGatewayTimeout
	} else {
		rc = http.StatusBadRequest
	}

	errMsg = fmt.Sprintf("Unable to intialize Solver for day %s", day)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	// try to solve
	result, err := slvr.SolveCtx(ctx, part_converted)

	// solution took too long or solver error?
	if errors.Is(err, solver.ErrTimeout) {
		rc = http.StatusGatewayTimeout
	} else {
		rc = http.StatusInternalServerError
	}

	errMsg = fmt.Sprintf("Unable to solve for day %s part %s", day, part)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	// prepare response
	b, err := json.Marshal(SolveResult{Output: result})
	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("unable to Marshal result: %s", err)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// SolverListing godoc
//
//	@Summary		Solve List
//	@Description	Lists days which the solver can solve
//	@Tags			solverList
//	@Accepts		json
//	@Produces		json
//	@Param			Authorization	header		string						true	"Bearer format, prefix with Bearer"
//	@Success		200				{array}		solver.RegistryItemPublic	"Result"
//	@Failure		401				{object}	weberrors.AoCError			"Unathorized"
//	@Failure		429				{object}	weberrors.AoCError			"Request was Rate limited"
//	@Failure		500				{object}	weberrors.AoCError			"Internal Server Error"
//	@Router			/solvers														[GET]
//	@Security		OAuth2AccessCode [read]
//
// Handles solver listing API endpoint
func SolverListing(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)

	// prepare response headers
	w.Header().Set("Content-Type", "application/json")

	// gets registry of ctx supporting solvers
	registryItems := solver.ListRegistryItemsWithCtx()

	// prepare response body
	b, err := json.Marshal(registryItems)

	rc := http.StatusInternalServerError
	errMsg := "unable to marchal registered items solvers"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// OAuthCodeExchange godoc
//
//	@Summary		Exchanged OAuth code for a JWT token
//	@Description	Exchanges OAuth code for a JWT token
//	@Tags			Oauth_code_exchange
//	@Accepts		json
//	@Produces		json
//	@Success		200				{object}		JWTToken	"JWT Token"
//	@Failure		401				{object}	weberrors.AoCError			"Unathorized"
//	@Failure		500				{object}	weberrors.AoCError			"Internal Server Error"
//	@Router			/public/auth_token														[POST]
//
// Handles solver listing API endpoint
func OAuthCodeExchange(w http.ResponseWriter, r *http.Request) {
	// prepare response headers
	w.Header().Set("Content-Type", "application/json")

	var rc int
	var errMsg string

	logger := middleware.GetLogger(r)
	config, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = "configuration error: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	// read request
	// limit the size of read response
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	rc = http.StatusBadRequest
	errMsg = "unable to read body"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	// unmarshall request body
	var codeExchangeRequest CodeExchangeRequest
	err = json.Unmarshal(body, &codeExchangeRequest)

	rc = http.StatusBadRequest
	errMsg = "unable to read body: Invalid JSON"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	// get provider name
	providerName := codeExchangeRequest.Provider

	// get provider config
	provider, exists := config.OAuthProviders[providerName]

	rc = http.StatusBadRequest
	errMsg = fmt.Sprintf("unknown Oauth provider %s", providerName)
	if weberrors.HandleError(w, logger, weberrors.OkToError(exists), rc, errMsg) != nil {
		return
	}

	switch provider.Name() {
	case "github":

		// code is missing
		rc = http.StatusBadRequest
		errMsg = fmt.Sprintf("unable to exchange code for token with %s: code is missing", provider.Name())
		ok := codeExchangeRequest.Code != ""
		if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
			return
		}

		// exchange code
		token, err := exchangeCodeForToken(&provider, codeExchangeRequest.Code)

		// error => local error
		rc = http.StatusInternalServerError
		errMsg = fmt.Sprintf("unable to exchange code for token with %s: %v", provider.Name(), err)
		if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
			return
		}

		// error nil, but error response from provider
		if token.IsError() {
			rc = http.StatusBadRequest
			errMsg = fmt.Sprintf("unable to exchange code for token with %s: %v", provider.Name(), err)
			if weberrors.HandleError(w, logger, token, rc, errMsg) != nil {
				return
			}
		}

		// prepare response for client
		w.Header().Set("Content-Type", "application/json")

		// generate JWT Token
		jwtToken, err := middleware.GenerateJWT(provider.Name(), []byte(config.JWTSecret), config.JWTTokenValidity)

		// unable to generate token
		rc = http.StatusInternalServerError
		errMsg = "unable to create jwt token"
		if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
			return
		}

		// write response to client
		// prepare response
		b, err := json.Marshal(JWTToken{Token: jwtToken})
		rc = http.StatusInternalServerError
		errMsg = fmt.Sprintf("unable to Marshal token: %s", err)
		if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

// Exchanges code with OAuth provider for a token
func exchangeCodeForToken(provider *config.OAuthProvider, code string) (middleware.OAuthResponse, error) {
	// no provider
	if provider == nil {
		return nil, fmt.Errorf("unable to find empty provider")
	}

	// know providres
	switch (*provider).Name() {
	case "github":
		// extract required information from client request
		data := url.Values{}
		data.Set("client_id", (*provider).ClientID())
		data.Set("client_secret", (*provider).ClientSecret())
		data.Set("code", code)
		// TODO: send redirect_uri to github

		// create request to provider
		req, err := http.NewRequest(
			"POST",
			(*provider).TokenURL(),
			strings.NewReader(data.Encode()))

		if err != nil {
			return nil, fmt.Errorf("unable to create OAuth request")
		}

		// set headers and make the request
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)

		// err from client
		if err != nil {
			return nil, fmt.Errorf("unable to exchange code for token with %s: %v", (*provider).Name(), err)
		}

		// work only with non-nil response
		defer resp.Body.Close()

		// only process OK responses
		if resp.StatusCode != http.StatusOK {
			// if there is a response include part of it
			limited := io.LimitReader(resp.Body, 80)
			data, err := io.ReadAll(limited)
			if err != nil {
				data = []byte("")
			}
			return nil, fmt.Errorf("unable to exchange code for token with %s: %s", (*provider).Name(), data)
		}

		// decode token
		var token middleware.OAuthGithubReply

		err = json.NewDecoder(resp.Body).Decode(&token)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal %s response", (*provider).Name())
		}

		return token, nil
	// unknown provider
	default:
		return nil, fmt.Errorf("unable to find provider %s", (*provider).Name())
	}
}

func Info(w http.ResponseWriter, r *http.Request) {
	var rc int
	var errMsg string

	logger := middleware.GetLogger(r)
	cfg, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = "configuration error: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	// prepare response headers
	w.Header().Set("Content-Type", "application/json")

	var info InfoResponse

	info.Version = cfg.Version

	if cfg.OAuth {
		info.Authentication = "oauth"
	} else {
		info.Authentication = "none"
	}

	// prepare response body
	b, err := json.Marshal(info)

	rc = http.StatusInternalServerError
	errMsg = "unable to mashal info"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
