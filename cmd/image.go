package cmd

import (
	"fmt"
	"strings"

	"image-combiner/utils"
)

type Image struct {
	name     string
	arch     string
	source   string
	manifest string
}

func (i *Image) Name() string {
	return i.name
}

func (i *Image) Arch() string {
	return i.arch
}

func (i *Image) Source() string {
	return i.arch
}

func (i *Image) Manifest() string {
	return i.manifest
}

func (i *Image) Load(file, arch, registry string) error {
	var (
		res   []byte
		image string
		err   error
	)

	i.arch = arch

	if res, _, err = utils.Exec([]string{
		"docker", "load", "-i", file,
	}); err != nil {
		return err
	}

	image = strings.Replace(strings.SplitN(string(res), `/`, 2)[1], "\n", "", -1)
	i.source = strings.Join([]string{registry, image}, `/`)
	if insecure {
		i.manifest = strings.Join([]string{registry, ":80", `/`, image}, "")
	} else {
		i.manifest = i.source
	}
	i.name = strings.Join([]string{i.manifest, ".", i.arch}, "")

	return nil
}

func (i *Image) Tag() error {
	var (
		res []byte
		err error
	)

	if res, _, err = utils.Exec([]string{
		"docker", "tag", i.source, i.name,
	}); err != nil {
		fmt.Println(string(res))
		return err
	}
	return nil
}

func (i *Image) Push(image string) error {
	var (
		err error
	)

	if _, _, err = utils.Exec([]string{
		"docker", "push", image,
	}); err != nil {
		return err
	}
	return nil
}

func (i *Image) Remove(image string) error {
	var (
		err error
	)

	if _, _, err = utils.Exec([]string{
		"docker", "rmi", image,
	}); err != nil {
		return err
	}
	return nil
}

type Images struct {
	registry  string
	insecure  bool
	manifests map[string][]*Image
}

func (is *Images) Initialize(registry string, insecure bool) *Images {
	is.registry = registry
	is.insecure = insecure
	is.manifests = make(map[string][]*Image)

	return is
}

func (is *Images) Load(arch, dir string) error {
	var (
		files []string
		err   error
	)

	if files, err = utils.ListFile(dir); err != nil {
		return err
	}

	for _, file := range files {
		if err = is.load(file, arch); err != nil {
			return err
		}

	}
	return nil
}

func (is *Images) load(file, arch string) error {
	var (
		image = &Image{}
		exist bool
		err   error
	)

	if err = image.Load(file, arch, is.registry); err != nil {
		fmt.Printf("load image from %s fail\n", file)
		return err
	}

	if err = image.Tag(); err != nil {
		fmt.Printf("tag image %s fail\n", image.source)
		return err
	}
	image.Remove(image.source)

	if err = image.Push(image.name); err != nil {
		fmt.Printf("push image %s fail\n", image.name)
		return err
	}
	image.Remove(image.name)

	if _, exist = is.manifests[image.manifest]; !exist {
		is.manifests[image.manifest] = make([]*Image, 0)
	}
	is.manifests[image.manifest] = append(
		is.manifests[image.manifest],
		image,
	)
	return nil
}

func (is *Images) Manifests() map[string][]*Image {
	return is.manifests
}
