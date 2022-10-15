package greeting

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintBuildInfo(t *testing.T) {
	for _, tc := range []struct {
		name     string
		version  string
		date     string
		commit   string
		expected string
	}{
		{
			name:    "Simple test",
			version: "0.0.1", date: "2022/07/19", commit: "01234567",
			expected: "Build version: 0.0.1\nBuild date: 2022/07/19\nBuild commit: 01234567\n",
		},
		{
			name:     "Empty fields",
			expected: "Build version: N/A\nBuild date: N/A\nBuild commit: N/A\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var builder strings.Builder

			err := PrintBuildInfo(&builder, tc.version, tc.date, tc.commit)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, builder.String())
		})
	}
}
