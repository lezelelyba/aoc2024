package api

import (
	"advent2024/pkg/solver"
	"advent2024/web/middleware"
	"advent2024/web/weberrors"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type SolvePayload struct {
	Input string `json:"input" format:"base64" example:"MyAgIDQKNCAgIDMKMiAgIDUKMSAgIDMKMyAgIDkKMyAgIDMK"`
}

type SolveResult struct {
	Output string `json:"output"`
}

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
//	@Param		input					body		SolvePayload		true	"Solve Base64 encoded input"
//	@Success	200						{object}	SolveResult			"Result"
//	@Failure	400						{object}	weberrors.AoCError	"Bad Request"
//	@Failure	401						{object}	weberrors.AoCError	"Unathorized"
//	@Failure	404						{object}	weberrors.AoCError	"Solver for the day not found"
//	@Failure	429						{object}	weberrors.AoCError	"Request was Rate limited"
//	@Failure	500						{object}	weberrors.AoCError	"Internal Server Error"
//	@Router		/solvers/{day}/{part}	[post]
//	@Security	OAuth2AccessCode [read]
func Solve(w http.ResponseWriter, r *http.Request) {

	logger := middleware.GetLogger(r)

	day := r.PathValue("day")
	part := r.PathValue("part")

	part_converted, err := strconv.Atoi(part)

	var rc int
	var errMsg string

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

	var p SolvePayload
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

	solver, ok := solver.New(day)

	rc = http.StatusNotFound
	errMsg = fmt.Sprintf("Solver for day %s part %s not implemented: day not implemented", day, part)
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	err = solver.Init(strings.NewReader(string(decoded_body)))

	rc = http.StatusBadRequest
	errMsg = fmt.Sprintf("Unable to intialize Solver for day %s", day)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	result, err := solver.Solve(part_converted)

	rc = http.StatusInternalServerError
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
//	@Success		200			{array}		solver.RegistryItemPublic	"Result"
//	@Failure		401			{object}	weberrors.AoCError						"Unathorized"
//	@Failure		429			{object}	weberrors.AoCError			"Request was Rate limited"
//	@Failure		500			{object}	weberrors.AoCError			"Internal Server Error"
//	@Router			/solvers	[GET]
//	@Security		OAuth2AccessCode [read]
func SolverListing(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)

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
