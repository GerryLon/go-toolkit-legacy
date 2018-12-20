package argparser

import (
	"github.com/GerryLon/go-toolkit/common"
	"strings"
)

// 1: ^[a-zA-Z0-9]$
// 2: ^([a-zA-Z0-9]{2} | \-[a-zA-Z0-9])$ => ^((-|[a-zA-Z0-9])[a-zA-Z0-9])$
// 3: ^((-[a-zA-Z][a-zA-Z0-9]) | [a-zA-Z0-9]{3} )$
// 4: ^(([a-zA-Z0-9]{4}) | (-[a-zA-Z][a-zA-Z0-9]{2}) | (--[a-zA-Z][a-zA-Z0-9]))$
// 5: ^([a-zA-Z0-9]{5} | -[a-zA-Z][a-zA-Z0-9]{3} | --[a-zA-Z][a-zA-Z0-9]{2} | )$

// valid: is arg valid
// isOption: is arg an option
func isValidArg(arg string) (valid bool, isOption bool) {
	argLen := len(arg)

	if argLen >= 1 { // start without dash, must be value
		if !isDash(arg[0]) {
			return true, false
		}
	}

	switch argLen {
	case 0:
		return false, false

	case 1: // valid format: [^\-]
		return false, false // -

	case 2: // valid format: -f xx x-(filename ends with dash)
		if isDash(arg[1]) {
			return false, false // --
		} else if common.IsAlphaDigit(arg[1]) {
			return true, true // -f -4
		} else {
			return false, false // -[other symbol]
		}

	case 3: // valid format: fff, f-f f-- ff- -xx
		if common.IsAlpha(arg[1]) && common.IsAlpha(arg[2]) {
			return true, true
		} else {
			return false, false
		}

	case 4: // valid format: ffff -xxx --xx,
		if common.IsAlpha(arg[2]) && common.IsAlpha(arg[3]) { // -?xx
			if common.IsAlpha(arg[1]) || isDash(arg[1]) {
				return true, true
			}
		}

		// TODO: --w-afdsa
	default:
		return isValidArg(arg[0:4])
	}

	return false, false
}

// optionName: -f -> f, --list -> list --without-xxx -> without-xxx
// return optionName optionValue(only for --xx=yy style)
func getOptionName(arg string) (string, string) {
	valid, isOption := isValidArg(arg)
	if !valid || !isOption {
		return "", ""
	}
	// --with-https will cause error
	// return arg[strings.LastIndex(arg, "-")+1:]
	var i int
	for i, _ = range arg {
		if arg[i] != '-' {
			break
		}
	}

	optionName := arg[i:]
	var index int

	if len(optionName) == 1 { // short option
		return optionName, ""
	}

	// --xx --xx=yy
	if index = strings.IndexByte(optionName, '='); index < 0 {
		return optionName, ""
	} else {
		// TODO: if optionName[index:] == ""
		return optionName[:index], optionName[index+1:]
	}

}
