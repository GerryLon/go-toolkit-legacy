package argparser

import (
	"fmt"
	"github.com/GerryLon/go-toolkit/common"
	"os"
	"strings"
)

// command line arguments parser

var osArgs []string // argv
var osArgsStr string
var programName string   // argv[0]
var programArgs []string // argv[1:]
var ap *ArgParser

// TODO, ignore mainFile
var mainFile string // xx.exe -l -s config.conf, mainFile is config.conf

// initialize ap
func init() {
	osArgs = os.Args
	osArgsStr = strings.Join(osArgs, " ")
	programName = osArgs[0]
	programArgs = osArgs[1:]

	if ap == nil {
		ap = new(ArgParser)
	}
}

func formatErr() {
	fmt.Printf("format err: %s\n", osArgsStr)
}

// func (ap *ArgParser) config(opts []option) map[string]*option {
// 	if ap.args == nil {
// 		ap.args = make(map[string]*option, 0)
// 	}
//
// 	for _, opt := range opts {
// 		if common.IsEmptyString(opt.name) {
// 			panic("option name can not be empty!")
// 		}
// 		ap.args[opt.name] = &opt
// 	}
//
// 	return ap.args
// }

func exit(code int) {
	os.Exit(code)
}

// TODO xx.exe config.conf xx.conf yy.conf
// xx.exe -l -s --config config.conf
// TODO: xx.exe config.conf -sl(-l -s)
func parse() map[string]string {
	argc := len(programArgs)
	// 先指定长度, 避免后面不必要扩容
	programArgsParsed := make(map[string]string, 0)
	optValue := ""
	optName := ""
	var valid, isOption bool
	// [-s "", -l "", --xx "xx"]

	for i, arg := range programArgs {
		// fmt.Println(i, arg, programArgsParsed, optValue)
		if valid, isOption = isValidArg(arg); !valid {
			formatErr()
			goto End
		}

		// -x, --xx, --xx=yy
		if isOption {
			optName, optValue = getOptionName(arg)

			// --xx=yy
			if optValue != "" {
				programArgsParsed[optName] = optValue
				continue
			}

			if i < argc-1 {
				if valid, isOption = isValidArg(programArgs[i+1]); !valid {
					formatErr()
					goto End
				}

				if isOption { // -x -y
					programArgsParsed[optName] = ""
				} else { // -x f
					programArgsParsed[optName] = programArgs[i+1]
				}

			} else { // last option
				programArgsParsed[optName] = ""
			}
		} else {
			if i == 0 { // xx.exe f
				formatErr()
				goto End
			} else {
				if valid, isOption = isValidArg(programArgs[i-1]); !valid || !isOption {
					formatErr()
					goto End
				}
			}
		}
	}

End:
	return programArgsParsed
}

func Parse() {
	optsParsed := parse()

	for k1, v1 := range optsParsed {
		valueType := ap.args[k1].valueType
		switch valueType {
		case "string":
			ap.args[k1].value = v1
		case "bool":
			ap.args[k1].value = true
		}
	}
}

func usage() {
	u := fmt.Sprintf("Usage of %s:\n", programName)

	for k, v := range ap.args {
		u += fmt.Sprintf("-%s %s\n", k, v.valueType)
	}
	u += "\n"
	fmt.Println(u)
}

func Option(name string, usage string) *option {
	if ap.args == nil {
		ap.args = make(map[string]*option, 0)
	}

	if common.IsEmptyString(name) {
		fmt.Println("name can not be null")
		return nil
	}

	if _, ok := ap.args[name]; !ok {
		opt := option{
			name:  name,
			usage: usage,
		}
		ap.args[name] = &opt
		return &opt
	}

	return ap.args[name]
}

// default value
func (opt *option) Default(defaultValue interface{}) *option {
	name := opt.name
	if _, ok := ap.args[name]; !ok {
		fmt.Printf("error: option %s is not exist\n", name)
		return opt
	}

	switch t := defaultValue.(type) {
	case string:
		ap.args[name].value = defaultValue.(string)
	case int:
		ap.args[name].value = defaultValue.(int)
	case bool:
		ap.args[name].value = defaultValue.(bool)
	default:
		panic(fmt.Sprintf("invalid value type %v of option %s", t, name))
	}

	return opt
}

// short name
func (opt *option) Short(shortName string) *option {
	name := opt.name
	if _, ok := ap.args[name]; !ok {
		fmt.Printf("error: option %s is not exist\n", name)
		return opt
	}

	if common.IsEmptyString(shortName) {
		fmt.Printf("error: shortName %s is empty\n", shortName)
		return opt
	}
	ap.args[name].shortName = shortName
	return opt
}

// set option to be required
func (opt *option) Required() *option {
	name := opt.name
	if _, ok := ap.args[name]; !ok {
		// fmt.Printf("error: option %s is not exist\n", name)
		usage()
		return opt
	}

	ap.args[name].required = true
	return opt
}

func (opt *option) Value() interface{} {
	name := opt.name
	if _, ok := ap.args[name]; !ok {
		usage()
		return nil
	}

	if ap.args[name].value == nil && ap.args[name].required {
		usage()
		return nil
	}

	switch ap.args[name].valueType {
	case "string":
		return ap.args[name].value.(string)
	case "bool":
		return ap.args[name].value.(bool)
	}

	return nil
}

// option value to string
func (opt *option) String() *option {

	name := opt.name
	if _, ok := ap.args[name]; !ok {
		fmt.Printf("error: option %s is not exist\n", name)
		return opt
	}

	ap.args[name].valueType = "string"
	return opt
}

// option value to bool
func (opt *option) Bool() *option {

	name := opt.name
	if _, ok := ap.args[name]; !ok {
		// fmt.Printf("error: option %s is not exist\n", name)
		return opt
	}

	ap.args[name].valueType = "bool"
	// return ap.args[name].value.(bool)
	// for bool type, just option name indicate true
	return opt
}
