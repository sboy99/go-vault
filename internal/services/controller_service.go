package services

import (
	"path/filepath"

	"github.com/sboy99/go-nester/pkg/utils"
)

const (
	_HEXAGONAL_CONTROLLER_DIR_PATH           = "src/controllers"
	_HEXAGONAL_CONTROLLER_PORT_DIR_PATH      = "src/controllers/ports"
	_HEXAGONAL_CONTROLLER_REGISTRY_FILE_PATH = "src/apps/controllers.ts"
	_INDEX_FILE_NAME                         = "index.ts"
)

func CreateController(name string) {
	if err := createControllerPort(_HEXAGONAL_CONTROLLER_PORT_DIR_PATH, name); err != nil {
		panic(err)
	}
	if err := createControllerFile(_HEXAGONAL_CONTROLLER_DIR_PATH, name); err != nil {
		panic(err)
	}
}

// -----------------------------------PRIVATE----------------------------------------- //

// Create controller port file
func createControllerPort(path, name string) error {
	// Create controller port file
	if err := createControllerPortInPortDir(path, name); err != nil {
		return err
	}
	if err := updateControllerPortIndexFile(path, name); err != nil {
		return err
	}
	return nil
}

// Create controller port file in port directory
func createControllerPortInPortDir(path, name string) error {
	fileName := utils.ToFileName(name, utils.CONTROLLER_PORT)
	className := utils.ToClassName(name, utils.CONTROLLER_PORT)

	content := GetPortBoilerPlateContent(className)
	return utils.WriteFile(
		filepath.Join(path, fileName),
		content,
	)
}

// Update controller port index file
func updateControllerPortIndexFile(path, name string) error {
	fileName := utils.ToFileName(name, utils.CONTROLLER_PORT)
	content := GetIndexBoilerPlateContent(fileName)
	return utils.AppendToFile(
		filepath.Join(path, _INDEX_FILE_NAME),
		content,
	)
}

// Create controller file
func createControllerFile(path, name string) error {
	// Create file
	if err := createControllerFileInControllerDir(path, name); err != nil {
		return err
	}
	if err := updateControllerIndexFile(path, name); err != nil {
		return err
	}
	return nil
}

// Create controller file in controller directory
func createControllerFileInControllerDir(path, name string) error {
	fileName := utils.ToFileName(name, utils.CONTROLLER)
	className := utils.ToClassName(name, utils.CONTROLLER)

	content := GetControllerBoilerPlateContent(className)
	return utils.WriteFile(
		filepath.Join(path, fileName),
		content,
	)
}

// Update controller index file
func updateControllerIndexFile(path, name string) error {
	fileName := utils.ToFileName(name, utils.CONTROLLER)
	content := GetIndexBoilerPlateContent(fileName)
	return utils.AppendToFile(
		filepath.Join(path, _INDEX_FILE_NAME),
		content,
	)
}
