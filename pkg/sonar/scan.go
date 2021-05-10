package sonar

import (
	"fmt"
)

type BitbucketServer struct {
	URL        string
	Token      string
	Project    string
	Repository string
}

type PullRequest struct {
	Key    string
	Branch string
	Base   string
}

func (c *Client) Scan(project, branch, commit string, bb *BitbucketServer, pr *PullRequest) (string, error) {
	scannerParams := []string{
		fmt.Sprintf("-Dsonar.host.url=%s", c.clientConfig.BaseURL),
		"-Dsonar.scm.provider=git",
		fmt.Sprintf("-Dsonar.projectKey=%s", project),
		fmt.Sprintf("-Dsonar.projectName=%s", project), // TODO: allow to overwrite to cater for multi-repo?
		fmt.Sprintf("-Dsonar.projectVersion=%s", commit),
	}
	if c.clientConfig.Debug {
		scannerParams = append(scannerParams, "-X")
	}
	if bb != nil && pr != nil {
		scannerParams = append(
			scannerParams,
			"-Dsonar.pullrequest.provider=Bitbucket Server",
			fmt.Sprintf("-Dsonar.pullrequest.bitbucketserver.serverUrl=%s", bb.URL),
			fmt.Sprintf("-Dsonar.pullrequest.bitbucketserver.token.secured=%s", bb.Token),
			fmt.Sprintf("-Dsonar.pullrequest.bitbucketserver.project=%s", bb.Project),
			fmt.Sprintf("-Dsonar.pullrequest.bitbucketserver.repository=%s", bb.Repository),
			fmt.Sprintf("-Dsonar.pullrequest.key=%s", pr.Key),
			fmt.Sprintf("-Dsonar.pullrequest.branch=%s", pr.Branch),
			fmt.Sprintf("-Dsonar.pullrequest.base=%s", pr.Base),
		)
	} else if c.clientConfig.ServerEdition != "community" {
		scannerParams = append(scannerParams, fmt.Sprintf("-Dsonar.branch.name=%s", branch))
	}

	fmt.Printf("scan params: %v", scannerParams)
	scannerParams = append(scannerParams, fmt.Sprintf("-Dsonar.login=%s", c.clientConfig.APIToken))
	stdout, stderr, err := runCmd("sonar-scanner", scannerParams)
	if err != nil {
		fmt.Println(string(stdout))
		fmt.Println(string(stderr))
		return "", fmt.Errorf("scanning failed: %w", err)
	}
	return string(stdout), nil
}

//     def getQualityGateJSON(String projectKey) {
//         withSonarServerConfig { hostUrl, authToken ->
//             script.sh(
//                 label: 'Get status of quality gate',
//                 script: "curl -s -u ${authToken}: ${hostUrl}/api/qualitygates/project_status?projectKey=${projectKey}",
//                 returnStdout: true
//             )
//         }
//     }

// }
