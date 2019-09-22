package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommit(t *testing.T) {
	clearTestStorage()
	initTestFile(t)

	commitCommand := &commitCommand{
		getStorage:    getTestStorageClosure(),
		commitMessage: "test commit",
	}

	t.Run("returns error when file is not tracked", func(t *testing.T) {
		commitCommand.fileName = notTrackedFile
		assert.Error(t, commitCommand.run(nil))
	})

	t.Run("ok", func(t *testing.T) {
		commitCommand.fileName = trackedFileAlias
		assert.NoError(t, commitCommand.run(nil))
	})
}