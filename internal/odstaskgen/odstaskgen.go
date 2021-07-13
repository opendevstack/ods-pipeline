package odstaskgen

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// GenerateODSTaskFiles will render and place the files needed to get started creating a specific-technology ODS Task.
func GenerateODSTaskFiles(technology string) {

	// source, target
	var filesAndDirs = map[string]string{
		"internal/odstaskgen/files/Dockerfile.technology":         fmt.Sprintf("build/package/Dockerfile.%s", technology),
		"internal/odstaskgen/files/build-technology.sh":           fmt.Sprintf("build/package/scripts/build-%s.sh", technology),
		"internal/odstaskgen/files/bc-ods-build-technology.yml":   fmt.Sprintf("deploy/central/images/bc-ods-build-%s.yml", technology),
		"internal/odstaskgen/files/is-ods-build-technology.yml":   fmt.Sprintf("deploy/central/images/is-ods-build-%s.yml", technology),
		"internal/odstaskgen/files/task-ods-build-technology.yml": fmt.Sprintf("deploy/central/tasks/task-ods-build-%s.yml", technology),
	}

	for srcFilePath, dstFilePath := range filesAndDirs {

		fmt.Printf("Creating file %s\n", dstFilePath)

		var _, err = os.Stat(dstFilePath)

		if os.IsNotExist(err) {

			// create file
			file, err := os.Create(dstFilePath)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			// read template file
			templateFile, err := ioutil.ReadFile(srcFilePath)
			if err != nil {
				panic(err)
			}

			// Replace 'TECHNOLOGY'
			templateFile = []byte(strings.ReplaceAll(string(templateFile), "TECHNOLOGY", technology))

			// write file
			fmt.Fprintf(file, "%s\n", string(templateFile))
		}

	}

	// append content to existing files

	appendToFiles := map[string]string{
		"deploy/central/images/kustomization.yaml": fmt.Sprintf("  - is-ods-build-%s.yml\n", technology),
		"deploy/central/tasks/kustomization.yaml":  fmt.Sprintf("  - task-ods-build-%s.yaml\n", technology),
	}

	for filePath, fileContent := range appendToFiles {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		fmt.Printf("Appending file %s\n", filePath)

		if _, err = f.WriteString(fileContent); err != nil {
			panic(err)
		}
	}

	fmt.Println("You're all set to write contribute with your technofmty-specific ODS Task.")
	fmt.Println("Now to create a sample application that will be used to test the Task under test/testdata/workspaces/.")
	fmt.Println("Once you're done implementing and testing your Task, review the checklist in docs/creating-an-ods-task.adoc to make sure you're not missing anything.")
}
