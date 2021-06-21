0. Fill in values.custon.yaml
Careful: the username in the bitbucket must match the auth token
1. Upgrade Helm chart
```
cd deploy/cd-namespace
helm -n michaeltest-cd upgrade --install --values=./chart/values.custom.yaml ods-pipeline ./chart
```
2. Add secret to pipeline serviceaccount
```
oc -n michaeltest-cd patch sa pipeline --type json -p '[{"op":"add","path":"/secrets","value":[{"name":"ods-bitbucket-auth"}]}]'
```
3. Apply pipeline
4. Start pipeline
