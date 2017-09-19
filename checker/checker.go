package main

import (
	"context"
	"time"

	"github.com/paralin/grpc-liveness/statussvc"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

var cliArgs struct {
	// Endpoint is the endpoint to contact
	Endpoint string
	// TimeoutSeconds is the maximum amount of time to wait for the call to complete.
	TimeoutSeconds int
	// FailFast indicates that we won't wait for the maximum timeout to fail.
	FailFast bool
}

// performChecks does the actual checking
func performChecks(isLiveness bool) error {
	ctx, ctxCancel := context.WithTimeout(
		context.Background(),
		time.Duration(cliArgs.TimeoutSeconds)*time.Second,
	)
	defer ctxCancel()

	conn, err := grpc.DialContext(
		ctx,
		cliArgs.Endpoint,
		grpc.WithInsecure(),
		grpc.FailOnNonTempDialError(true),
	)
	if err != nil {
		return err
	}

	callOptions := []grpc.CallOption{
		grpc.FailFast(cliArgs.FailFast),
	}

	statusClient := statussvc.NewStatusServiceClient(conn)
	if isLiveness {
		_, err = statusClient.GetLiveness(
			ctx,
			&statussvc.GetLivenessRequest{},
			callOptions...,
		)
	} else {
		_, err = statusClient.GetReadiness(
			ctx,
			&statussvc.GetReadinessRequest{},
			callOptions...,
		)
	}
	return err
}

func main() {
	app := cli.NewApp()
	app.Name = "checker"
	app.Usage = "Checks to see if a GRPC service is online or not, by calling a status service."
	app.HideVersion = true
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "ready",
			Usage: "checks to see if the GRPC service is ready",
			Action: func(c *cli.Context) error {
				return performChecks(false)
			},
		},
		cli.Command{
			Name:  "live",
			Usage: "checks to see if the GRPC service is alive",
			Action: func(c *cli.Context) error {
				return performChecks(true)
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "endpoint",
			Usage:       "the grpc endpoint to contact",
			Destination: &cliArgs.Endpoint,
			Value:       "localhost:5000",
		},
		cli.BoolTFlag{
			Name:        "fail-fast",
			Usage:       "If set, don't wait for the timeout to fail, but fail quickly.",
			Destination: &cliArgs.FailFast,
		},
		cli.IntFlag{
			Name:        "timeout-seconds",
			Usage:       "sets the seconds to wait for the call to go through",
			Destination: &cliArgs.TimeoutSeconds,
			Value:       3,
		},
	}
	app.RunAndExitOnError()
}
