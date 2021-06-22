* Figure out solution for SonarQube developer edition
* Use custom built image for SonarQube to have proper plugins installed
* Stream log output from task while it is running
* Test cases for one task currently run within same namespace. That requires to clean up. Should we create one namespace per test case instead?
* Handle if BB status endpoint returns error -> stop right away
* Extract BB wait-for-ready function
* Write central install logic (images for tasks, images for webhook interceptor, cluster tasks)
* Write local install logic (configmaps, secrets, event listener, webhook interceptor). could be a Helm chart or an operator or such. How to handle secrets (e.g. Bitbucket token)?


PRIO 1
* finish task
* helm task basic
* defining environment / release concept
* aqua scanning of built image -> allow to run it locally, disable scanning in github


PRIO 2
* local install, central install logix
* run manual tests on openshift
* release concept / multi-repo
