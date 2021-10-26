package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "rerun-fail-test",
		Description: "rerun fail go test with go test logfile",
		Flags: []cli.Flag{
			&cli.Int64Flag{Name: "retry-times", Usage: "how many times to retry", Value: 3},
		},
		Action: run,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func run(c *cli.Context) error {
	retryTimes := c.Int("retry-times")
	args := c.Args().Slice()

	failTestCases, err := filterFailTestCase(io.TeeReader(os.Stdin, os.Stdout))
	if err != nil {
		return err
	}

	for _, testCase := range failTestCases {
		var err error
		for i := 0; i < retryTimes; i++ {
			err = runGoTest(testCase, args)
			if err == nil {
				fmt.Printf("[rerun-fail-test] %s run success at %d\n", testCase, i+1)
				break
			} else {
				fmt.Printf("[rerun-fail-test] %s run fail at %d\n", testCase, i+1)
				if i < retryTimes-1 {
					continue
				} else {
					return err
				}
			}
		}
	}

	return nil
}

func filterFailTestCase(reader io.Reader) ([]string, error) {
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, v := range regFailTestCase.FindAllStringSubmatch(string(bs), -1) {
		if len(v) == 2 {
			res = append(res, v[1])
		}
	}
	return res, nil
}

var regFailTestCase = regexp.MustCompile(`(?m)--- FAIL: (.*?) \(`)

func runGoTest(testCase string, args []string) error {
	args = append(append([]string{"test"}, args...), "-test.run", testCase)
	cmd := exec.Command("go", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
