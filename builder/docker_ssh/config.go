package docker_ssh

import (
    "fmt"
    "github.com/mitchellh/packer/common"
    "github.com/mitchellh/packer/packer"
)

type Config struct {
    common.PackerConfig `mapstructure:",squash"`

    ExportPath string `mapstructure:"export_path"`
    Image      string
    Pull       bool
    RunCommand []string `mapstructure:"run_command"`

    Port              int    `mapstructure:"port"`
    SSHUsername       string `mapstructure:"ssh_username"`
    SSHPassword       string `mapstructure:"ssh_password"`
    SSHPrivateKeyFile string `mapstructure:"ssh_private_key_file"`

    tpl *packer.ConfigTemplate
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
    c := new(Config)
    md, err := common.DecodeConfig(c, raws...)
    if err != nil {
        return nil, nil, err
    }

    c.tpl, err = packer.NewConfigTemplate()
    if err != nil {
        return nil, nil, err
    }

    c.tpl.UserVars = c.PackerUserVars

    // Defaults
    if len(c.RunCommand) == 0 {
        c.RunCommand = []string{
            "run",
            "-d", //"-i", "-t",
            "-v", "{{.Volumes}}",
            "{{.Image}}",
            "/sbin/init",
        }
    }

    // Default Pull if it wasn't set
    hasPull := false
    for _, k := range md.Keys {
        if k == "Pull" {
            hasPull = true
            break
        }
    }

    if !hasPull {
        c.Pull = true
    }

    // Default ssh port
    if c.Port == 0 {
        c.Port = 22
    }

    errs := common.CheckUnusedConfig(md)

    templates := map[string]*string{
        "export_path": &c.ExportPath,
        "image":       &c.Image,
    }

    for n, ptr := range templates {
        var err error
        *ptr, err = c.tpl.Process(*ptr, nil)
        if err != nil {
            errs = packer.MultiErrorAppend(
                errs, fmt.Errorf("Error processing %s: %s", n, err))
        }
    }

    if c.ExportPath == "" {
        errs = packer.MultiErrorAppend(errs,
            fmt.Errorf("export_path must be specified"))
    }

    if c.Image == "" {
        errs = packer.MultiErrorAppend(errs,
            fmt.Errorf("image must be specified"))
    }

    if c.SSHUsername == "" {
        errs = packer.MultiErrorAppend(errs,
            fmt.Errorf("ssh_username must be specified"))
    }

    if c.SSHPassword == "" && c.SSHPrivateKeyFile == "" {
        errs = packer.MultiErrorAppend(errs,
            fmt.Errorf("one of ssh_password and ssh_private_key_file must be specified"))
    }

    if c.SSHPassword != "" && c.SSHPrivateKeyFile != "" {
        errs = packer.MultiErrorAppend(errs,
            fmt.Errorf("only one of ssh_password and ssh_private_key_file must be specified"))
    }

    if errs != nil && len(errs.Errors) > 0 {
        return nil, nil, errs
    }

    return c, nil, nil
}
