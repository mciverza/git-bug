package identity

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/repository"
	"github.com/MichaelMure/git-bug/util/lamport"
)

func makeIdentityTestRepo(t *testing.T) repository.ClockedRepo {
	repo := repository.NewMockRepoForTest()

	clock1, err := repo.GetOrCreateClock("foo")
	require.NoError(t, err)
	err = clock1.Witness(42) // clock goes to 43
	require.NoError(t, err)

	clock2, err := repo.GetOrCreateClock("bar")
	require.NoError(t, err)
	err = clock2.Witness(34) // clock goes to 35
	require.NoError(t, err)

	return repo
}

func TestVersionSerialize(t *testing.T) {
	repo := makeIdentityTestRepo(t)

	keys := []*Key{
		{
			Fingerprint: "fingerprint1",
			PubKey:      "pubkey1",
		},
		{
			Fingerprint: "fingerprint2",
			PubKey:      "pubkey2",
		},
	}

	before, err := newVersion(repo, "name", "email", "login", "avatarUrl", keys)
	require.NoError(t, err)

	before.SetMetadata("key1", "value1")
	before.SetMetadata("key2", "value2")

	expected := &version{
		id:        entity.UnsetId,
		name:      "name",
		email:     "email",
		login:     "login",
		avatarURL: "avatarUrl",
		unixTime:  time.Now().Unix(),
		times: map[string]lamport.Time{
			"foo": 43,
			"bar": 35,
		},
		keys:  keys,
		nonce: before.nonce,
		metadata: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	require.Equal(t, expected, before)

	data, err := json.Marshal(before)
	assert.NoError(t, err)

	var after version
	err = json.Unmarshal(data, &after)
	assert.NoError(t, err)

	// make sure we now have an Id
	expected.Id()

	assert.Equal(t, expected, &after)
}
