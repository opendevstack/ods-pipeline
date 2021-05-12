# Naming Conventions

## Tasks
hyphenated, e.g. `buildah-v0-22-0`

Reason: seems to be de-facto standard (catalog tasks are named like this)

## Parameters
hyphenated, e.g. `params: sonarqube-branch`

Reason: no de-facto standard seems to exist, hyphenated aligns with task names, works well with prefix/suffix, avoids upper/lowercase weirdness for e.g. `enableHTTPProxy` etc.)
