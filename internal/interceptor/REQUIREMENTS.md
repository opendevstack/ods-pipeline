# Requirements

* Run in OpenShift
* Read Bitbucket webhook payload
* Read `ods.yml` files from basic-auth protected endpoints
* Assemble Tekton pipeline based on tasks defined in `ods.yml` file
* Create Tekton pipeline if it does not exist yet
* Update Tekton pipeline if it exists
* Update request body with `extensions` field
