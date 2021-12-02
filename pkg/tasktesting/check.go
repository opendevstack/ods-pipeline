package tasktesting

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/internal/command"
)

type Service string

const (
	Bitbucket Service = "7990"
	Nexus     Service = "8081"
	SonarQube Service = "9000"
)

var serviceMapping = map[Service]string{
	Bitbucket: "Bitbucket",
	Nexus:     "Nexus",
	SonarQube: "SonarQube",
}

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

func CheckServices(t *testing.T, requiredServices []Service) {
	t.Logf("Trying to reach the required services...")
	for _, port := range requiredServices {
		service := serviceMapping[port]
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s", port))
		if err != nil {
			t.Fatalf("%s needs to run for this test to be executable, but it could not be reached: %s", service, err)
		}
		t.Logf("%s reached successfully.", service)
		defer resp.Body.Close()
	}
}
