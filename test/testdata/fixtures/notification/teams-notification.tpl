{
  "@type": "MessageCard",
  "@context": "http://schema.org/extensions",
  "themeColor": {{if eq .OverallStatus "Succeeded"}}"237b4b"{{else}}"c4314b"{{end }},
  "summary": "{{.ODSContext.Project}} - ODS Pipeline Run {{.PipelineRunName}} finished with status {{.OverallStatus}}",
  "sections": [
    {
      "activityTitle": "ODS Pipeline Run {{.PipelineRunName}} finished with status {{.OverallStatus}}",
      "activitySubtitle": "On Project {{.ODSContext.Project}}",
      "activityImage": "https://avatars.githubusercontent.com/u/38974438?s=200&v=4",
      "facts": [
        {
          "name": "GitRef",
          "value": "{{.ODSContext.GitRef}}"
        }
      ],
      "markdown": true
    }
  ],
  "potentialAction": [
    {
      "@type": "OpenUri",
      "name": "Go to PipelineRun",
      "targets": [
        {
          "os": "default",
          "uri": "{{.PipelineRunURL}}"
        }
      ]
    }
  ]
}
