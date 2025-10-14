package api

import (
	"advent2024/pkg/solver"
	"advent2024/web/middleware"
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
//	@Param			day		path		string			true	"Day, format d[0-9]*"	example(d1)
//	@Param			part	path		int				true	"Problem part"			example(1)
//	@Param			input	body		SolvePayload	true	"Solve Base64 encoded input"
//	@Success		200		{object}	SolveResult		"Result"
//	@Failure		404
//	@Router			/solve/{day}/{part}	[post]
func Solve(w http.ResponseWriter, r *http.Request) {

	logger := middleware.GetLogger(r.Context())

	day := r.PathValue("day")
	part := r.PathValue("part")

	part_converted, err := strconv.Atoi(part)

	if err != nil {
		logger.Printf("Part is not numerical: %d", part_converted)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Solver for day %s part %s not implemented\n", day, part)))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		logger.Printf("Unable to read body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to read body")))
		return
	}

	var p SolvePayload
	err = json.Unmarshal(body, &p)

	if err != nil {
		logger.Printf("Unable to unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to unmarshal request")))
		return
	}

	decoded_body, err := base64.StdEncoding.DecodeString(string(p.Input))

	if err != nil {
		logger.Printf("Unable to decode Base64")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to read body: Invalid Base64 encoding")))
		return
	}

	solver, ok := solver.New(day)

	if !ok {
		logger.Printf("Unable to find solver for day %s", day)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Solver for day %s not implemented\n", day)))
		return
	}

	err = solver.Init(strings.NewReader(string(decoded_body)))

	if err != nil {
		logger.Printf("Unable to intialize Solver for day %s", day)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to initialize Solver for day %s\n", day)))
		return
	}

	result, err := solver.Solve(part_converted)

	if err != nil {
		logger.Printf("Unable to solve problem for day %s part %s", day, part)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unable to solve for day %s\n", day)))
		return
	}

	b, err := json.Marshal(SolveResult{Output: result})
	if err != nil {
		logger.Printf("Unable to marshal result")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
