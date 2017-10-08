package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"

	"code.cloudfoundry.org/cli/plugin"
)

type joke struct {
	Id     string `json:"id"`
	Text   string `json:"joke"`
	Status int    `json:"status"`
}

// TellMeAJokeCmd struct for the plugin
type TellMeAJokeCmd struct {
	ui            terminal.UI
	jokeServerURL string
}

// GetMetadata shows metadata for the plugin
func (c *TellMeAJokeCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "tell-me-a-joke-plugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 23,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "tell-me-a-joke",
				HelpText: "Ties to make you laugh",
				UsageDetails: plugin.Usage{
					Usage: "cf tell-me-a-joke",
				},
			},
		},
	}
}

func main() {
	serverURL := "https://icanhazdadjoke.com/"
	plugin.Start(NewTellMeAJokeCmd(serverURL))
}

func NewTellMeAJokeCmd(jokeServerURL string) *TellMeAJokeCmd {
	return &TellMeAJokeCmd{jokeServerURL: jokeServerURL}
}

// Run will be executed when cf tell-me-a-joke gets invoked
func (c *TellMeAJokeCmd) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] != "tell-me-a-joke" {
		return
	}

	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	c.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

	client := &http.Client{}
	req, err := http.NewRequest("GET", c.jokeServerURL, nil)
	req.Header.Add("User-Agent", "https://github.com/jtuchscherer/tell-me-joke-plugin")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		c.ui.Failed("Sad Day - Cannot connect to the Joke Server")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.ui.Failed(fmt.Sprintf("Sad Day - Could not read joke from the Joke Server - %s", err.Error()))
		return
	}

	j := joke{}

	err = json.Unmarshal(body, &j)
	if err != nil {
		c.ui.Failed(fmt.Sprintf("Sad Day - Could not read joke from the Joke Server - %s", err.Error()))
		return
	}

	c.ui.Say(j.Text)

}
