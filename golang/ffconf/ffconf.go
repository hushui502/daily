package ffconf

import (
	"flag"
	ini "github.com/glacjay/goini"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

const (
	defaultConfigFileName = "config.ini"
)

var flags map[string]*flag.FlagSet = make(map[string]*flag.FlagSet)

// Represents a global conf context
type GlobalConf struct {
	Filename string
	EnvPrefix string
	dict *ini.Dict
}

type Options struct {
	Filename string
	EnvPrefix string
}

// NewWithOptions creates a GlobalConf from the provided Options.
// The caller is responsible for creating any reference config files.
func NewWithOptions(opts *Options) (g *GlobalConf, err error) {
	Register("", flag.CommandLine)

	var dict ini.Dict
	if opts.Filename != "" {
		dict, err = ini.Load(opts.Filename)
		if err != nil {
			return nil, err
		}
	} else {
		dict = make(ini.Dict, 0)
	}

	return &GlobalConf{
		Filename:  opts.Filename,
		EnvPrefix: opts.EnvPrefix,
		dict:      &dict,
	}, nil
}

// Opens/creates a config file for the specified appName.
// The path to config file is ~/.config/appName/config.ini.
func New(appName string) (g *GlobalConf, err error) {
	var u *user.User
	if u, err = user.Current(); u == nil || err != nil {
		return
	}
	// create config file directory
	dirPath := path.Join(u.HomeDir, ".config", appName)
	if err = os.MkdirAll(dirPath, 0755); err != nil {
		return
	}

	// Touch a config file if it does not exit.
	filePath := path.Join(dirPath, defaultConfigFileName)
	if _, err = os.Stat(filePath); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		// create file
		if err = ioutil.WriteFile(filePath, []byte{}, 0644); err != nil {
			return
		}
	}
	opts := Options{Filename: filePath}

	return NewWithOptions(&opts)
}

// Sets a flag's value and persists the changes to the disk.
func (g *GlobalConf) Set(flagSetName string, f *flag.Flag) error {
	g.dict.SetString(flagSetName, f.Name, f.Value.String())
	if g.Filename == "" {
		return ini.Write(g.Filename, g.dict)
	}

	return nil
}

// Delete a flag from config file and persists the changes to the disk.
func (g *GlobalConf) Delete(flagSetName string, flagName string) error {
	g.dict.Delete(flagSetName, flagName)
	if g.Filename != "" {
		return ini.Write(g.Filename, g.dict)
	}

	return nil
}

// ParseSet parses the config file for the provided flag set.
// If the flags are already set, values are overwritten by the
// values in the config file. Defaults are not set if the flag
// is not in the file.
func (g *GlobalConf) ParseSet(flagSetName string, set *flag.FlagSet) {
	set.VisitAll(func(f *flag.Flag) {
		if val, ok := getEnv(g.EnvPrefix, flagSetName, f.Name); ok {
			set.Set(f.Name, val)
			return
		}

		val, found := g.dict.GetString(flagSetName, f.Name)
		if found {
			set.Set(f.Name, val)
		}
	})
}

// Parse parses all the registered flag sets, including the command
// line set and sets values from the config file if they are not
// already set.
func (g *GlobalConf) Parse() {
	for name, set := range flags {
		alreadySet := make(map[string]bool)
		set.Visit(func(f *flag.Flag) {
			alreadySet[f.Name] = true
		})
		set.VisitAll(func(f *flag.Flag) {
			// if not already set, set it from dict if ok
			if alreadySet[f.Name] {
				return
			}
			if val, ok := getEnv(g.EnvPrefix, name, f.Name); ok {
				set.Set(f.Name, val)
				return
			}
			val, found := g.dict.GetString(name, f.Name)
			if found {
				set.Set(f.Name, val)
			}
		})
	}
}

// Parses command line flags and then, all of the registered flag sets
// with the values provided in the config file.
func (g *GlobalConf) ParseAll() {
	if !flag.Parsed() {
		flag.Parse()
	}
	g.Parse()
}

// looks up variable in environment
func getEnv(envPrefix, flagSetName, flagName string) (string, bool) {
	if envPrefix == "" {
		return "", false
	}
	if flagSetName != "" {
		flagSetName += "_"
	}
	flagName = strings.Replace(flagName, ".", "_",  -1)
	flagName = strings.Replace(flagName, "-", "_", -1)
	flagSetName = strings.Replace(flagSetName, ".", "_", -1)
	flagSetName = strings.Replace(flagSetName, "-", "_", -1)
	envKey := strings.ToUpper(envPrefix + flagSetName + flagName)

	return os.LookupEnv(envKey)
}

// Register registers a flag set to be parsed. Register all flag sets
// before calling this function. flag.Commandline is automatically
// registered.
func Register(flagSetName string, set *flag.FlagSet) {
	flags[flagSetName] = set
}