package main

import "testing"

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedName    string
		expectedTarget  string
		expectedPayload string
	}{
		{
			name:            "Normale Broadcast",
			input:           "ciao a tutti",
			expectedName:    "BROADCAST",
			expectedTarget:  "",
			expectedPayload: "ciao a tutti",
		},
		{
			name:            "Join Canale",
			input:           "/join #general",
			expectedName:    "JOIN",
			expectedTarget:  "#general",
			expectedPayload: "",
		},
		{
			name:            "Privmsg a utente",
			input:           "/privmsg mario ciao come stai?",
			expectedName:    "PRIVMSG",
			expectedTarget:  "mario",
			expectedPayload: "ciao come stai?",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := parseCommand(tc.input)

			if result.name != tc.expectedName {
				t.Errorf("Comando errato. Atteso: %s, Ottenuto: %s", tc.expectedName, result.name)
			}
			if result.target != tc.expectedTarget {
				t.Errorf("Target errato. Atteso: %s, Ottenuto: %s", tc.expectedTarget, result.target)
			}
			if result.payload != tc.expectedPayload {
				t.Errorf("Payload errato. Atteso: '%s', Ottenuto: '%s'", tc.expectedPayload, result.payload)
			}
		})
	}
}
