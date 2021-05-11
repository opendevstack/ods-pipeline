# Running Tekton Tests with GitHub actions

Tekton Tasks should be easily testeable. Every Tekton Task shall have a `TaskRun` associated that will serve to test the functionality of the `Task`.

## Running GitHub actions locally with [`act`](https://github.com/nektos/act)

Prerequistes:

- [Docker](https://www.docker.com/get-started)
- [KinD](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [`act`](https://github.com/nektos/act)

Configure proxy env. variables if running it from a corporate environment:

```yaml
http_proxy: <YOUR_CORPORATE_PROXY>
https_proxy: <YOUR_CORPORATE_PROXY>
HTTP_PROXY: <YOUR_CORPORATE_PROXY>
HTTPS_PROXY: <YOUR_CORPORATE_PROXY>
```

```cli
kind delete cluster --name kind-ods-github-actions && \
act -s GITHUB_TOKEN=<YOUR_GITHUB_PERSONAL_ACCESS_TOKEN>
```