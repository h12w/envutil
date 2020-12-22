package envutil

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/multierr"
)

type Reader struct {
	prefix string
	errs   []error
}

type EnvError struct {
	Name string
	Err  error
}

func NewReader(prefix string) *Reader {
	return &Reader{prefix: prefix}
}

func (r *Reader) Env(name string) (string, bool) {
	return os.LookupEnv(r.prefix + name)
}

func (r *Reader) String(name string, defaultValue ...string) (ret string) {
	if r.MoreThanOneError(name, len(defaultValue)) {
		return
	}
	value, ok := r.Env(name)
	if !ok {
		if r.NoDefaultError(name, len(defaultValue)) {
			return
		}
		return defaultValue[0]
	}
	return value
}

func (r *Reader) Bool(name string, defaultValue ...bool) (ret bool) {
	if r.MoreThanOneError(name, len(defaultValue)) {
		return
	}
	value, ok := r.Env(name)
	if !ok {
		if r.NoDefaultError(name, len(defaultValue)) {
			return
		}
		return defaultValue[0]
	}

	ret, err := strconv.ParseBool(value)
	if err != nil {
		r.AddError(name, err)
		return
	}
	return
}

func (r *Reader) Int(name string, defaultValue ...int) (ret int) {
	if r.MoreThanOneError(name, len(defaultValue)) {
		return
	}
	value, ok := r.Env(name)
	if !ok {
		if r.NoDefaultError(name, len(defaultValue)) {
			return
		}
		return defaultValue[0]
	}

	ret, err := strconv.Atoi(value)
	if err != nil {
		r.AddError(name, err)
		return
	}
	return
}

func (r *Reader) Duration(name string, defaultValue ...time.Duration) (ret time.Duration) {
	if r.MoreThanOneError(name, len(defaultValue)) {
		return
	}
	value, ok := r.Env(name)
	if !ok {
		if r.NoDefaultError(name, len(defaultValue)) {
			return
		}
		return defaultValue[0]
	}

	ret, err := time.ParseDuration(value)
	if err != nil {
		r.AddError(name, err)
		return
	}
	return
}

func (r *Reader) MoreThanOneError(name string, numDefaults int) bool {
	if numDefaults > 1 {
		r.AddErrorf(name, "more than one default value")
		return true
	}
	return false
}

func (r *Reader) NoDefaultError(name string, numDefaults int) bool {
	if numDefaults == 0 {
		r.AddErrorf(name, "no value set")
		return true
	}
	return false
}

func (r *Reader) AddError(name string, err error) {
	if err != nil {
		r.errs = append(r.errs, EnvError{Name: r.prefix + name, Err: err})
	}
}

func (r *Reader) AddErrorf(name string, format string, args ...interface{}) {
	r.AddError(name, fmt.Errorf(format, args...))
}

func (r *Reader) Err() error {
	return multierr.Combine(r.errs...)
}

func (e EnvError) Error() string {
	return fmt.Sprintf("%s: %v", e.Name, e.Err)
}
