# Namespaced installation

Date: 2022-03-08

## Status

Accepted

Supersedes [Release name suffix](20211025-release-name-suffix.md)

## Context

Build tasks are not easily customizable (e.g. with respect to resource consumption/quotas or sidecars) in a centralized deployment model. In addition, for every change to a task, cluster admin rights are required.

## Decision

Abolish a central installation altogether. Provide images from a central registry (ghcr.io).

## Consequences

Every project needs to build images themselves, resulting in the installation being 

1. simpler to understand as there is no admin/user distinction
2. simpler to install as no collaboration with a cluster admin is required at all
3. simpler to customize since images are then controlled by users

Projects can use or base their images on public `ods-pipeline` images available via ghcr.io, thus greatly lowering the burden this approach entails.
