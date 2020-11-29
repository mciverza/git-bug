package bug

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
)

func TestAddCommentSerialize(t *testing.T) {
	repo := repository.NewMockRepoClock()

	rene, err := identity.NewIdentity(repo, "René Descartes", "rene@descartes.fr")
	require.NoError(t, err)

	unix := time.Now().Unix()
	before := NewAddCommentOp(rene, unix, "message", nil)

	data, err := json.Marshal(before)
	assert.NoError(t, err)

	var after AddCommentOperation
	err = json.Unmarshal(data, &after)
	assert.NoError(t, err)

	// enforce creating the ID
	before.Id()

	// Replace the identity stub with the real thing
	assert.Equal(t, rene.Id(), after.base().Author.Id())
	after.Author = rene

	assert.Equal(t, before, &after)
}
