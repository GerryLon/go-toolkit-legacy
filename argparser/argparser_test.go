package argparser

import (
	"testing"
)

func Test_isValidArg(t *testing.T) {
	// var arg string = "-f"
	var isValid, isOption bool
	var optName, optValue string

	if isValid, isOption = isValidArg("-f"); !isValid || !isOption {
		t.Error("isValidArg() err")
	}

	if optName, optValue = getOptionName("-f"); optName != "f" {
		t.Error("getOptionName() err")
	}

	if optName, optValue = getOptionName("--wi"); optName != "wi" {
		t.Error("getOptionName() err")
	}

	if optName, optValue = getOptionName("--with-https"); optName != "with-https" {
		t.Error("getOptionName() err")
	}

	if optName, optValue = getOptionName("--config=conf"); optName != "config" {
		t.Error("getOptionName() err")
	}
	if optValue != "conf" {
		t.Error("getOptionName() err")
	}

}

func Test_Parse(t *testing.T) {
	//Parse()
	// opts := Parse(strings.Split("-f file -d --config=conf -e -x abc --with-https", " "))
	// fmt.Printf("%+v\n", opts)
	// for k, v := range opts {
	// 	fmt.Printf("%s=%s\n", k, v)
	// }
}
