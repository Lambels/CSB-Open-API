package dialogue

import (
	"bufio"
	"context"
	"io"
	"strings"
	"sync"
)

type Dialogue struct {
	Prefix string
	R      io.Reader
	W      io.Writer

	Ctx    context.Context
	ctx    context.Context
	cancel context.CancelFunc

	NotFoundHandler ExecFunc

	runningMu sync.Mutex
	running   bool
	pipe      io.Closer

	// error shared by token generator and run loop.
	err error

	// commands holds all the root commands.
	commands   map[string]*Command
	commandsWg sync.WaitGroup
}

// Start starts the main processing thread where values from the dialogue are dispatched
// to their own handlers.
//
// It reports any error returned from the ReadWriter or from parsing the
// dialogue.
func (d *Dialogue) Start() error {
	d.runningMu.Lock()
	d.running = true
	d.runningMu.Unlock()

	if d.Ctx == nil {
		d.Ctx = context.Background()
	}
	d.ctx, d.cancel = context.WithCancel(d.Ctx)

	// pipe used to decouple the processing thread from the possibly blocking
	// read writer.
	pR, pW := io.Pipe()
	d.pipe = pR

	// forward tokens to the pipe writer, they will be available in the pipe reader.
	go d.forwardTokens(pW)
	tokens := d.fetchTokens(pR)

	for tkn := range tokens {
		d.commandsWg.Add(1)
		cmd, args, err := parseRawCmd(tkn)
		if err != nil {
			d.commandsWg.Done()
			return err
		}

		if err := d.dispatchHandler(d.ctx, cmd, args); err != nil {
			d.commandsWg.Done()
			return err
		}
		d.commandsWg.Done()

		if err := d.writePrefix(); err != nil {
			return err
		}
	}
	return d.err
}

func (d *Dialogue) forwardTokens(w io.WriteCloser) {
	defer w.Close()

	scan := bufio.NewScanner(d.R)
	if err := d.writePrefix(); err != nil {
		d.err = err
		return
	}
	for scan.Scan() {
		b := scan.Bytes()
		b = append(b, byte('\n'))
		if _, err := w.Write(b); err != nil {
			d.err = err
			return
		}
	}

	d.err = scan.Err()
}

// fetchTokens fetches tokens from the reader sending them down the
// returned channel.
//
// the returned channel is closed when:
//
//  1. The reader is closed.
//
//  2. An error is encountered when writing / reading from the read writer.
func (d *Dialogue) fetchTokens(r io.Reader) <-chan []string {
	ch := make(chan []string)

	go func() {
		defer close(ch)

		scan := bufio.NewScanner(r)
		for scan.Scan() {
			s := scan.Text()
			ch <- strings.Fields(s)
		}

		d.err = scan.Err()
	}()

	return ch
}

func (d *Dialogue) writePrefix() error {
	_, err := d.W.Write([]byte(d.Prefix))
	return err
}

// dispatchHandler dispatches the handler for cmd if it exits or the not found handler.
// finally it returns any error from the handlers.
func (d *Dialogue) dispatchHandler(ctx context.Context, cmd string, args []string) error {
	command, ok := d.commands[cmd]
	if !ok {
		tmp := make([]string, 1, len(args)+1)
		tmp[0] = cmd
		copy(tmp[1:], args)

		return d.NotFoundHandler(ctx, tmp)
	}

	return command.ParseAndRunInversly(ctx, args)
}

// Stop stops the dialogue asap.
//
// If the dialogue is stuck on processing any value then its stopped immediately.
//
// If the dialogue is stuck on reading / writing to the ReadWriter,
// its stopped as soon as a value goes through.
func (d *Dialogue) Stop() {
	d.runningMu.Lock()
	defer d.runningMu.Unlock()

	if !d.running {
		return
	}

	d.pipe.Close()
	d.cancel()
	d.running = false
	d.commandsWg.Wait()
}

// StopGracefully stops any future dialogue from happening whilst waiting
// for the current dialogue to ensue.
//
// If no current dialogue is happening, the behaviour is the same as with Stop.
func (d *Dialogue) StopGracefully(ctx context.Context) {
	d.runningMu.Lock()
	defer d.runningMu.Unlock()

	d.pipe.Close()
	select {
	case <-ctx.Done():
		d.cancel()
		d.commandsWg.Wait()

	case <-signalWg(&d.commandsWg):
		d.cancel()
	}
	d.running = false
}

func signalWg(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()
	return ch
}

// Handle registers hand as a handler for val. If a handler is already registered
// for val then the old value is replaced.
//
// Handle is no-op after the dialogue started.
func (d *Dialogue) RegisterCommands(cmds ...*Command) {
	d.runningMu.Lock()
	defer d.runningMu.Unlock()

	if d.running {
		return
	}

	for _, c := range cmds {
		d.commands[c.Name] = c
	}
}
