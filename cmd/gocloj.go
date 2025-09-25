package main

import (
	"flag"
	"fmt"
	//"gocloj/data/hashmap"
	"gocloj/gocloj"
	"gocloj/lib"
	"gocloj/log"
	"gocloj/runtime"
	"io"
	"os"
	"time"
)

var mainLogger = log.Get("main")

func makeEnv() *runtime.Env {
	env := runtime.NewEnv()

	// TODO: which is better?
	// env.AddLib(lib.Core)
	// or:
	lib.AddCore(env)
	lib.AddMath(env)

	return env
}

func makeParser(r io.Reader, file string) gocloj.AtomIterator {
	tz := gocloj.NewTokenizer(r, file)
	return gocloj.NewParser(tz)
}

func run(r io.Reader, file string) {
	env := makeEnv()
	p := makeParser(r, file)

	i := 0
	for p.Next() {
		if p.Err() == nil {
			value := p.Value()
			mainLogger.Infof("%d: %s", i, value.String())
			i++

			res, err := env.Eval(value)
			if err != nil {
				mainLogger.Infof("-> err: %s", err)
				break
			} else {
				mainLogger.Infof("-> %v", res)
			}
		} else {
			break
		}
	}

	if err := p.Err(); err != nil {
		mainLogger.Infof("-> parse err: %s", p.Err())
	}
}

func dump() {
	// TODO
	// phm.Dump()
}

func token(r io.Reader, file string) {
	tz := gocloj.NewTokenizer(r, file)
	for tz.Next() {
		_ = tz.Value()
	}
}

func main() {
	// TODO: NewStdoutSink
	log.AddSink(log.StdoutSink())
	log.SetLevel(log.Info)

	var filename string

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	dumpCmd := flag.NewFlagSet("dump", flag.ExitOnError)
	// runCmd.StringVar(&filename, "file", "", "path to config file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s COMMAND [OPTIONS]\n\n", os.Args[0])
		fmt.Fprint(os.Stderr, "Commands:\n")
		fmt.Fprint(os.Stderr, "  run - run files\n")
		fmt.Fprint(os.Stderr, "  dump - dumps stuff\n")
		flag.PrintDefaults()
	}

	runCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s run [OPTIONS] [FILENAME]\n\n", os.Args[0])
		fmt.Fprint(os.Stderr, "Options:\n")
		runCmd.PrintDefaults()
	}

	dumpCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s dump\n", os.Args[0])
		runCmd.PrintDefaults()
	}

	if len(os.Args[1:]) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	var cmd *flag.FlagSet

	switch os.Args[1] {
	case "run":
		cmd = runCmd
		cmd.Parse(os.Args[2:])

		names := cmd.Args()

		if len(names) == 1 {
			filename = names[0]
		} else {
			fmt.Fprint(os.Stderr, "filename must be specified.\n")
			cmd.Usage()
			os.Exit(1)
		}

	case "dump":
		cmd = dumpCmd
		cmd.Parse(os.Args[2:])

		dump()
		os.Exit(0)

	default:
		flag.Usage()
		os.Exit(1)
	}

	if filename == "" {
		fmt.Fprint(os.Stderr, "filename must be specified.\n")
		cmd.Usage()
		os.Exit(1)
	}

	f, err := os.Open(filename)
	if err != nil {
		mainLogger.Fatal("unable to open file: " + err.Error())
	}

	var (
		start   time.Time
		elapsed time.Duration
	)

	const repeat = 1

	const benchToken = false
	const benchRun = false
	const doRun = true

	if benchToken {
		log.SetLevel(log.Warn)

		start = time.Now()
		for i := 0; i < repeat; i++ {
			f.Seek(0, 0)
			token(f, filename)
		}
		elapsed = time.Since(start)
		mainLogger.Warnf("token took %0.3fms", float64(elapsed)/float64(time.Millisecond))

		/*
			start = time.Now()
			for i := 0; i < repeat; i++ {
				f.Seek(0, 0)
				token2(f, filename)
			}
			elapsed = time.Since(start)
			mainLogger.Warnf("token2 took %0.3fms", float64(elapsed)/float64(time.Millisecond))
		*/
	}

	if benchRun {
		log.SetLevel(log.Warn)

		start = time.Now()
		for i := 0; i < repeat; i++ {
			f.Seek(0, 0)
			run(f, filename)
		}
		elapsed = time.Since(start)
		mainLogger.Warnf("run took %0.3fms", float64(elapsed)/float64(time.Millisecond))

		start = time.Now()
		for i := 0; i < repeat; i++ {
			f.Seek(0, 0)
			run(f, filename)
		}
		elapsed = time.Since(start)
		mainLogger.Warnf("run took %0.3fms", float64(elapsed)/float64(time.Millisecond))
	}

	if doRun {
		log.SetLevel(log.Info)
		f.Seek(0, 0)
		start = time.Now()
		run(f, filename)
		elapsed = time.Since(start)
		mainLogger.Infof("run took %0.3fms", float64(elapsed)/float64(time.Millisecond))
	}

	f.Close()
}
