package myssh

import (
	"fmt"
	"github.com/cnlubo/promptx"
	"github.com/pkg/errors"
)

func serverLogin(server *ServerConfig) error {
	if server == nil {
		return errors.New("Bad Server")
	}
	// utils.Clear()
	err := server.Terminal()
	if err != nil {
		return fmt.Errorf("failed to Login server: %v", err)
	}
	return nil
}

func InteractiveSetContext(env *Environment) error {

	cfg := &promptx.SelectConfig{
		ActiveTpl:    `»  {{ .Name | cyan }}`,
		InactiveTpl:  `  {{ .Name | white }}`,
		SelectPrompt: "Context",
		SelectedTpl:  `{{ "» " | green }}{{ .Name | green }}`,
		DisPlaySize:  9,
		DetailsTpl: `
--------- Context ----------
{{ "Name:" | faint }} {{ .Name | faint }}
{{ "ClusterConfig:" | faint }} {{ .ClusterConfig | faint }}
{{ "SSHConfig:" | faint }} {{ .SSHConfig | faint }}`,
	}

	s := &promptx.Select{
		Items:  Main.Contexts,
		Config: cfg,
	}
	idx := s.Run()
	err := SetContext(Main.Contexts[idx].Name,env)
	if err != nil {
		return errors.Wrapf(err, "set configfile failed")
	}
	return nil
}
