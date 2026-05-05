package core

import (
	"bufio"
	"io"
	"os"
)

// Scanner processes an io.Reader line-by-line and checks for matches.
type Scanner struct {
	Matcher Matcher
	After   int
	Before  int
}

// Scan reads from an io.Reader and sends matches to the results channel.
func (s *Scanner) Scan(reader io.Reader, fileName string, results chan<- MatchResult) error {
	scanner := bufio.NewScanner(reader)
	lineNumber := 1

	// Circular buffer for Before context
	var beforeBuf []MatchResult
	if s.Before > 0 {
		beforeBuf = make([]MatchResult, 0, s.Before)
	}
	bufIdx := 0

	afterCount := 0
	lastEmittedLine := 0 // Tracks last emitted line to prevent overlapping

	for scanner.Scan() {
		line := scanner.Text()
		isMatch := s.Matcher.Match(line)

		if isMatch {
			// 1. Emit Before context if applicable
			if s.Before > 0 {
				numItems := len(beforeBuf)
				startIdx := 0
				if numItems == s.Before { // Buffer is full, starts at bufIdx
					startIdx = bufIdx
				}
				for i := 0; i < numItems; i++ {
					item := beforeBuf[(startIdx+i)%s.Before]
					if item.LineNumber > lastEmittedLine {
						results <- item
						lastEmittedLine = item.LineNumber
					}
				}
			}

			// 2. Emit actual match
			mr := MatchResult{
				FileName:   fileName,
				LineNumber: lineNumber,
				Line:       line,
				IsMatch:    true,
			}
			if mr.LineNumber > lastEmittedLine {
				results <- mr
				lastEmittedLine = mr.LineNumber
			}

			// 3. Reset after counter
			afterCount = s.After
		} else {
			// Not a match
			// 1. Emit After context if we are counting down
			if afterCount > 0 {
				mr := MatchResult{
					FileName:   fileName,
					LineNumber: lineNumber,
					Line:       line,
					IsMatch:    false,
				}
				if mr.LineNumber > lastEmittedLine {
					results <- mr
					lastEmittedLine = mr.LineNumber
				}
				afterCount--
			}

			// 2. Store in before circular buffer
			if s.Before > 0 {
				mr := MatchResult{
					FileName:   fileName,
					LineNumber: lineNumber,
					Line:       line,
					IsMatch:    false,
				}
				if len(beforeBuf) < s.Before {
					beforeBuf = append(beforeBuf, mr)
				} else {
					beforeBuf[bufIdx] = mr
				}
				bufIdx = (bufIdx + 1) % s.Before
			}
		}
		lineNumber++
	}
	return scanner.Err()
}

// ScanFile opens a file and scans it.
func (s *Scanner) ScanFile(filePath string, results chan<- MatchResult) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return s.Scan(file, filePath, results)
}
