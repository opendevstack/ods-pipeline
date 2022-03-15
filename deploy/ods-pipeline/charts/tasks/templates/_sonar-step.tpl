{{- define "sonar-step"}}
- name: scan-with-sonar
  # Image is built from build/package/Dockerfile.sonar.
  image: '{{.Values.registry}}/{{.Values.namespace}}/ods-sonar:{{.Values.global.imageTag}}'
  env:
    - name: HOME
      value: '/tekton/home'
    - name: SONAR_URL
      valueFrom:
        configMapKeyRef:
          key: url
          name: ods-sonar
    - name: SONAR_EDITION
      valueFrom:
        configMapKeyRef:
          key: edition
          name: ods-sonar
    - name: SONAR_AUTH_TOKEN
      valueFrom:
        secretKeyRef:
          key: password
          name: ods-sonar-auth
    - name: DEBUG
      valueFrom:
        configMapKeyRef:
          key: debug
          name: ods-pipeline
  resources: {}
  script: |
    if [ "$(params.sonar-skip)" = "true" ]; then
      echo "Skipping SonarQube analysis"
    else
      mkdir -p .ods/artifacts/sonarqube-analysis
      # sonar is built from cmd/sonar/main.go.
      sonar \
        -working-dir=$(params.working-dir) \
        -quality-gate=$(params.sonar-quality-gate)
    fi
  workingDir: $(workspaces.source.path)
{{- end}}
