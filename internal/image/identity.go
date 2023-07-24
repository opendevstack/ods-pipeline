package image

import (
	"fmt"

	"github.com/opendevstack/ods-pipeline/pkg/artifact"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
)

// registry/imageNamespace/imageStream:<tag>
type Identity struct {
	ImageNamespace string
	ImageStream    string
	GitCommitSHA   string // our Digest not docker digest.
}

// StreamSha renders ImageStream:GitCommitSHA
func (iid *Identity) StreamSha() string {
	return fmt.Sprintf("%s:%s", iid.ImageStream, iid.GitCommitSHA)
}

// NamespaceStreamSha renders ImageNamespace/ImageStream:GitCommitSHA
func (iid *Identity) NamespaceStreamSha() string {
	return fmt.Sprintf("%s:%s", iid.NamespaceStream(), iid.GitCommitSHA)
}

// NamespaceStream renders ImageNamespace/ImageStream aka Repository in docker terms
func (iid *Identity) NamespaceStream() string {
	return fmt.Sprintf("%s/%s", iid.ImageNamespace, iid.ImageStream)
}

func CreateImageIdentity(ctxt *pipelinectxt.ODSContext, imageNamespace, imageStream string) Identity {
	n := imageNamespace
	if len(n) == 0 {
		n = ctxt.Namespace
	}
	s := imageStream
	if len(s) == 0 {
		s = ctxt.Component
	}
	return Identity{
		ImageNamespace: n,
		ImageStream:    s,
		GitCommitSHA:   ctxt.GitCommitSHA,
	}
}

type IdentityWithTag struct {
	ImageIdentity *Identity
	Tag           string
}

// NamespaceStreamTag renders ImageNamespace/ImageStream:Tag
func (idt *IdentityWithTag) NamespaceStreamTag() string {
	return fmt.Sprintf("%s:%s", idt.ImageIdentity.NamespaceStream(), idt.Tag)
}

// NamespaceStreamSha renders ImageNamespace/ImageStream:GitCommitSHA
func (idt *IdentityWithTag) NamespaceStreamSha() string {
	return idt.ImageIdentity.NamespaceStreamSha()
}

func (iid *Identity) Tag(tag string) IdentityWithTag {
	return IdentityWithTag{
		ImageIdentity: iid,
		Tag:           tag,
	}
}

func (iid *Identity) ImageRefWithSha(registry string) string {
	return fmt.Sprintf("%s/%s", registry, iid.NamespaceStreamSha())
}

func (iid *Identity) ArtifactImage(registry string, imageDigest string) artifact.Image {
	return artifact.Image{
		Ref:        iid.ImageRefWithSha(registry),
		Registry:   registry,
		Repository: iid.ImageNamespace,
		Name:       iid.ImageStream,
		Tag:        iid.GitCommitSHA,
		Digest:     imageDigest,
	}
}

// ImageRef renders Registry/ImageNamespace/ImageStream:Tag
func (idt *IdentityWithTag) ImageRef(registry string) string {
	return fmt.Sprintf("%s/%s", registry, idt.NamespaceStreamTag())
}

// imageRef renders Registry/ImageNamespace/ImageStream:GitCommitSHA
func (idt *IdentityWithTag) ImageRefWithSha(registry string) string {
	return fmt.Sprintf("%s/%s", registry, idt.NamespaceStreamSha())
}

func (idt *IdentityWithTag) ArtifactImage(registry string, imageDigest string) artifact.Image {
	return artifact.Image{
		Ref:        idt.ImageRef(registry),
		Registry:   registry,
		Repository: idt.ImageIdentity.ImageNamespace,
		Name:       idt.ImageIdentity.ImageStream,
		Tag:        idt.Tag,
		Digest:     imageDigest,
	}
}
