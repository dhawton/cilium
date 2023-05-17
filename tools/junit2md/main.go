// Copyright Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type TestSuites struct {
	XMLName  xml.Name    `xml:"testsuites"`
	Tests    int         `xml:"tests,attr"`
	Failures int         `xml:"failures,attr"`
	Skipped  int         `xml:"skipped,attr"`
	Suites   []Testsuite `xml:"testsuite"`
}

type TestSuiteJenkins struct {
	XMLName xml.Name `xml:"testsuite"`
	Testsuite
}

type Testsuite struct {
	Name      string     `xml:"name,attr"`
	Tests     int        `xml:"tests,attr"`
	Failures  int        `xml:"failures,attr"`
	Errors    int        `xml:"errors,attr"`
	ID        int        `xml:"id,attr"`
	Hostname  string     `xml:"hostname,attr"`
	Time      float64    `xml:"time,attr"`
	Timestamp string     `xml:"timestamp,attr"`
	Testcases []Testcase `xml:"testcase"`
}

type Testcase struct {
	Name      string   `xml:"name,attr"`
	Classname string   `xml:"classname,attr"`
	Time      float64  `xml:"time,attr"`
	SystemOut string   `xml:"system-out"`
	Failure   *Failure `xml:"failure,omitempty"`
	XMLName   xml.Name `xml:"testcase"`
	Skipped   *Skipped `xml:"skipped,omitempty"`
	Error     *Error   `xml:"error,omitempty"`
	Filename  string   `xml:"filename,omitempty"`
}

type Skipped struct{}

type Failure struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

type Error struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

type TestCaseResult struct {
	TestSuite string
	TestCase  string
	ClassName string
	Status    string
	Time      float64
}

const (
	StatusPass    = "Pass"
	StatusSkipped = "Skipped"
	StatusFail    = "Fail"
	StatusError   = "Error"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the input file.")
		return
	}

	inputFile := flag.String("i", "", "Input file path")
	outputFile := flag.String("o", "", "Output file path")
	failed := flag.Bool("f", true, "Show failed tests")
	passed := flag.Bool("p", true, "Show passed tests")
	skipped := flag.Bool("s", true, "Show skipped tests")
	errored := flag.Bool("e", true, "Show errored tests")
	flag.Parse()

	// Open the file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Extract test case results
	var testCaseResults []TestCaseResult

	// Create a struct to store the unmarshalled data
	var testsuites TestSuites
	// Read and decode the XML from the file
	err = xml.NewDecoder(file).Decode(&testsuites)
	if err == nil {
		for _, suite := range testsuites.Suites {
			if suite.Name == "" {
				continue
			}
			testSuiteStatus := StatusPass
			if len(suite.Testcases) == 0 {
				testSuiteStatus = StatusSkipped
			}
			for _, testcase := range suite.Testcases {
				testCaseStatus := status(testcase)
				switch testCaseStatus {
				case StatusPass, StatusSkipped:
				default:
					testSuiteStatus = testCaseStatus
				}
			}
			testCaseResults = addTestCase(testCaseResults, suite.Name, testSuiteStatus, suite.Time, passed, skipped, failed, errored)
		}
	} else {
		_, err := file.Seek(0, io.SeekStart)
		if err != nil {
			// try Jenkins
			fmt.Println("Error seeking:", err)
			return
		}
		// try Jenkins
		var tsj TestSuiteJenkins
		// Read and decode the XML from the file
		err2 := xml.NewDecoder(file).Decode(&tsj)
		if err2 != nil {
			// try Jenkins
			fmt.Println("Error decoding XML:", err2)
			return
		}
		testsuites = TestSuites{
			XMLName:  tsj.XMLName,
			Tests:    tsj.Tests,
			Failures: tsj.Failures,
			Suites:   []Testsuite{tsj.Testsuite},
		}
		for _, suite := range testsuites.Suites {
			if suite.Name == "" {
				continue
			}
			for _, testcase := range suite.Testcases {
				testCaseStatus := status(testcase)
				testCaseResults = addTestCase(testCaseResults, testcase.Name, testCaseStatus, testcase.Time, passed, skipped, failed, errored)
			}
		}
	}
	if len(testCaseResults) == 0 {
		return
	}

	// Sort test case results
	sort.Slice(testCaseResults, func(i, j int) bool {
		if testCaseResults[i].Status != testCaseResults[j].Status {
			switch testCaseResults[i].Status {
			case StatusError, StatusFail:
				switch testCaseResults[j].Status {
				case StatusError, StatusFail:
					return strings.Compare(testCaseResults[i].TestSuite, testCaseResults[j].TestSuite) < 0
				}
				return true
			}
		}
		return strings.Compare(testCaseResults[i].TestSuite, testCaseResults[j].TestSuite) < 0
	})

	fout := os.Stdout
	// Write the table to the output
	if outputFile != nil && len(*outputFile) != 0 {
		fout, err = os.Create(*outputFile)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		defer fout.Close()
	}

	// Generate markdown table
	_, err = io.WriteString(fout, "| Status | Package | Time (seconds) |\n")
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}
	_, err = io.WriteString(fout, "|--------|---------|----------------|\n")
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}
	for _, result := range testCaseResults {
		var statusEmoji string
		switch result.Status {
		case StatusPass:
			statusEmoji = ":heavy_check_mark:"
		case StatusSkipped:
			statusEmoji = ":white_check_mark:"
		case StatusFail:
			statusEmoji = ":x:"
		case StatusError:
			statusEmoji = ":warning:"
		}
		row := fmt.Sprintf("| %-6s | %-10s | %-14.3f |\n", statusEmoji, result.TestSuite, result.Time)

		_, err = io.WriteString(fout, row)
		if err != nil {
			fmt.Println("Error writing output file:", err)
			return
		}
	}

	if outputFile != nil && len(*outputFile) != 0 {
		fmt.Println("Markdown table saved to", *outputFile)
	}
}

func status(testcase Testcase) string {
	var status string
	switch {
	case testcase.Skipped != nil:
		status = StatusSkipped
	case testcase.Failure != nil:
		status = StatusFail
	case testcase.Error != nil:
		status = StatusError
	default:
		status = StatusPass
	}
	return status
}

func addTestCase(testCaseResults []TestCaseResult, name, status string, timeElapsed float64, passed, skipped, failed, errored *bool) []TestCaseResult {
	testCaseResult := TestCaseResult{
		TestSuite: name,
		Status:    status,
		Time:      timeElapsed,
	}
	switch status {
	case StatusPass:
		if *passed {
			testCaseResults = append(testCaseResults, testCaseResult)
		}
	case StatusSkipped:
		if *skipped {
			testCaseResults = append(testCaseResults, testCaseResult)
		}
	case StatusFail:
		if *failed {
			testCaseResults = append(testCaseResults, testCaseResult)
		}
	case StatusError:
		if *errored {
			testCaseResults = append(testCaseResults, testCaseResult)
		}
	}
	return testCaseResults
}
