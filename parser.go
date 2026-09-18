package main

import "strings"

type Command struct {
	name    string
	target  string
	payload string
}

func parseCommand(text string) Command {
	if !strings.HasPrefix(text, "/") {
		return Command{name: "BROADCAST", target: "", payload: text}
	}

	text = strings.TrimPrefix(text, "/")

	parts := strings.SplitN(text, " ", 2)

	command := strings.ToUpper(parts[0])

	payload := ""
	target := ""
	if len(parts) > 1 {
		payload = parts[1]
		parts = strings.SplitN(payload, " ", 2)
		if len(parts) > 1 {
			target = parts[0]
			payload = parts[1]
		}
	}

	return Command{name: command, target: target, payload: payload}
}
