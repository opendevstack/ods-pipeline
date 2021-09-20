package tasktesting

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
)

// Safeguard against running outside KinD
func CheckCluster(t *testing.T, outsideKindAllowed bool) {
	if !outsideKindAllowed {
		stdout, stderr, err := command.Run("kubectl", []string{"config", "current-context"})
		if err != nil {
			t.Fatalf("could not check current Kube context: %s, err: %s", string(stderr), err)
		}
		gotContext := strings.TrimSpace(string(stdout))
		wantContext := "kind-kind"
		if gotContext != wantContext {
			t.Fatalf("Not running tests outside KinD cluster ('%s') without -outside-kind! Current context: %s", wantContext, gotContext)
		}
	}
}
