package notify

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/stackdiff/internal/diff"
)

// Channel represents a notification destination.
type Channel string

const (
	ChannelStdout Channel = "stdout"
	ChannelSlack  Channel = "slack"
	ChannelFile   Channel = "file"
)

// Config holds notification settings.
type Config struct {
	Channel    Channel
	Destination string // file path or webhook URL
	Writer     io.Writer
}

// DefaultConfig returns a Config that writes to stdout.
func DefaultConfig() Config {
	return Config{
		Channel: ChannelStdout,
		Writer:  os.Stdout,
	}
}

// Notifier sends drift notifications.
type Notifier struct {
	cfg Config
}

// New creates a Notifier from cfg.
func New(cfg Config) *Notifier {
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}
	return &Notifier{cfg: cfg}
}

// Notify sends a notification if any entries represent drift.
func (n *Notifier) Notify(entries []diff.Entry) error {
	var drifted []diff.Entry
	for _, e := range entries {
		if e.IsDrift() {
			drifted = append(drifted, e)
		}
	}
	if len(drifted) == 0 {
		return nil
	}
	msg := buildMessage(drifted)
	switch n.cfg.Channel {
	case ChannelSlack:
		return sendSlack(n.cfg.Destination, msg)
	default:
		_, err := fmt.Fprintln(n.cfg.Writer, msg)
		return err
	}
}

func buildMessage(entries []diff.Entry) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("stackdiff: %d drift(s) detected\n", len(entries)))
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n", e.Status, e.Key))
	}
	return strings.TrimRight(sb.String(), "\n")
}
