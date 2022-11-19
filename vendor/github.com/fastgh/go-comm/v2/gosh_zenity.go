package comm

import (
	"context"
	"fmt"
	"time"

	"github.com/ncruces/zenity"
	"mvdan.cc/sh/v3/interp"
)

func GoshExecHandler(killTimeout time.Duration) interp.ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		cmd := args[0]

		var cmdArgs []string
		if len(args) > 1 {
			cmdArgs = args[1:]
		} else {
			cmdArgs = []string{}
		}

		switch cmd {
		case "zenity":
			return ExecZenity(ctx, interp.HandlerCtx(ctx), cmdArgs)
		default:
			return interp.DefaultExecHandler(killTimeout)(ctx, args)
		}
	}
}

func ExecZenity(ctx context.Context, hc interp.HandlerContext, args []string) error {
	subCmd := args[0]

	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	switch subCmd {
	case "--error":
		return ZenityError(ctx, hc, args)
	case "--info":
		return ZenityInfo(ctx, hc, args)
	case "--warning":
		return ZenityWarning(ctx, hc, args)
	case "--question":
		return ZenityWarning(ctx, hc, args)
	default:
		fmt.Fprintln(hc.Stdout, "unknown sub command")
		return nil
	}
}

func ZenityError(ctx context.Context, hc interp.HandlerContext, args []string) error {
	var text string
	if len(args) > 0 {
		text = args[0]
	} else {
		text = ""
	}

	return zenity.Error(text)
}

func ZenityInfo(ctx context.Context, hc interp.HandlerContext, args []string) error {
	var text string
	if len(args) > 0 {
		text = args[0]
	} else {
		text = ""
	}

	return zenity.Info(text)
}

func ZenityWarning(ctx context.Context, hc interp.HandlerContext, args []string) error {
	var text string
	if len(args) > 0 {
		text = args[0]
	} else {
		text = ""
	}

	return zenity.Warning(text)
}

func ZenityQuestion(ctx context.Context, hc interp.HandlerContext, args []string) error {
	var text string
	if len(args) > 0 {
		text = args[0]
	} else {
		text = ""
	}

	return zenity.Question(text)
}
