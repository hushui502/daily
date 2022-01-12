package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Args[1:]); err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	// Parse command line flags.
	fs := flag.NewFlagSet("kubectl-super-cp", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	} else if fs.NArg() != 2 {
		return fmt.Errorf("usage: kubectl-super-cp SOURCE DESTINATION")
	}

	// Parse source pod & remote path.
	a := strings.SplitN(fs.Arg(0), ":", 2)
	if len(a) != 2 {
		return fmt.Errorf("must specify pod in source path")
	}
	podName, remotePath := a[0], a[1]

	// Open local file to write to.
	f, err := os.Create(fs.Arg(1))
	if err != nil {
		return err
	}
	defer f.Close()

	// Continuously try to resume file copy.
	for {
		if err := copyToFile(ctx, f, podName, remotePath); err != nil {
			fmt.Fprintln(os.Stderr, "copy failed, resuming")
			continue
		}
		break
	}

	// Sync and close file.
	if err := f.Sync(); err != nil {
		return err
	} else if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func copyToFile(ctx context.Context, f *os.File, podName, remotePath string) error {
	// Determine current local file size.
	fi, err := f.Stat()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "copying from offset %d\n", fi.Size())

	// Open a "kubectl exec" to write from our last file position to our local file.
	cmd := exec.CommandContext(ctx, `kubectl`, `exec`, podName, `--`, `tail`, `-c`, fmt.Sprintf("+%d", fi.Size()), remotePath)
	cmd.Stdout = f
	return cmd.Run()
}
