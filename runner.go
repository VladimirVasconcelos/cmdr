package cmdr

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func run(name string, args ...string) ([]byte, error) {

	cmd := exec.Command(name, args...)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return []byte{}, err
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return []byte{}, err
	}

	// Start listening commands
	err = cmd.Start()
	if err != nil {
		return []byte{}, err
	}

	// Create buffers and listen the typing pipe.
	var outBuff bytes.Buffer
	var errBuff bytes.Buffer

	//Don't know how long it' s gonna take. Channeling...
	done := make(chan bool)

	// Copy buffered content from pipe stream
	go func() {
		io.Copy(&outBuff, outPipe)
		done <- true
	}()

	io.Copy(&errBuff, errPipe)

	// Need to wait until 'done' signal is send.
	_ = <-done

	err = cmd.Wait()

	if err != nil {
		return outBuff.Bytes(), err
	}

	if len(errBuff.Bytes()) != 0 {
		return outBuff.Bytes(), fmt.Errorf(`😢Error: %s`, errBuff.Bytes())
	}

	// Everything fine, command me!
	return outBuff.Bytes(), nil
}
