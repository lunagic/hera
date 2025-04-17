package hera

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"sync"
	"syscall"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/creack/pty"
	"github.com/fatih/color"
	"github.com/lunagic/hera/hera/internal/utils"
)

func newCommandTab(
	name string,
	command string,
	onUpdate func(),
) *commandTab {
	vp := viewport.New(0, 0)
	vp.SetContent("")
	vp.GotoBottom()
	return &commandTab{
		Title:          name,
		Command:        command,
		triggerRefresh: onUpdate,
		mu:             &sync.Mutex{},
		viewport:       vp,
		status:         "🟣",
		processTracker: utils.NewProcessTracker(),
	}
}

type commandTab struct {
	Title          string
	Command        string
	triggerRefresh func()
	cmd            *exec.Cmd
	commandOutput  string
	status         string
	viewport       viewport.Model
	mu             *sync.Mutex
	processTracker *utils.ProcessTracker
}

func (c *commandTab) Init() tea.Cmd {
	if c.Command == "" {
		return nil
	}

	c.processTracker.KillAll()

	return func() tea.Msg {
		messageColor := color.New(color.Bold)

		c.mu.Lock()
		defer c.mu.Unlock()
		cmd := exec.Command("bash", "-c", c.Command)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

		c.cmd = cmd

		c.Write(
			"🔵",
			[]byte(messageColor.Sprintf("Running: %s", c.Command)),
		)

		ptyFile, err := pty.Start(cmd)
		if err != nil {
			c.Write(
				"🔴",
				[]byte(messageColor.Sprintf(
					"Error starting command: %s",
					err.Error(),
				)),
			)
			return nil
		}
		c.processTracker.Add(cmd.Process.Pid)

		reader := func(r io.Reader) {
			buf := bufio.NewReader(r)
			for {
				line, err := buf.ReadBytes(byte('\n'))
				if len(line) > 0 {
					c.Write("🔵", line)
				}
				if err != nil {
					if err == io.EOF {
						break
					}
					break
				}
			}
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			reader(ptyFile)
		}()

		wg.Wait()

		if err := c.cmd.Wait(); err != nil {
			if c.cmd.ProcessState.ExitCode() == -1 {
				c.Write(
					"🟡",
					[]byte(messageColor.Sprint("Was killed (probably to restart)")),
				)
			} else {
				c.Write(
					"🔴",
					[]byte(messageColor.Sprintf(
						"Exited with status code: %d",
						c.cmd.ProcessState.ExitCode(),
					)),
				)
			}
		} else {
			c.Write(
				"🟢",
				[]byte(messageColor.Sprint("Completed Successfully")),
			)
		}

		return nil
	}
}

func (c *commandTab) Write(status string, s []byte) {
	if !bytes.HasSuffix(s, []byte("\n")) {
		s = append(s, byte('\n'))
	}

	// check if we were at the bottom before we change anything
	atBottom := c.viewport.AtBottom()

	// Actually make the changes to the state
	c.status = status
	c.commandOutput += string(s)

	// Pre-wrap the lines so jit line wrapping doesn't confuse the viewport knowing how many lines to the bottom
	wrapped := lipgloss.NewStyle().Width(c.viewport.Width).Render(c.commandOutput)

	c.viewport.SetContent(wrapped)

	// Move to the bottom if we were already at the bottom
	if atBottom {
		c.viewport.GotoBottom()
	}

	c.triggerRefresh()
}
