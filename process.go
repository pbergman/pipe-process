package pipe

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"strings"
)

type Process []*exec.Cmd

func (p *Process) valid() error {
	if 0 == p.Len() {
		return errors.New("exec: No command(s) defined")
	}
	return nil
}

func (p *Process) Run() error {
	if e := p.valid(); e != nil {
		return e
	}
	for i, c := 1, len(*p); i < c; i++ {
		if err := (*p)[i].Start(); err != nil {
			return err
		}
	}
	if err := (*p)[0].Run(); err != nil {
		return err
	}
	for i, c := 1, len(*p); i < c; i++ {
		if err := (*p)[i].Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Process) Len() int {
	return len(*p)
}

func (p *Process) Start() error {
	if e := p.valid(); e != nil {
		return e
	}
	for i, c := 1, len(*p); i < c; i++ {
		if err := (*p)[i].Start(); err != nil {
			return err
		}
	}
	if err := (*p)[0].Start(); err != nil {
		return err
	}
	return nil
}

func (p *Process) Wait() error {
	if e := p.valid(); e != nil {
		return e
	}
	return (*p)[0].Wait()
}

func (p *Process) Output() ([]byte, error) {
	if e := p.valid(); e != nil {
		return nil, e
	}
	last := (*p)[p.Len()-1]
	if last.Stdout != nil {
		return nil, errors.New("exec: Stdout already set")
	}
	buf := new(bytes.Buffer)
	last.Stdout = buf
	capture := last.Stderr == nil
	if capture {
		last.Stderr = new(bytes.Buffer)
	}
	err := p.Run()
	if nil != err && capture {
		if ee, ok := err.(*exec.ExitError); ok {
			ee.Stderr = buf.Bytes()
		}
	}
	return buf.Bytes(), err
}

func (p *Process) CombinedOutput() ([]byte, error) {
	if e := p.valid(); e != nil {
		return nil, e
	}
	var last = (*p)[p.Len()-1]
	var buf = new(bytes.Buffer)
	if last.Stdout != nil {
		return nil, errors.New("exec: Stdout already set")
	}
	for i, c := 0, len(*p); i < c; i++ {
		if (*p)[i].Stderr != nil {
			return nil, errors.New("exec: Stderr already set")
		}
	}
	for i, c := 0, len(*p); i < c; i++ {
		(*p)[i].Stderr = buf
	}
	last.Stdout = buf
	err := p.Run()
	return buf.Bytes(), err
}

func (p *Process) StdinPipe() (io.WriteCloser, error) {
	if e := p.valid(); e != nil {
		return nil, e
	}
	return (*p)[0].StdinPipe()
}


func (p *Process) StdoutPipe() (io.ReadCloser, error) {
	if e := p.valid(); e != nil {
		return nil, e
	}
	return (*p)[p.Len()-1].StdoutPipe()
}

func NewProcess(cmd string) (*Process, error) {
	return NewProcessContext(nil, cmd)
}

func NewProcessContext(ctx context.Context, cmd string) (*Process, error) {
	var parts = strings.Split(cmd, "|")
	var stdin io.Reader
	var count = len(parts) - 1
	var process = Process(make([]*exec.Cmd, count+1))
	for i := 0; i <= count; i++ {
		var cmd *exec.Cmd
		args, err := parse(parts[i])
		if err != nil {
			return nil, err
		}
		if len(args) > 1 {
			if nil == ctx {
				cmd = exec.Command(args[0], args[1:]...)
			} else {
				cmd = exec.CommandContext(ctx, args[0], args[1:]...)
			}
		} else {
			if nil == ctx {
				cmd = exec.Command(args[0])
			} else {
				cmd = exec.CommandContext(ctx, args[0])
			}
		}
		if nil != stdin {
			cmd.Stdin = stdin
		}
		if i != count {
			pipe , err := cmd.StdoutPipe()
			if err != nil {
				return nil, err
			}
			stdin = pipe
		}
		process[i] = cmd
	}
	return &process, nil
}

