package study

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/joaovictornsv/cards-cli/internal/queue"
	"golang.org/x/term"
)

var ErrQuit = errors.New("study session quit")

type Input interface {
	WaitReveal() error
	ReadGrade() (queue.Grade, error)
}

type TerminalInput struct {
	in  io.Reader
	fd  int
	raw bool
}

func NewTerminalInput(in io.Reader) *TerminalInput {
	ti := &TerminalInput{in: in}
	if f, ok := in.(interface{ Fd() uintptr }); ok && term.IsTerminal(int(f.Fd())) {
		ti.fd = int(f.Fd())
		ti.raw = true
	}
	return ti
}

func (t *TerminalInput) WaitReveal() error {
	if t.raw {
		return t.readRawKey(isRevealKey)
	}
	return t.readLineKey(isRevealRune)
}

func (t *TerminalInput) ReadGrade() (queue.Grade, error) {
	if t.raw {
		var grade queue.Grade
		err := t.readRawKey(func(b []byte) (bool, error) {
			g, ok, err := parseGradeBytes(b)
			if ok {
				grade = g
				return true, nil
			}
			return false, err
		})
		if err != nil {
			return "", err
		}
		return grade, nil
	}

	scanner := bufio.NewScanner(t.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return parseGradeLine(scanner.Text())
}

func (t *TerminalInput) readRawKey(match func([]byte) (done bool, err error)) error {
	oldState, err := term.MakeRaw(t.fd)
	if err != nil {
		return fmt.Errorf("enable raw mode: %w", err)
	}
	defer func() {
		_ = term.Restore(t.fd, oldState)
	}()

	buf := make([]byte, 3)
	for {
		n, err := io.ReadFull(t.in, buf[:1])
		if err != nil {
			return err
		}
		if n == 0 {
			continue
		}

		key := buf[:1]
		if key[0] == 27 {
			n, err = io.ReadFull(t.in, buf[1:3])
			if err != nil {
				return err
			}
			key = buf[:1+n]
		}

		done, err := match(key)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}

func (t *TerminalInput) readLineKey(match func(rune) bool) error {
	scanner := bufio.NewScanner(t.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}
		return io.EOF
	}
	line := scanner.Text()
	if len(line) == 0 {
		return t.readLineKey(match)
	}
	for _, r := range line {
		if match(r) {
			return nil
		}
	}
	return fmt.Errorf("invalid key")
}

func isRevealKey(b []byte) (bool, error) {
	if len(b) == 1 {
		switch b[0] {
		case ' ', '\r', '\n':
			return true, nil
		case 'q':
			return false, ErrQuit
		}
	}
	return false, nil
}

func isRevealRune(r rune) bool {
	return r == ' ' || r == '\r' || r == '\n'
}

func parseGradeBytes(b []byte) (queue.Grade, bool, error) {
	if len(b) == 1 {
		switch b[0] {
		case '1':
			return queue.GradeAgain, true, nil
		case '2':
			return queue.GradeEasy, true, nil
		case 'r', 'R':
			return queue.GradeReplace, true, nil
		case 'q':
			return "", false, ErrQuit
		}
	}
	if len(b) == 3 && b[0] == 27 && b[1] == '[' {
		switch b[2] {
		case 'A', 'D':
			return queue.GradeAgain, true, nil
		case 'C':
			return queue.GradeEasy, true, nil
		}
	}
	return "", false, nil
}

func parseGradeLine(line string) (queue.Grade, error) {
	for _, r := range line {
		switch r {
		case '1':
			return queue.GradeAgain, nil
		case '2':
			return queue.GradeEasy, nil
		case 'r', 'R':
			return queue.GradeReplace, nil
		case 'q', 'Q':
			return "", ErrQuit
		}
	}
	return "", fmt.Errorf("invalid grade key")
}

type ScriptedInput struct {
	Grades []queue.Grade
	quitAt int // index at which WaitReveal returns ErrQuit; -1 disables
	idx    int
}

func NewScriptedInput(grades []queue.Grade) *ScriptedInput {
	return &ScriptedInput{Grades: grades, quitAt: -1}
}

func (s *ScriptedInput) WithQuitAt(index int) *ScriptedInput {
	s.quitAt = index
	return s
}

func (s *ScriptedInput) WaitReveal() error {
	if s.quitAt >= 0 && s.idx == s.quitAt {
		return ErrQuit
	}
	return nil
}

func (s *ScriptedInput) ReadGrade() (queue.Grade, error) {
	if s.idx >= len(s.Grades) {
		return "", fmt.Errorf("unexpected grade read at index %d", s.idx)
	}
	g := s.Grades[s.idx]
	s.idx++
	return g, nil
}
