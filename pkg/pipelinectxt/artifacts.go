package pipelinectxt

import (
	"fmt"
	"io/ioutil"
	"log"
)

func ReadArtifactsDir() (map[string][]string, error) {

	artifactsDir := ".ods/artifacts/"
	artifactsMap := map[string][]string{}

	items, err := ioutil.ReadDir(artifactsDir)
	if err != nil {
		return artifactsMap, fmt.Errorf("%w", err)
	}

	for _, item := range items {
		if item.IsDir() {
			// artifact subdir here, e.g. "xunit-reports"
			subitems, err := ioutil.ReadDir(artifactsDir + item.Name())
			if err != nil {
				log.Fatalf("Failed to read dir %s, %s", item.Name(), err.Error())
			}
			filesInSubDir := []string{}
			for _, subitem := range subitems {
				if !subitem.IsDir() {
					// artifact file here, e.g. "report.xml"
					log.Println(item.Name() + "/" + subitem.Name())
					filesInSubDir = append(filesInSubDir, subitem.Name())
				}
			}

			artifactsMap[item.Name()] = filesInSubDir
		}
	}

	log.Println("artifactsMap: ", artifactsMap)

	// map of directories and files under .ods/artifacts, e.g
	// [
	//	"xunit-reports": ["report.xml"]
	//	"sonarqube-analysis": ["analysis-report.md", "issues-report.csv"],
	// ]

	return artifactsMap, nil
}
