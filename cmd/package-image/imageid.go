package main

import (
	"fmt"

	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

// registry/imageNamespace/imageStream:<tag>
type imageIdentity struct {
	ImageNamespace string
	ImageStream    string
	GitCommitSHA   string // our Digest not docker digest.
}

// streamSha renders ImageStream:GitCommitSHA
func (iid *imageIdentity) streamSha() string {
	return fmt.Sprintf("%s:%s", iid.ImageStream, iid.GitCommitSHA)
}

// nsStreamSha renders ImageNamespace/ImageStream:GitCommitSHA
//
// ns is an mnemonic for namespace
func (iid *imageIdentity) nsStreamSha() string {
	return fmt.Sprintf("%s:%s", iid.nsStream(), iid.GitCommitSHA)
}

// nsStream renders ImageNamespace/ImageStream aka Repository in docker terms
//
// ns is an mnemonic for namespace
func (iid *imageIdentity) nsStream() string {
	return fmt.Sprintf("%s/%s", iid.ImageNamespace, iid.ImageStream)
}

func createImageIdentity(ctxt *pipelinectxt.ODSContext, opts *options) imageIdentity {
	imageNamespace := opts.imageNamespace
	if len(imageNamespace) == 0 {
		imageNamespace = ctxt.Namespace
	}
	imageStream := opts.imageStream
	if len(imageStream) == 0 {
		imageStream = ctxt.Component
	}
	return imageIdentity{
		ImageNamespace: imageNamespace,
		ImageStream:    imageStream,
		GitCommitSHA:   ctxt.GitCommitSHA,
	}
}

type imageIdentityWithTag struct {
	ImageIdentity *imageIdentity
	Tag           string
}

// nsStreamTag renders ImageNamespace/ImageStream:Tag
//
// ns is an mnemonic for namespace
func (idt *imageIdentityWithTag) nsStreamTag() string {
	return fmt.Sprintf("%s:%s", idt.ImageIdentity.nsStream(), idt.Tag)
}

// nsStreamSha renders ImageNamespace/ImageStream:GitCommitSHA
//
// ns is an mnemonic for namespace
func (idt *imageIdentityWithTag) nsStreamSha() string {
	return idt.ImageIdentity.nsStreamSha()
}

func (iid *imageIdentity) tag(tag string) imageIdentityWithTag {
	return imageIdentityWithTag{
		ImageIdentity: iid,
		Tag:           tag,
	}
}

func (iid *imageIdentity) shaTag() imageIdentityWithTag {
	return imageIdentityWithTag{
		ImageIdentity: iid,
		Tag:           iid.GitCommitSHA,
	}
}

func (iid *imageIdentity) imageRefWithSha(registry string) string {
	return fmt.Sprintf("%s/%s", registry, iid.nsStreamSha())
}

// imageRef renders Registry/ImageNamespace/ImageStream:Tag
//
// ns is an mnemonic for namespace
func (idt *imageIdentityWithTag) imageRef(registry string) string {
	return fmt.Sprintf("%s/%s", registry, idt.nsStreamTag())
}

// imageRef renders Registry/ImageNamespace/ImageStream:GitCommitSHA
//
// ns is an mnemonic for namespace
func (idt *imageIdentityWithTag) imageRefWithSha(registry string) string {
	return fmt.Sprintf("%s/%s", registry, idt.nsStreamSha())
}
