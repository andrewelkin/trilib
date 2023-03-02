package logger

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

// FileWriter implements io.Writer and can be used with as a logger output to write logs to a log folder, in files named
// according to the UTC day.
type FileWriter struct {
	mux sync.Mutex

	basePath string
	prefix   string
	suffix   string

	lastLineRaw         []byte
	lastLineCmp         []byte
	lastLineRepeatCt    int
	lastLineRepeatStart time.Time

	lastLineSkipQ bool // if true the repeating lines will be merged

	currentFile *os.File
}

// NewFileWriter creates a new file writer that writes to files in basePath
// If prefix and/or suffix are set, they will be applied to the file names.
func NewFileWriter(basePath string, prefix, suffix *string, skipRepeating bool) (*FileWriter, error) {
	var p, s string
	if prefix != nil {
		p = *prefix
	}
	if suffix != nil {
		s = *suffix
	}

	// make sure the base directory exists
	baseExists, err := exists(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to check base path: %w", err)
	}
	if !baseExists {
		if err := os.Mkdir(basePath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	return &FileWriter{basePath: basePath, prefix: p, suffix: s, lastLineSkipQ: skipRepeating}, nil
}

// Write implements io.Writer
func (fw *FileWriter) Write(p []byte) (n int, err error) {
	file, err := fw.getOrCreateFile(fw.getCurrentFileName())
	if err != nil {
		return 0, err
	}

	var line []byte
	if fw.lastLineSkipQ {
		// perform repeat check and line parsing
		if skipLine, modLine := fw.updateLastLine(p); !skipLine {
			line = modLine
		} else {
			return 0, nil
		}
	} else {
		line = p
	}
	return file.Write(line)
}

func (fw *FileWriter) updateLastLine(rawLine []byte) (skip bool, modifiedLine []byte) {
	fw.mux.Lock()
	defer fw.mux.Unlock()

	// copy rawline; skip dup tracking if not "log format"
	tmp := make([]byte, len(rawLine))
	copy(tmp, rawLine)
	spl := bytes.Split(tmp, []byte("]: "))
	if len(spl) < 2 {
		return false, rawLine
	}
	lineCmp := spl[1]

	// no last line saved: first log line handled (special case)
	if fw.lastLineRaw == nil {
		fw.lastLineCmp = lineCmp
		fw.lastLineRaw = tmp
		fw.lastLineRepeatCt = 0
		return false, rawLine
	}

	// repeating line contents: increase count and indicate to not write anything
	if bytes.Compare(lineCmp, fw.lastLineCmp) == 0 {
		if fw.lastLineRepeatCt == 0 {
			fw.lastLineRepeatStart = time.Now()
		}

		fw.lastLineRepeatCt += 1
		return true, nil
	}

	// past here, this is a new line
	// - if it as a new line after a new line, proceed to log unmodified as normal
	// - if it is a new line after a series of repeats, write a summary line for the repeat and the new line
	if fw.lastLineRepeatCt == 0 {
		fw.lastLineCmp = lineCmp
		fw.lastLineRaw = tmp
		return false, rawLine
	}

	// past here, we have to create a replacement line with the summary of the
	// repeating line (ensuring to strip out the original '\n') as well as the
	// current logged line appended at the end
	l := fmt.Sprintf(
		"%s (repeated %v times over %s)\n%s",
		string(fw.lastLineRaw[0:len(fw.lastLineRaw)-1]),
		fw.lastLineRepeatCt,
		time.Now().Sub(fw.lastLineRepeatStart),
		string(tmp),
	)
	fw.lastLineCmp = lineCmp
	fw.lastLineRaw = tmp
	fw.lastLineRepeatCt = 0
	return false, []byte(l)
}

func (fw *FileWriter) getCurrentFileName() string {
	return fmt.Sprintf("%s%s%s", fw.prefix, time.Now().UTC().Format("2006-01-02"), fw.suffix)
}

func (fw *FileWriter) getOrCreateFile(name string) (*os.File, error) {
	fw.mux.Lock()
	defer fw.mux.Unlock()

	fullPath := path.Join(fw.basePath, name)
	if fw.currentFile != nil {
		if path.Clean(fw.currentFile.Name()) == path.Clean(fullPath) {
			return fw.currentFile, nil
		}
		fw.currentFile.Close()
		fw.currentFile = nil
	}

	// check if the file already exists; if so, we append a separator later
	fileExists, err := exists(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check if file exists: %w", err)
	}

	file, err := os.OpenFile(fullPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// apply the separator if this is not a new file
	if fileExists {
		if _, err := file.WriteString(separator); err != nil {
			return nil, fmt.Errorf("failed to write separator to existing file: %w", err)
		}
	}
	fw.currentFile = file
	return file, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

const separator = "======\n"
