# Pipeline-as-code Interceptor

## Why?

* Store pipeline in repo
* Allow easy navigation in OpenShift through pipeline to pipeline runs
* SCM-independent CI/CD solution ... not much would be needed to make this work with GitLab or GitHub. Depends more on the tasks used by the pipeline ...
* Stop-gap solution until something like https://github.com/tektoncd/pipeline/issues/859 gets worked on

## Notes

* should we read the body, or can we already make use of the extraction?
* if we can rely on the extraction, do we need to pass the params as headers?
* mirror HTTP headers
* read secret from file?
* problem: sanitized branch might not match pipelineref
* protect webhook with secret
* the workspace is predefined in the trigger template. this means that the size is the same for all repos, and that there is only one workspace config for all repos. is that good? and the name is hardcoded too. solution: do not specify PVC template at all ... probably don't want that anyway. but how big is the emptydir then? That takes up memory which we might not want. We can also link to an existing PVC, but that makes everything sequential!!

... we could have another file which is downloaded from the repo which defines the workspaces?
... we could check for repo-specific PVCs in the interceptor and use that in the trigger template if present?
... we could even auto-create according to template ...


Experiment with template (cleaned up when run is deleted):
```yml
volumeClaimTemplate:
    spec:
    accessModes:
        - ReadWriteOnce
    resources:
        requests:
        storage: 500Mi
```

## TODO

### Sanitized branch problem
See https://github.com/tektoncd/triggers/blob/master/docs/eventlisteners.md#chaining-interceptors.

We need our interceptor to store the pipeline name in the `extensions` field of the response (next to the original body). We create a custom binding then which can extract this. The custom binding then needs to be added to the event listener. This gives us access to the pipeline name in the trigger template, solving the issue of the sanitized branch.

### Where to get project, component from?

* Could write it under extensions like the sanitized branch name
* According to the docs ([The returned request (body and headers) is used as the new event payload by the EventListener and passed on the TriggerBinding](https://github.com/tektoncd/triggers/blob/master/docs/eventlisteners.md#webhook-interceptors)), the bindings run after the interceptors, this means we need to extract info from bitbucket/gitlab/github ourselves. if that's the case, we might as well store all what we extracted into the extensions field, and write a custom binding which just looks at the extensions field. then we can get project, component from there.

# Limitations

* Only one workspace, which is mapped to the same PVC for all repos right now. One could only work around this for now by having different event template triggers / event listeners. This is doable (re-using the paci interceptor) but not particularly nice. We could at least auto-create one PVC per repo ... using a default size, people can re-create using another anyway.


# Getting Started

## One-time setup for the namespace

Create a `Secret` with a Bitbucket access token (at least read perm). Annotate that secret with `tekton.dev/git-0: 'https://bitbucket.acme.org'`

Associate that secret withe the `pipeline` serviceaccount.

Create a PVC (shared workspace) to mount in each pipeline:
```yml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pipeline-<COMPONENT>
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
```

Create a trigger binding:
```yml
apiVersion: triggers.tekton.dev/v1alpha1
kind: TriggerBinding
metadata:
  name: ods-pipeline
spec:
  params:
    - name: pipeline-name
      value: $(body.extensions.name)
    - name: pipeline-pvc
      value: pipeline-$(body.extensions.pvc)

```

Create a trigger template:
```yml
apiVersion: triggers.tekton.dev/v1alpha1
kind: TriggerTemplate
metadata:
  name: ods-pipeline
spec:
  params:
    - default: pipeline-go-helm
      description: The PVC backing up the pipeline
      name: pipeline-pvc
    - description: Name of the pipeline as determined by the interceptor
      name: pipeline-name
  resourcetemplates:
    - apiVersion: tekton.dev/v1beta1
      kind: PipelineRun
      metadata:
        name: $(tt.params.pipeline-name)-$(uid)
      spec:
        params:
          - name: pipeline-run-name
            value: $(tt.params.pipeline-name)-$(uid)
        pipelineRef:
          name: $(tt.params.pipeline-name)
        serviceAccountName: pipeline
        workspaces:
          - name: shared-workspace
            persistentVolumeClaim:
              claimName: $(tt.params.pipeline-pvc)
```

Create an event listener:
```yml
apiVersion: triggers.tekton.dev/v1alpha1
kind: EventListener
metadata:
  name: ods-pipeline
spec:
  serviceAccountName: pipeline
  triggers:
    - bindings:
        - kind: TriggerBinding
          ref: ods-pipeline
      interceptors:
        - webhook:
            objectRef:
              apiVersion: v1
              kind: Service
              name: paci
              namespace: michael-cd
      name: ods-pipeline
      template:
        name: ods-pipeline
```

Expose event listener service:
```
oc expose svc el-ods-pipeline # Note the "el-" prefix!
```


## Per-repo setup

In your repository, create a file named `ods.yml` in your repository. Example:
```yml
phases:
  build:
    runPolicy: Always
    tasks:
    - name: ods-build-go
      params:
      - name: DOCKER_CONTEXT
        value: docker
      - name: TLSVERIFY
        value: 'false'
      - name: IMAGE_STREAM
        value: go-helm
      runAfter:
      - fetch-repository
      taskRef:
        kind: Task
        name: ods-build-go
      workspaces:
      - name: source
        workspace: shared-workspace
  deploy:
    tasks:
    - name: ods-deploy-helm
      params:
      - name: RELEASE_NAME
        value: go-helm
      - name: RELEASE_NAMESPACE
        value: $(params.project)-dev
      - name: PROJECT
        value: $(params.project)
      - name: REPOSITORY
        value: michael-go-helm
      - name: NAMESPACE
        value: $(params.namespace)
      runAfter:
      - ods-build-go
      taskRef:
        kind: Task
        name: ods-deploy-helm
      workspaces:
      - name: source
        workspace: shared-workspace
```

Push.
