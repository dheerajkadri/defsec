package documentdb

import (
	"testing"

	"github.com/aquasecurity/defsec/provider/aws/documentdb"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/state"
	"github.com/aquasecurity/defsec/types"
	"github.com/stretchr/testify/assert"
)

func TestCheckEnableLogExport(t *testing.T) {
	tests := []struct {
		name     string
		input    documentdb.DocumentDB
		expected bool
	}{
		{
			name: "DocDB Cluster not exporting logs",
			input: documentdb.DocumentDB{
				Metadata: types.NewTestMetadata(),
				Clusters: []documentdb.Cluster{
					{
						Metadata: types.NewTestMetadata(),
						EnabledLogExports: []types.StringValue{
							types.String("", types.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "DocDB Cluster exporting audit logs",
			input: documentdb.DocumentDB{
				Metadata: types.NewTestMetadata(),
				Clusters: []documentdb.Cluster{
					{
						Metadata: types.NewTestMetadata(),
						EnabledLogExports: []types.StringValue{
							types.String(documentdb.LogExportAudit, types.NewTestMetadata()),
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "DocDB Cluster exporting profiler logs",
			input: documentdb.DocumentDB{
				Metadata: types.NewTestMetadata(),
				Clusters: []documentdb.Cluster{
					{
						Metadata: types.NewTestMetadata(),
						EnabledLogExports: []types.StringValue{
							types.String(documentdb.LogExportProfiler, types.NewTestMetadata()),
						},
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.DocumentDB = test.input
			results := CheckEnableLogExport.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() != rules.StatusPassed && result.Rule().LongID() == CheckEnableLogExport.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}