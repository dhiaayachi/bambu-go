package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewReportEvent(t *testing.T) {
	t.Run("Create PrintReport", func(t *testing.T) {
		event := NewReportEvent(PrintType)
		assert.NotNil(t, event)
		assert.Equal(t, PrintType, event.GetType())
	})

	t.Run("Create InfoReport", func(t *testing.T) {
		event := NewReportEvent(InfoType)
		assert.NotNil(t, event)
		assert.Equal(t, InfoType, event.GetType())
	})

	t.Run("Create UpgradeReport", func(t *testing.T) {
		event := NewReportEvent(UpgrateType)
		assert.NotNil(t, event)
		assert.Equal(t, UpgrateType, event.GetType())
	})

	t.Run("Invalid Event Type", func(t *testing.T) {
		event := NewReportEvent("unknown")
		assert.Nil(t, event)
	})
}
