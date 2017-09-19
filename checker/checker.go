package checker

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

// CheckerSubCommands are the subcommands for the checker cli.
var CheckerSubCommands = []cli.Command{
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

// CheckerFlags are the flags of the checker command.
var CheckerFlags = []cli.Flag{
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

// CheckerCommand is the command with live and readiness subcommands, and the root flags.
var CheckerCommand = cli.Command{
	Name:        "checker",
	Usage:       "sub-commands check to see if a status-enabled service is live or ready",
	Subcommands: CheckerSubCommands,
}
