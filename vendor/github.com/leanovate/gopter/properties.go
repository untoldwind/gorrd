package gopter

import "testing"

// Properties is a collection of properties that should be checked in a test
type Properties struct {
	parameters *TestParameters
	props      map[string]Prop
	propNames  []string
}

// NewProperties create new Properties with given test parameters.
// If parameters is nil default test parameters will be used
func NewProperties(parameters *TestParameters) *Properties {
	if parameters == nil {
		parameters = DefaultTestParameters()
	}
	return &Properties{
		parameters: parameters,
		props:      make(map[string]Prop, 0),
		propNames:  make([]string, 0),
	}
}

func (p *Properties) Property(name string, prop Prop) {
	p.propNames = append(p.propNames, name)
	p.props[name] = prop
}

func (p *Properties) Run(reporter Reporter) bool {
	success := true
	for _, propName := range p.propNames {
		prop := p.props[propName]

		result := prop.Check(p.parameters)

		reporter.ReportTestResult(propName, result)
		if !result.Passed() {
			success = false
		}
	}
	return success
}

func (p *Properties) TestingRun(t *testing.T) {
	if !p.Run(ConsoleReporter(true)) {
		t.Fail()
	}
}
