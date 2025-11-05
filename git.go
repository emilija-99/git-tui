package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func Run(ctx context.Context, dir string, args ...string) (string, string, error) {
	fmt.Printf("run: %v\n %v\n %v\n", ctx, dir, args)
	c := exec.CommandContext(ctx, "git", args...)
	// current dir
	c.Dir = dir

	var out, errBytes bytes.Buffer
	c.Stdout, c.Stderr = &out, &errBytes
	fmt.Printf("output: %v\n", c)

	err := c.Run()
	fmt.Printf("err: %s \n", err)

	return out.String(), errBytes.String(), err
}

func Status(ctx context.Context, dir string) ([]string, error) {
	out, _, error := Run(ctx, dir, "status", "--short")
	if error != nil {
		return nil, error
	}

	// status, _, error := Run(ctx, dir, "status", "--branch")
	if error != nil {
		return nil, error
	}

	// wrapLines := "status --short: " + out + "\nstatus --branch: " + status

	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i, l := range lines {
		lines[i] = strings.ReplaceAll(l, "??", "- ")
	}
	// fmt.Printf("%s", lines)
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	return lines, nil
}
func TimeoutCtx(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}

func Add(ctx context.Context, dir, path string) error {
	_, _, err := Run(ctx, dir, "add", "--", path)
	return err
}

func Unstage(ctx context.Context, dir, path string) error {
	_, _, err := Run(ctx, dir, "restore", "--staged", "--", path)
	return err
}

func Commit(ctx context.Context, dir, msg string) error {
	_, _, err := Run(ctx, dir, "commit", "-m", msg)
	return err
}

func Push(ctx context.Context, dir string) (string, error) {
	out, _, err := Run(ctx, dir, "push")
	return out, err
}

func Pull(ctx context.Context, dir string) (string, error) {
	out, _, err := Run(ctx, dir, "pull", "--ff-only")
	return out, err
}

func Diff(ctx context.Context, dir, path string, staged bool) (string, error) {
	args := []string{"diff"}
	if staged {
		args = []string{"diff", "--cached"}
	}
	args = append(args, "--", path)
	out, _, err := Run(ctx, dir, args...)
	return out, err
}
