* Figure out solution for SonarQube developer edition
* Use custom built image for SonarQube to have proper plugins installed
* Enable script settings in Nexus out of the box to avoid waiting
* Stream log output from task while it is running
* Test cases for one task currently run within same namespace. That requires to clean up. Should we create one namespace per test case instead?
* Handle if BB status endpoint returns error -> stop right away
* Extract BB wait-for-ready function
