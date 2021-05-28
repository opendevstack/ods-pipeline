ROOTDIR=.

# Key-Value pairs.
# Dockerfile and ODS image name.
Dockerfile.go-toolset := ods-build-go
Dockerfile.sonar := ods-sonar

create-kind-with-registry:
	cd scripts && ./kind-with-registry.sh
.PHONY: create-kind-with-registry

# TODO: Find out how to build only specific targets.
# Make them parameterized?
# Move it onto a file
build-push-images: $(ROOTDIR)/build/package/Dockerfile.*
		for file in $^ ; do \
			imageName=$$(basename $$file); \
			echo "Image Name: $$imageName"; \
			suffix="$(suffix $$(imageName))"; \
			imageName=$(Dockerfile.$(KEY))
			docker build -f $$file -t localhost:5000/ods/$$imageName:latest . \
			docker push localhost:5000/ods/$$imageName:latest \
        done
.PHONY: build-push-images

load-images-in-kind:
	cd scripts && ./load-images-in-kind.sh
.PHONY: load-images-in-kind

install-tekton-pipelines:
	cd scripts && ./install-tekton-pipelines.sh
.PHONY: install-tekton-pipelines

prepare-local-env: create-kind-with-registry load-images-in-kind install-tekton-pipelines 
.PHONY: prepare-local-env