package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
	"k8s.io/utils/strings/slices"
)

// → skopeo inspect --tls-verify=false --format="{{.Digest}} {{.RepoTags}}" docker://localhost:5000/ods/ods-buildah
// sha256:ccf11cd07321c135d6b1949d367c80614ca49f37ade987148c24c2b704dd9426 [latest cool]

func (p *packageImage) skopeoMissingTags() ([]string, error) {
	p.logger.Infof("Inspecting image %s for tags", p.imageShaTag.ImageIdentity.nsStreamSha())
	tlsVerify := p.opts.tlsVerify
	// TLS verification of the KinD registry is not possible at the moment as
	// requests error out with "server gave HTTP response to HTTPS client".
	if strings.HasPrefix(p.opts.registry, "kind-registry.kind") {
		tlsVerify = false
	}
	args := []string{
		"inspect",
		`--format={{.RepoTags}}`,
		fmt.Sprintf("--tls-verify=%v", tlsVerify),
	}
	if tlsVerify {
		args = append(args, fmt.Sprintf("--cert-dir=%v", p.opts.certDir))
	}
	if p.opts.debug {
		args = append(args, "--debug")
	}
	imageRef := fmt.Sprintf("docker://%s", p.imageShaTag.imageRefWithSha(p.opts.registry))
	args = append(args, imageRef)

	p.logger.Infof("skopeo inspect %s", imageRef)
	p.logger.Infof("skopeo %+q", args)
	stdout, _, err := command.RunBuffered("skopeo", args)
	if err != nil {
		return nil, fmt.Errorf("skopeo inspect: %w", err)
	}
	tags, err := p.parseSkopeoInspectDigestTags(string(stdout))
	if err != nil {
		return nil, err
	}
	p.logger.Infof("skopeo tags=%v", tags)
	missingTags := []string{}
	for _, extraTag := range p.parsedExtraTags {
		if !slices.Contains(tags, extraTag) {
			missingTags = append(missingTags, extraTag)
		}
	}
	return missingTags, nil
}

func (p *packageImage) parseSkopeoInspectDigestTags(out string) ([]string, error) {
	t := strings.TrimSpace(out)
	p.logger.Debugf("skopeo output=%s", t)
	if !(strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]")) {
		return nil, fmt.Errorf("skopeo inspect: unexpected tag response expecting tags to be in brackets %s", t)
	}
	t = t[1 : len(t)-1]
	// expecting t to have space separated tags.
	tags := strings.Split(t, " ")
	return tags, nil
}

func (p *packageImage) skopeoTag(idt *imageIdentityWithTag, outWriter, errWriter io.Writer) error {
	imageRef := idt.imageRefWithSha(p.opts.registry)
	p.logger.Infof("Tagging image %s with %s", imageRef, idt.Tag)
	tlsVerify := p.opts.tlsVerify
	// TLS verification of the KinD registry is not possible at the moment as
	// requests error out with "server gave HTTP response to HTTPS client".
	if strings.HasPrefix(p.opts.registry, "kind-registry.kind") {
		tlsVerify = false
	}
	args := []string{
		"copy",
		fmt.Sprintf("--src-tls-verify=%v", tlsVerify),
		fmt.Sprintf("--dest-tls-verify=%v", tlsVerify),
	}
	if tlsVerify {
		args = append(args,
			fmt.Sprintf("--src-cert-dir=%v", p.opts.certDir),
			fmt.Sprintf("--dest-cert-dir=%v", p.opts.certDir))
	}
	if p.opts.debug {
		args = append(args, "--debug")
	}
	source := fmt.Sprintf("docker://%s", imageRef)
	destination := fmt.Sprintf("docker://%s", idt.imageRef(p.opts.registry))

	args = append(args, source, destination)

	p.logger.Infof("skopeo copy %s %s", source, destination)
	err := command.Run("skopeo", args, []string{}, outWriter, errWriter)
	if err != nil {
		return fmt.Errorf("skopeo copy %s to %s: %w", source, destination, err)
	}
	return nil
}
