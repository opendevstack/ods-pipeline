# Executing step scripts

Tasks YAML should not include the script, the script should be a separate file (executable):

* This makes it easier to maintain the script (e.g. you get syntax highlighting in editors)
* The file location should be referenced in a comment to make it easier to find

The executable can be a Shell script or it could be a Go binary, etc. Consider

* Small things could be a shell script (makes it easier to modify the contents)
* Bigger things should not be a shell script. It could be Go (static complilation, easier to test)
