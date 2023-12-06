package disk

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

func ReadJsonFile[E any](path string, obj *E) error {
	if obj == nil {
		return errors.New("obj is nil")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, obj)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonFile[E any](obj E, path string, perm fs.FileMode) error {
	towrite, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, towrite, perm)
	if err != nil {
		return err
	}

	return nil
}
