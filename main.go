package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ali.ghanem/generic-fizzbuzz/zhttp"
	"github.com/ali.ghanem/generic-fizzbuzz/zlogging"
)

type API struct {
	Counter    *Counter
	converters chan *ValueConverter
}

func main() {
	// read config
	logPath := os.Getenv("API_LOG_PATH")
	if len(logPath) == 0 {
		log.Fatal("api log path not defined")
	}

	apiPort := os.Getenv("API_PORT")
	if len(apiPort) == 0 {
		log.Fatal("api port not defined")
	}

	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// add multiple log outputs
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	var a = API{
		Counter: &Counter{
			RWMutex: sync.RWMutex{},
			Metrics: make(map[string]*Metric),
		},
		converters: make(chan *ValueConverter),
	}
	go a.consumeConverters()

	// setup endpoints
	http.HandleFunc("/generate", a.generate)
	http.HandleFunc("/statistics", a.statistics)
	err = http.ListenAndServe(fmt.Sprintf(":%s", apiPort), nil)
	if err != nil {
		zlogging.WriteError("stopped", err)
	}

	close(a.converters)
}

func (a *API) consumeConverters() {
	for converter := range a.converters {
		a.Counter.Inc(converter)
	}
}

func (a *API) statistics(w http.ResponseWriter, _ *http.Request) {
	max := a.Counter.GetMaxHits()
	if max == nil {
		zlogging.WriteError("max number of hits not found", nil)
		a.writeError(w, http.StatusNotFound, "not_found", errors.New("max number of hits not found"))
		return
	}

	a.write(w, max)
}

func (a *API) generate(w http.ResponseWriter, r *http.Request) {
	converter, errs := readAndValidate(r)
	if len(errs) > 0 {
		zlogging.WriteError("invalid request", flattenErrors(errs))
		err := zhttp.Write(w, http.StatusBadRequest, nil, errs...)
		if err != nil {
			zlogging.WriteError("cannot write response", err)
		}
		return
	}

	results, err := converter.Generate()
	if err != nil {
		zlogging.WriteError("cannot generate results", err)
		a.writeError(w, http.StatusInternalServerError, "internal_error", fmt.Errorf("cannot generate results: %w", err))
		return
	}

	a.write(w, results)

	// send converter to the chan to save it metric
	a.converters <- converter
}

func readAndValidate(r *http.Request) (*ValueConverter, []error) {
	var errs []error
	// check if params are set
	int1, err := zhttp.ReadIntParam(r, "int1")
	if err != nil {
		errs = append(errs, fmt.Errorf("cannot read int1: %w", err))
	}

	int2, err := zhttp.ReadIntParam(r, "int2")
	if err != nil {
		errs = append(errs, fmt.Errorf("cannot read int2: %w", err))
	}

	limit, err := zhttp.ReadIntParam(r, "limit")
	if err != nil {
		errs = append(errs, fmt.Errorf("cannot read limit: %w", err))
	}

	str1, err := zhttp.ReadStringParam(r, "str1")
	if err != nil {
		errs = append(errs, fmt.Errorf("cannot read str1: %w", err))
	}

	str2, err := zhttp.ReadStringParam(r, "str2")
	if err != nil {
		errs = append(errs, fmt.Errorf("cannot read str2: %w", err))
	}

	if len(errs) > 0 {
		return nil, errs
	}

	// validate values
	if int1 <= 0 {
		errs = append(errs, errors.New("int1 must be greater than 0"))
	}

	if int2 <= 0 {
		errs = append(errs, errors.New("int2 must be greater than 0"))
	}

	if limit <= 1 {
		errs = append(errs, errors.New("limit must be greater than 1"))
	}

	if len(str1) == 0 {
		errs = append(errs, errors.New("str1 must not be empty"))
	}

	if len(str2) == 0 {
		errs = append(errs, errors.New("str2 must not be empty"))
	}

	if len(errs) > 0 {
		return nil, errs
	}

	converter := ValueConverter{
		Value1: int1,
		Word1:  str1,
		Value2: int2,
		Word2:  str2,
		Limit:  limit,
	}
	return &converter, nil
}

func (a *API) write(w http.ResponseWriter, data interface{}) {
	err := zhttp.Write(w, http.StatusOK, data)
	if err != nil {
		zlogging.WriteError("cannot write http response", err)
	}
}

func (a *API) writeError(w http.ResponseWriter, status int, code string, err error) {
	var detail string
	if err != nil {
		detail = err.Error()
	}
	err = zhttp.WriteError(w, status, code, detail)
	if err != nil {
		zlogging.WriteError("cannot write response", err)
	}
}

func flattenErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var multierr strings.Builder
	for _, err := range errs {
		multierr.WriteString(fmt.Sprintf("%s\n", err.Error()))
	}

	return errors.New(multierr.String())
}
