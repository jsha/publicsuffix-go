package publicsuffix

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/idna"
)

type pslTestCase struct {
	input  string
	output string
	error  bool
}

func TestPsl(t *testing.T) {
	f, err := os.Open("../fixtures/tests.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	testCases := []pslTestCase{}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == "":
			break
		case strings.HasPrefix(line, "//"):
			break
		default:
			xy := strings.Split(line, " ")
			tc := pslTestCase{}
			tc.input = xy[0]
			if xy[1] == "null" {
				tc.error = true
			} else {
				tc.error = false
				tc.output = xy[1]
			}
			testCases = append(testCases, tc)
		}
	}

	for _, testCase := range testCases {
		input := testCase.input
		if strings.Contains(testCase.input, "xn--") {
			input, _ = idna.ToUnicode(input)
		}

		got, err := Domain(input)
		if strings.Contains(testCase.input, "xn--") {
			got, _ = idna.ToASCII(got)
		}

		if testCase.error && err == nil {
			t.Errorf("PSL(%v) should have returned error, got: %v", testCase.input, got)
			continue
		}
		if !testCase.error && err != nil {
			t.Errorf("PSL(%v) returned error: %v", testCase.input, err)
			continue
		}
		if testCase.output != got {
			t.Errorf("PSL(%v) = %v, want %v", testCase.input, got, testCase.output)
			continue
		}
	}
}
