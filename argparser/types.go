package argparser

// map[string]Value =  argparser.Parse()
type argParser interface {
	Parse() []map[string]option
}

// option struct
type option struct {
	name      string
	usage     string
	required  bool
	shortName string
	value     interface{}
	valueType string
}

type ArgParser struct {
	args map[string]*option
}

func (ap *ArgParser) Parse() {

}
