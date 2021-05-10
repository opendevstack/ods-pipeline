# Requirements

* Run in OpenShift
* Read Bitbucket webhook payload
* Read `ods.yml` files from basic-auth protected endpoints
* Assemble Tekton pipeline based on phases and tasks defined in (multiple) `ods.yml` files
    * Enforce proper value of `runAfter`
    * Enforce proper value of `subPath`
* Create Tekton pipeline if it does not exist yet
* Update Tekton pipeline if it exists
* Update request body with `extensions` field
