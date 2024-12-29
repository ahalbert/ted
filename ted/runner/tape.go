package runner

import (
	"errors"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/edsrzf/mmap-go"
)

type Tape interface {
	Split(seperator string)
	Scan() bool
	Text() string
	Seek(int, int) (int, error)
	Prev() bool
	Next() bool
}

var ErrEof = errors.New("EOF")
var ErrBof = errors.New("BOF")

type StringTape struct {
	input     string
	groups    []string
	offset    int
	maxOffset int
	seperator string
}

func NewStringTape(in string) *StringTape {
	return &StringTape{input: in,
		offset:    0,
		maxOffset: 0,
		seperator: "\n",
		groups:    strings.Split(in, "\n"),
	}
}

func (ss *StringTape) Text() string {
	return ss.groups[ss.offset]
}

func (ss *StringTape) Split(seperator string) {
	ss.seperator = seperator
	ss.groups = strings.Split(ss.input, ss.seperator)
}

func (ss *StringTape) Prev() bool {
	if ss.offset <= 0 {
		return false
	}
	ss.offset--
	return true
}

func (ss *StringTape) Next() bool {
	if ss.offset >= len(ss.groups) {
		return false
	}
	ss.offset++
	ss.maxOffset = max(ss.offset, ss.maxOffset)
	return true
}

func (ss *StringTape) Scan() bool {
	ss.maxOffset++
	ss.offset = ss.maxOffset
	if ss.offset == len(ss.groups) {
		return false
	}
	return true
}

func (ss *StringTape) Seek(offset int, whence int) (int, error) {
	var whenceOffset int
	switch whence {
	case io.SeekStart:
		whenceOffset = 0
	case io.SeekCurrent:
		whenceOffset = ss.offset
	case io.SeekEnd:
		whenceOffset = len(ss.groups)
	}
	newOffset := offset + whenceOffset
	if newOffset >= len(ss.groups) {
		return 0, ErrEof
	}
	if newOffset < 0 {
		return 0, ErrBof
	}
	ss.offset = newOffset
	return newOffset, nil
}

type stringPosition struct {
	begin int
	end   int
}

type ReversibleScanner struct {
	mmap      mmap.MMap
	pos       int
	curr      string
	offsets   map[int]stringPosition
	seperator string
	offset    int
	maxOffset int
	readAll   bool
}

func NewReversibleScanner(m mmap.MMap) *ReversibleScanner {
	ofs := make(map[int]stringPosition)
	return &ReversibleScanner{mmap: m,
		pos:       0,
		offsets:   ofs,
		seperator: "\n",
		readAll:   false,
		offset:    -1,
		maxOffset: -1,
	}
}

func (rs *ReversibleScanner) Split(sep string) {
	rs.pos = 0
	rs.seperator = sep
	rs.offsets = make(map[int]stringPosition)
	rs.readAll = false
	rs.offset = -1
	rs.maxOffset = -1
}

func (rs *ReversibleScanner) Text() string {
	return rs.curr
}

func (rs *ReversibleScanner) Prev() bool {
	rs.offset--
	if rs.offset < 0 {
		rs.curr = ""
		return false
	}
	_, err := rs.Seek(rs.offset, io.SeekStart)
	if err != nil {
		return false
	}
	return true
}

func (rs *ReversibleScanner) Next() bool {
	rs.offset++
	_, err := rs.Seek(rs.offset, io.SeekStart)
	if err != nil {
		return false
	}
	return true
}

func (rs *ReversibleScanner) Seek(offset int, whence int) (int, error) {
	var whenceOffset int
	switch whence {
	case io.SeekStart:
		whenceOffset = 0
	case io.SeekCurrent:
		whenceOffset = rs.offset
	case io.SeekEnd:
		for rs.Scan() {
		}
		whenceOffset = rs.maxOffset
	}
	newOffset := offset + whenceOffset
	if newOffset < 0 {
		return 0, ErrBof
	}
	if newOffset > rs.maxOffset {
		for rs.Scan() && newOffset > rs.maxOffset {
		}
		if newOffset > rs.maxOffset {
			return 0, ErrEof
		}
	}
	position, ok := rs.offsets[newOffset]
	if ok {
		rs.offset = newOffset
		rs.curr = string(rs.mmap[position.begin:position.end])
		return int(position.begin), nil
	}
	return 0, errors.New("Could not find offset")
}

func (rs *ReversibleScanner) Scan() bool {
	if rs.readAll {
		return false
	}
	rs.curr = ""
	var begin int
	currentPos, ok := rs.offsets[rs.maxOffset]
	if ok {
		begin = int(currentPos.end) + len(rs.seperator)
		rs.pos = begin
	} else {
		begin = 0
	}
	end := begin
	for rs.curr[max(0, len(rs.curr)-len(rs.seperator)):] != rs.seperator {
		nextRune, size, err := rs.readRune()
		if err != nil {
			rs.readAll = true
			if begin == end {
				return false
			}
		}
		rs.curr += string(nextRune)
		rs.pos += size
		end += size
	}
	rs.maxOffset++
	rs.curr = rs.curr[:len(rs.curr)-len(rs.seperator)]
	rs.offsets[rs.maxOffset] = stringPosition{begin: begin, end: end - len(rs.seperator)}
	return true
}

func (rs *ReversibleScanner) readRune() (rune, int, error) {
	if rs.pos >= len(rs.mmap) {
		return rune(0), 0, ErrEof
	}

	if rs.mmap[rs.pos] < utf8.RuneSelf {
		return rune(rs.mmap[rs.pos]), 1, nil
	}

	r, width := utf8.DecodeRune(rs.mmap[rs.pos:])
	if width > 1 {
		return r, width, nil
	}
	return utf8.RuneError, 1, nil

}
