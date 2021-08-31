# Git Subtree Based Install

## Status

Accepted

## Context

Users should be able to control the installation sources in Git in order to have a trace what was installed when. Further, they should be able to consume updates easily and understand the difference between versions.

## Decision

Support and recommend a `git subtree` based workflow for both admins and users.

## Alternatives

An alternative would be `git submodule`, but it is not ideal:

* the whole repo would be referenced, leading to confusion what is relevant for each role
* users cannot modify resources (e.g. customizing tasks) as they wish

## Consequences

The user installation guide as well as the admin installation guide will describe how to setup a local repository that uses sources from the `ods-pipeline` repo via `git subtree`.
