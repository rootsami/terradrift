package tfstack

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func Test_driftCalculator(t *testing.T) {
	state := &tfjson.Plan{
		ResourceChanges: []*tfjson.ResourceChange{
			{
				Change: &tfjson.Change{
					Actions: []tfjson.Action{
						tfjson.ActionCreate,
					},
				},
			},
		},
	}

	want := &DriftSum{
		Drift:   true,
		Add:     1,
		Change:  0,
		Destroy: 0,
	}

	got, err := driftCalculator(state)

	assert.NoError(t, err, "Unexpected error from driftCalculator")
	assert.Equal(t, want, got, "DriftSum does not match expected output")
}
