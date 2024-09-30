package benchmarks

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arago-dsp/gojay"
	"github.com/arago-dsp/gojay/benchmarks"
)

func TestGoJayDecodeObjMedium(t *testing.T) {
	result := benchmarks.MediumPayload{}
	err := gojay.Unmarshal(benchmarks.MediumFixture, &result)
	require.NoError(t, err)
	assert.Equal(t, "Leonid Bugaev", result.Person.Name.FullName, "result.Person.Name.FullName should be Leonid Bugaev")
	assert.Equal(t, 95, result.Person.Github.Followers, "result.Person.Github.Followers should be 95")
	assert.Len(t, result.Person.Gravatar.Avatars, 1, "result.Person.Gravatar.Avatars should have 1 item")
}
