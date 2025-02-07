package hera

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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
	c.processTracker.KillAll()

	return func() tea.Msg {
		messageColor := color.New(color.Bold)

		c.mu.Lock()
		defer c.mu.Unlock()
		cmd := exec.Command("bash", "-c", c.Command)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		c.cmd = cmd

		stdout, err := c.cmd.StdoutPipe()
		if err != nil {
			return nil
		}

		stderr, err := c.cmd.StderrPipe()
		if err != nil {
			return nil
		}

		c.Write(
			"🔵",
			messageColor.Sprintf("Running: %s", c.Command),
		)

		if err := c.cmd.Start(); err != nil {
			return nil
		}
		c.processTracker.Add(cmd.Process.Pid)

		reader := func(r io.Reader) {
			buf := bufio.NewReader(r)
			for {
				line, err := buf.ReadString('\n')
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
		wg.Add(2)
		go func() {
			defer wg.Done()
			reader(stdout)
		}()
		go func() {
			defer wg.Done()
			reader(stderr)
		}()
		wg.Wait()

		if err := c.cmd.Wait(); err != nil {
			if c.cmd.ProcessState.ExitCode() == -1 {
				c.Write(
					"🟡",
					messageColor.Sprint("Was killed (probably to restart)"),
				)
			} else {
				c.Write(
					"🔴",
					messageColor.Sprintf(
						"Exited with status code: %d",
						c.cmd.ProcessState.ExitCode(),
					),
				)
			}
		} else {
			c.Write(
				"🟢",
				messageColor.Sprint("Completed Successfully"),
			)
		}

		return nil
	}
}

func (c *commandTab) Write(status string, s string) {
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}

	// check if we were at the bottom before we change anything
	atBottom := c.viewport.AtBottom()

	// Actually make the changes to the state
	c.status = status
	c.commandOutput += s
	c.viewport.SetContent(c.commandOutput)

	// Move to the bottom if we were already at the bottom
	if atBottom {
		c.viewport.GotoBottom()
	}

	c.triggerRefresh()
}
