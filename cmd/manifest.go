package cmd

import (
	"fmt"

	"image-combiner/utils"
)

type Manifest struct {
	insecure bool
	name     string
	images   []*Image
}

func (m *Manifest) Initialize(name string, insecure bool, images []*Image) *Manifest {
	m.name = name
	m.insecure = insecure
	m.images = images

	return m
}

func (m *Manifest) Generate() error {
	var (
		err error
	)

	m.remove()

	if err = m.create(); err != nil {
		return err
	}

	for _, img := range m.images {
		if err = m.annotate(img); err != nil {
			return err
		}
	}

	if err = m.push(); err != nil {
		return err
	}
	return nil
}

func (m *Manifest) create() error {
	var (
		args       = []string{"docker", "manifest", "create", m.name}
		info, emsg []byte
		err        error
	)

	if insecure {
		args = append(args, "--insecure")
	}

	for _, image := range m.images {
		args = append(args, image.Name())
	}

	if info, emsg, err = utils.Exec(args); err != nil {
		fmt.Println(string(info))
		err = fmt.Errorf("create manifest error : %s  %s,", m.name, string(emsg))
		return err
	}
	return nil
}

func (m *Manifest) annotate(image *Image) error {
	var (
		args       = []string{"docker", "manifest", "annotate"}
		info, emsg []byte
		err        error
	)

	args = append(args, m.name, image.Name(), "--arch", image.arch)

	if info, emsg, err = utils.Exec(args); err != nil {
		fmt.Println(string(info))
		err = fmt.Errorf("annotate manifest error : %s", string(emsg))
		return err
	}
	return nil
}

func (m *Manifest) push() error {
	var (
		args       = []string{"docker", "manifest", "push", "--purge"}
		info, emsg []byte
		err        error
	)

	if m.insecure {
		args = append(args, "--insecure")
	}
	args = append(args, m.name)

	if info, emsg, err = utils.Exec(args); err != nil {
		fmt.Println(string(info))
		err = fmt.Errorf("push manifest error : %s", string(emsg))
		return err
	}
	return nil
}

func (m *Manifest) remove() error {
	var (
		args       = []string{"docker", "manifest", "rm"}
		info, emsg []byte
		err        error
	)

	args = append(args, m.name)

	if info, emsg, err = utils.Exec(args); err != nil {
		fmt.Println(string(info))
		err = fmt.Errorf("rm manifest error : %s", string(emsg))
		return err
	}
	return nil
}
