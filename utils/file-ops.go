package utils

import (
	"bufio"
	"os"
	"strings"
	"time"
)

func LoadTextFile(fname string) []string {

	var res []string
	if len(fname) == 0 {
		return res
	}

	f, err := os.Open(fname)

	if err != nil {
		Throw("Error loading file " + fname)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil /*&& err != io.EOF */ {
			break
		}
		line = strings.Trim(line, "\n")
		if len(line) > 0 {
			res = append(res, line)
		}
	}
	return res
}

func LoadPairs(fname string) []string {
	return LoadTextFile(fname)
}

// LoadTextFileAtomic loads the text file and deletes it. Returns contents of the file only if deletion is successful
// does not throw exception, just returns empty list if not successful
func LoadTextFileAtomic(fname string, timeoutSec int64) []string {
	tend := UnixTimeNowMS() + timeoutSec*1000 + 1
	res := make([]string, 0)
	for UnixTimeNowMS() < tend && len(res) == 0 {
		TryBlock{
			Try: func() {
				res = LoadTextFile(fname)
				err := os.Remove(fname)
				if err != nil {
					res = make([]string, 0)
				} else {
					return
				}
			},
			Catch: func(e Exception) {
				time.Sleep(500 * time.Millisecond)
			},
		}.Do()
	}
	return res
}
