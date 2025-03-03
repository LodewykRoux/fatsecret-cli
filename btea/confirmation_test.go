package btea

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
)

func TestConfirmation(t *testing.T) {
	var selectedAnswer string

	tm := teatest.NewTestModel(t, InitialConfirmationModel(&selectedAnswer),
		teatest.WithInitialTermSize(300, 100),
	)

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte(""))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("enter")})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))

	assert.Equal(t, selectedAnswer, "No")
}
