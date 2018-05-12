package tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/savaki/jq"
)

type Options struct {
	URL     string
	Path    []string
	Headers map[string]string
	StdOut  chan MessageCh
	Verbose bool
}

var suite *Suite

type Scan func(in []byte, s *Suite) error

func RunNew(options Options, scan Scan) error {
	APIURL = options.URL
	Headers = options.Headers
	Dump = options.Verbose
	StdOut = options.StdOut

	fl, err := FileList(options.Path)
	if err != nil {
		return fmt.Errorf("failed to read args, %s", err)
	}

	for _, f := range fl {
		if options.Verbose {
			Message("read file: %s\n", f)
		}

		b, err := ioutil.ReadFile(f)
		if err != nil {
			return fmt.Errorf("failed to read file %s, %s", f, err)
		}

		suite = &Suite{VarSet: make(Vars, 0)}
		if err = scan(b, suite); err != nil {
			return fmt.Errorf("failed to scan file %s, %s", f, err)
		}

		if err := runSuite(); err != nil {
			return fmt.Errorf("failed to run suite %s, %s", f, err)
		}
	}
	MessageExit()

	return nil
}

func runSuite() error {
	Message(string(suite.Title.Body))

	for _, t := range suite.Tests {
		if err := runTest(t); err != nil {
			Message("failed to run test, suite is break")
			Message(err.Error())
			break
		}
	}
	Message("")
	return nil
}

func runTest(test Test) error {
	var (
		actual = map[string]interface{}{}
		expect = map[string]interface{}{}
	)

	b, err := requestBodyBuilder(test)
	if err != nil {
		return err
	}

	b, err = sendRequest(b)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(b, &actual); err != nil {
		MessageError("failed to encode response, %s\nRESPONSE: %s", err.Error(), b)
		return fmt.Errorf("failed to encode response, %s", err.Error())
	}

	if test.HasExpectedResponseBody() {
		res := replaceVarsToValuesInBody(suite.VarSet, string(test.Response.Body))
		if err = json.Unmarshal([]byte(res), &expect); err != nil {
			MessageError("failed to encode response, %s\nRESPONSE: %s", err.Error(), test.Response.Body)
			return fmt.Errorf("failed to encode expected request, %s", err.Error())
		}
	}

	if ok, msg := assertRequestContains(actual, expect); !ok {
		MessageError("%s", test.Title.Body)
		Message(string(msg))
		Message("\n")
		return fmt.Errorf("failed to assert test")
	} else {
		MessageTest("%s", test.Title.Body)
	}

	for k, v := range test.ResponseVars {
		value, err := extractValueFromJson(v, b)
		if err != nil {
			MessageError("failed to extract value from json, var: %s, selector: %s", k, v)
			return fmt.Errorf("failed to extract value from json, %s", err)
		}
		suite.VarSet[k] = strings.Trim(string(value), "\"")
	}

	var pipe interface{}
	for _, c := range test.Conditions {
		pipe = b
		for i := 0; i < len(c); i++ {
			p, err := c[i].Cmd(pipe, c[i].Value)
			if err != nil {
				return fmt.Errorf("failed to resolve condition command, %s", err)
			}
			pipe = p
		}
	}

	for _, a := range test.Assertion {
		for _, selector := range a.NotEmpty {
			value, err := extractValueFromJson(selector, b)
			if err != nil {
				MessageError("failed to extract value from json, selector: %s", selector)
				return fmt.Errorf("failed to extract value from json, %s", err)
			}
			if len(strings.Trim(string(value), "\"")) < 1 {
				return fmt.Errorf("failed asserting not empty for %s\n", selector)
			}
		}
		for _, selector := range a.Empty {
			value, err := extractValueFromJson(selector, b)
			if err != nil {
				MessageError("failed to extract value from json, selector: %s", selector)
				return fmt.Errorf("failed to extract value from json, %s", err)
			}
			s := strings.Trim(string(value), "\"")
			if len(s) > 0 && s != "[]" {
				return fmt.Errorf("failed asserting empty for %s\n", selector)
			}
		}
	}

	return nil
}

func requestBodyBuilder(test Test) ([]byte, error) {
	s := map[string]interface{}{}
	if len(test.QueryVar.Body) > 0 {
		if !json.Valid(test.QueryVar.Body) {
			return nil, fmt.Errorf("json is invalid, line %d-%d, %s", test.QueryVar.StartLine, test.QueryVar.EndLine, string(test.QueryVar.Body))
		}

		res := replaceVarsToValuesInBody(suite.VarSet, string(test.QueryVar.Body))

		err := json.Unmarshal([]byte(res), &s)
		if err != nil {
			return nil, fmt.Errorf("failed to decode json, line %d-%d, %s", test.QueryVar.StartLine, test.QueryVar.EndLine, err)
		}
	}

	s["query"] = string(test.Query.Body)

	return json.Marshal(s)
}

func extractValueFromJson(selector string, data []byte) ([]byte, error) {
	op, err := jq.Parse(selector)
	if err != nil {
		return nil, err
	}
	return op.Apply(data)
}

func replaceVarsToValuesInBody(vars Vars, body string) string {
	for k, v := range vars {
		re := regexp.MustCompile(k + `(\W)`)
		body = re.ReplaceAllString(body, v+"$1")
	}
	return body
}
