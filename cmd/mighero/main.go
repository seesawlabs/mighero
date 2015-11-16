package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	_ "github.com/mattes/migrate/driver/bash"
	_ "github.com/mattes/migrate/driver/cassandra"
	_ "github.com/mattes/migrate/driver/mysql"
	_ "github.com/mattes/migrate/driver/postgres"
	_ "github.com/mattes/migrate/driver/sqlite3"
	"github.com/mattes/migrate/file"
	"github.com/mattes/migrate/migrate"
	"github.com/mattes/migrate/migrate/direction"
	pipep "github.com/mattes/migrate/pipe"
	"gopkg.in/yaml.v2"
)

const (
	cmdCreate  = "create"
	cmdMigrate = "migrate"
	cmdHelp    = "help"
	cmdVersion = "version"
	cmdUp      = "up"
	cmdDown    = "down"
	cmdReset   = "reset"
	cmdRedo    = "redo"
	cmdGoto    = "goto"
)

//Configuration struce is a map of the config file.
type Configuration struct {
	DB struct {
		IP           string `yaml:"ip"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Name         string `yaml:"name"`
		MigrationDir string `yaml:"migration_dir"`
		Driver       string `yaml:"driver"`
	} `yaml:"db"`
}

var path struct {
	ToDefault string
	ToEnv     string
}

func main() {
	flag.Usage = func() {
		helpCmd()
	}
	flag.StringVar(&path.ToDefault, "def", "env/default.yml", "the default configuration file.")
	flag.StringVar(&path.ToEnv, "env", "env/local.yml", "the environment configuration file.")
	flag.Parse()

	fmt.Println("Def config -", path.ToDefault)
	fmt.Println("Env config -", path.ToEnv)

	cmd := flag.Arg(0)

	c, err := initConfig(path.ToDefault, path.ToEnv)
	if err != nil {
		log.Fatal(err)
	}

	migrationDir := c.DB.MigrationDir

	url := fmt.Sprintf("%s://%s:%s@tcp(%s)/%s", c.DB.Driver, c.DB.User, c.DB.Password, c.DB.IP, c.DB.Name)

	switch cmd {
	case cmdCreate:

		name := flag.Arg(1)
		if name == "" {
			fmt.Println("Please specify name.")
			os.Exit(1)
		}

		migrationFile, err := migrate.Create(url, migrationDir, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Version %v migration files created in %v:\n", migrationFile.Version, migrationDir)
		fmt.Println(migrationFile.UpFile.FileName)
		fmt.Println(migrationFile.DownFile.FileName)
		////////// CREATE END ///////////////
	case cmdMigrate:
		relativeN := flag.Arg(1)
		relativeNInt, err := strconv.Atoi(relativeN)
		if err != nil {
			fmt.Println("Unable to parse param <n>.")
			os.Exit(1)
		}
		timerStart = time.Now()
		pipe := pipep.New()
		go migrate.Migrate(pipe, url, migrationDir, relativeNInt)
		ok := writePipe(pipe)
		printTimer()
		if !ok {
			os.Exit(1)
		}
		////////// MIGRATE END ///////////////
	case cmdUp:
		timerStart = time.Now()
		pipe := pipep.New()
		go migrate.Up(pipe, url, migrationDir)
		ok := writePipe(pipe)
		printTimer()
		if !ok {
			os.Exit(1)
		}
		////////// UP END ///////////////
	case cmdDown:
		timerStart = time.Now()
		pipe := pipep.New()
		go migrate.Down(pipe, url, migrationDir)
		ok := writePipe(pipe)
		printTimer()
		if !ok {
			os.Exit(1)
		}

		////////// DOWN END ///////////////
	case cmdVersion:
		version, err := migrate.Version(url, migrationDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(version)
		////////// VERSION END ///////////////

	case cmdReset:
		timerStart = time.Now()
		pipe := pipep.New()
		go migrate.Reset(pipe, url, migrationDir)
		ok := writePipe(pipe)
		printTimer()
		if !ok {
			os.Exit(1)
		}
		////////// RESET END ///////////////

	case cmdRedo:
		timerStart = time.Now()
		pipe := pipep.New()
		go migrate.Redo(pipe, url, migrationDir)
		ok := writePipe(pipe)
		printTimer()
		if !ok {
			os.Exit(1)
		}
		////////// REDO END ///////////////

	case cmdGoto:
		toVersion := flag.Arg(1)
		toVersionInt, err := strconv.Atoi(toVersion)
		if err != nil || toVersionInt < 0 {
			fmt.Println("Unable to parse param <v>.")
			os.Exit(1)
		}

		currentVersion, err := migrate.Version(url, migrationDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		relativeNInt := toVersionInt - int(currentVersion)

		timerStart = time.Now()
		pipe := pipep.New()
		go migrate.Migrate(pipe, url, migrationDir, relativeNInt)
		ok := writePipe(pipe)
		printTimer()
		if !ok {
			os.Exit(1)
		}
		////////// GOTO  END ///////////////

	default:
		fallthrough
	case cmdHelp:
		helpCmd()
	}
}

func writePipe(pipe chan interface{}) (ok bool) {
	okFlag := true
	if pipe != nil {
		for {
			select {
			case item, more := <-pipe:
				if !more {
					return okFlag
				}

				switch item.(type) {
				case string:
					fmt.Println(item.(string))

				case error:
					c := color.New(color.FgRed)
					c.Println(item.(error).Error())
					okFlag = false

				case file.File:
					f := item.(file.File)
					c := color.New(color.FgBlue)
					if f.Direction == direction.Up {
						c.Print(">")
					} else if f.Direction == direction.Down {
						c.Print("<")
					}
					fmt.Printf(" %s\n", f.FileName)

				default:
					text := fmt.Sprint(item)
					fmt.Println(text)
				}
			}
		}
	}
	return okFlag
}

var timerStart time.Time

func printTimer() {
	diff := time.Now().Sub(timerStart).Seconds()
	if diff > 60 {
		fmt.Printf("\n%.4f minutes\n", diff/60)
	} else {
		fmt.Printf("\n%.4f seconds\n", diff)
	}
}

func helpCmd() {
	os.Stderr.WriteString(
		`usage: mighero [-def=<path> -env=<path>] <command> [<args>]
Commands:
   create <name>  Create a new migration
   up             Apply all -up- migrations
   down           Apply all -down- migrations
   reset          Down followed by Up
   redo           Roll back most recent migration, then apply it again
   version        Show current migration version
   migrate <n>    Apply migrations -n|+n
   help           Show this help
'-path' defaults to the subdirectory env of current working directory.
`)
}

func initConfig(defaultConfigPath, envConfigPath string) (*Configuration, error) {
	def, err := ioutil.ReadFile(defaultConfigPath)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(def)

	env, err := ioutil.ReadFile(envConfigPath)
	if err != nil {
		return nil, err
	}

	buf.Write(env)

	cMap := map[string]*Configuration{}

	if err = yaml.Unmarshal(buf.Bytes(), cMap); err != nil {
		return nil, err
	}

	return cMap["config"], nil
}
