ROOTDIR=.

# Key-Value pairs.
# Dockerfile and ODS image name.
Dockerfile.go-toolset := ods-build-go
Dockerfile.sonar := ods-sonar

create-kind-with-registry:
	cd scripts && ./kind-with-registry.sh
.PHONY: create-kind-with-registry

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

# LOCAL ENVIRONMENT
# Prepare local environment by spinning up a Kubernetes cluster alongside a local registry with KinD.
# Build and push ODS images to the local registry.
prepare-local-env: create-kind-with-registry build-push-images
.PHONY: prepare-local-env