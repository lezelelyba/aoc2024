package api

import (
	"advent2024/pkg/solver"
	"advent2024/web/middleware"
	"advent2024/web/weberrors"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type SolveRequest struct {
	Input string `json:"input" format:"base64" example:"MyAgIDQKNCAgIDMKMiAgIDUKMSAgIDMKMyAgIDkKMyAgIDMK"`
} //@name Request

type SolveResult struct {
	Output string `json:"output"`
} //@name Response

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
func Solve(w http.ResponseWriter, r *http.Request) {

	var rc int
	var errMsg string

	logger := middleware.GetLogger(r)
	cfg, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = "configuration error: index: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	day := r.PathValue("day")
	part := r.PathValue("part")

	part_converted, err := strconv.Atoi(part)

	rc = http.StatusBadRequest
	errMsg = fmt.Sprintf("Solver for day %s part %s not implemented: part is not numerical", day, part)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	rc = http.StatusBadRequest
	errMsg = "unable to read body"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	var p SolveRequest
	err = json.Unmarshal(body, &p)

	rc = http.StatusBadRequest
	errMsg = "unable to read body: Invalid JSON"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	decoded_body, err := base64.StdEncoding.DecodeString(string(p.Input))

	rc = http.StatusBadRequest
	errMsg = "unable to read body: Invalid Base64 encoding"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	slvr, ok := solver.NewWithCtx(day)

	rc = http.StatusNotFound
	errMsg = fmt.Sprintf("Solver for day %s part %s not implemented: day not implemented", day, part)
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), cfg.SolverTimeout)
	defer cancel()

	err = slvr.InitCtx(ctx, strings.NewReader(string(decoded_body)))

	if errors.Is(err, solver.ErrTimeout) {
		rc = http.StatusGatewayTimeout
	} else {
		rc = http.StatusBadRequest
	}

	errMsg = fmt.Sprintf("Unable to intialize Solver for day %s", day)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	result, err := slvr.SolveCtx(ctx, part_converted)

	if errors.Is(err, solver.ErrTimeout) {
		rc = http.StatusGatewayTimeout
	} else {
		rc = http.StatusInternalServerError
	}

	errMsg = fmt.Sprintf("Unable to solve for day %s part %s", day, part)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

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
func SolverListing(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)

	w.Header().Set("Content-Type", "application/json")

	registryItems := solver.ListRegistryItems()

	b, err := json.Marshal(registryItems)

	rc := http.StatusInternalServerError
	errMsg := "unable to marchal registered items solvers"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
