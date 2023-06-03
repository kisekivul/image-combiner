package utils

import (
	"os"
	"path"
)

func ListFile(dir string) ([]string, error) {
	info, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var (
		list []string
		dict map[string]struct{}
	)

	if info != nil {
		dict = make(map[string]struct{})
		for _, file := range info {
			fileName := file.Name()
			if !file.IsDir() {
				// name list
				if _, ex := dict[fileName]; !ex {
					list = append(list, path.Join(dir, fileName))
				}
				// info dict
				dict[fileName] = struct{}{}
			}
		}
	}
	return list, err
}
