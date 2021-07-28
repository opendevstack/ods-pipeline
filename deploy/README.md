# deploy

This directory contains container orchestration deployment configurations and templates.

Manifests in `central` are applied once per cluster by an ODS administrator.

Manifests in `cd-namespace` are applied once per cd-namespace by an ODS user.
The resulting resources in the `cd-namespace` use the resources (e.g. the images)
installed centrally by an ODS administrator.
