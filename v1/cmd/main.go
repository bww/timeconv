package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

var errUnsupported = errors.New("Unsupported format")

type format []string

func (f format) Default() string {
	if len(f) > 0 {
		return f[0]
	} else {
		return ""
	}
}

func (f format) Matches(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	for _, e := range f {
		if e == v {
			return true
		}
	}
	return false
}

var (
	formatNano    = format{"nano", "nanos"}
	formatMicro   = format{"micro", "micros"}
	formatMilli   = format{"milli", "millis"}
	formatUnix    = format{"unix", "sec", "secs"}
	formatRFC3339 = format{"rfc3339"}
)

type Options struct {
	FromFormat string `long:"from" value-name:"FORMAT" description:"Input timestamp format"`
	ToFormat   string `long:"to" value-name:"FORMAT" description:"Output timestamp format"`
	Debug      bool   `long:"debug" description:"Enable debugging mode."`
	Verbose    []bool `short:"v" long:"verbose" description:"Be more verbose; specify repeatedly for greater verbosity."`
}

func (o Options) From(text string) (time.Time, error) {
	switch {
	case formatUnix.Matches(o.FromFormat):
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(v, 0), nil
	case formatMilli.Matches(o.FromFormat):
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.UnixMilli(v), nil
	case formatMicro.Matches(o.FromFormat):
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.UnixMicro(v), nil
	case formatNano.Matches(o.FromFormat):
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(0, v), nil
	default:
		return time.Time{}, errUnsupported
	}
}

func (o Options) To(t time.Time) (string, error) {
	switch {
	case formatRFC3339.Matches(o.ToFormat):
		return t.Format(time.RFC3339), nil
	default:
		return "", errUnsupported
	}
}

func parseArgs(dst interface{}, args []string) ([]string, error) {
	return flags.NewParser(dst, flags.HelpFlag|flags.PrintErrors|flags.PassAfterNonOption).ParseArgs(args)
}

func main() {
	err := exec()
	if err != nil {
		fmt.Println("* * *", err)
		os.Exit(1)
	}
}

func exec() error {
	opts := Options{
		FromFormat: formatUnix.Default(),
		ToFormat:   formatRFC3339.Default(),
	}

	args, err := parseArgs(&opts, os.Args[1:])
	if err != nil {
		return err
	}

	if len(args) == 0 {
		d, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		args = strings.Split(string(d), "\n")
	}

	for _, e := range args {
		t, err := opts.From(e)
		if err != nil {
			return fmt.Errorf("Input: %s as %v: %w", e, opts.FromFormat, err)
		}
		s, err := opts.To(t)
		if err != nil {
			return fmt.Errorf("Output: %s as %v: %w", t, opts.ToFormat, err)
		}
		fmt.Println(s)
	}

	return nil
}
