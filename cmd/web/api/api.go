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
type APIError struct {
	ErrorCode    int    `json:"errorcode"`
	ErrorMessage string `json:"errormessage"`
}

type RegisteredDay struct {
	Name string `json:"name"`
}

// Solve godoc
//
//	@Summary		Solver
//	@Description	Provides solution for the day and part based on input
//	@Tags			solver
//	@Accepts		json
//	@Produces		json
//	@Param			day					path		string			true	"Day, format d[0-9]*"	example(d1)
//	@Param			part				path		int				true	"Problem part"			example(1)
//	@Param			input				body		SolvePayload	true	"Solve Base64 encoded input"
//	@Success		200					{object}	SolveResult		"Result"
//	@Failure		400					{object}	APIError		"Bad Request"
//	@Failure		404					{object}	APIError		"Solver for the day not found"
//	@Failure		500					{object}	APIError		"Internal Server Error"
//	@Router			/solve/{day}/{part}	[post]
func Solve(w http.ResponseWriter, r *http.Request) {

	logger := middleware.GetLogger(r.Context())

	day := r.PathValue("day")
	part := r.PathValue("part")

	part_converted, err := strconv.Atoi(part)

	if err != nil {
		rc := http.StatusBadRequest
		errMsg := fmt.Sprintf("Solver for day %s part %s not implemented: part is not numerical", day, part)

		errJson, _ := json.Marshal(NewAPIError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)

		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		rc := http.StatusBadRequest
		errMsg := fmt.Sprintf("Unable to read body")

		errJson, _ := json.Marshal(NewAPIError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)

		return
	}

	var p SolvePayload
	err = json.Unmarshal(body, &p)

	if err != nil && p.Input == "" {
		rc := http.StatusBadRequest
		errMsg := fmt.Sprintf("Unable to read body: Invalid JSON")

		errJson, _ := json.Marshal(NewAPIError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)

		return
	}

	decoded_body, err := base64.StdEncoding.DecodeString(string(p.Input))

	if err != nil {
		rc := http.StatusBadRequest
		errMsg := fmt.Sprintf("Unable to read body: Invalid Base64 encoding")

		errJson, _ := json.Marshal(NewAPIError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)

		return
	}

	solver, ok := solver.New(day)

	if !ok {
		rc := http.StatusNotFound
		errMsg := fmt.Sprintf("Solver for day %s part %s not implemented: day not implemented", day, part)

		errJson, _ := json.Marshal(NewAPIError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)

		return
	}

	err = solver.Init(strings.NewReader(string(decoded_body)))

	if err != nil {
		rc := http.StatusBadRequest
		errMsg := fmt.Sprintf("Unable to intialize Solver for day %s", day)

		errJson, _ := json.Marshal(NewAPIError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)

		return
	}

	result, err := solver.Solve(part_converted)

	if err != nil {
		rc := http.StatusInternalServerError
		errMsg := fmt.Sprintf("Unable to solve for day %s part %s", day, part)

		errJson, _ := json.Marshal(NewAPIError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)

		return
	}

	b, err := json.Marshal(SolveResult{Output: result})
	if err != nil {
		rc := http.StatusInternalServerError
		errMsg := fmt.Sprintf("unable to Marshal result: %s", err)

		errJson, _ := json.Marshal(NewAPIError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)

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
//	@Success		200		{array}		RegisteredDay	"Result"
//	@Failure		500		{object}	APIError		"Internal Server Error"
//	@Router			/list	[GET]
func SolverListing(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r.Context())

	registered_keys := solver.ListRegister()

	response := make([]RegisteredDay, len(registered_keys))

	for i, key := range registered_keys {
		response[i] = RegisteredDay{Name: key}

	}

	b, err := json.Marshal(response)

	if err != nil {
		rc := http.StatusInternalServerError
		errMsg := "Unable to Marshal result"

		errJson, _ := json.Marshal(NewAPIError(rc, errMsg))

		logger.Println(errMsg)
		w.WriteHeader(rc)
		w.Write(errJson)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func NewAPIError(status int, message string) APIError {
	return APIError{ErrorCode: status, ErrorMessage: message}
}
