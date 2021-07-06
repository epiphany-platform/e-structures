package globals

type Initializer interface {

	// Init is responsible for building default version of structure to start from.
	Init(moduleVersion string)
}

type Backuper interface { // TODO that name is not a real word

	// Backup is responsible for storing current version of structure in new file. It should fail if there is
	// file already in place pointed by path.
	Backup(path string) error // TODO this should not take path as string
}

type Loader interface {

	// Load is responsible for loading file pointed with path argument. It should fail in following cases:
	// - if file provided in path argument is not present
	// - if version in file is not version expected by structure
	// - if validation of loaded structure failed
	//
	// In case of incorrect version fallback to Upgrader.Upgrade method should be applied by module.
	// In case of failed validation (provided by Validator.Valid method) it should be considered
	// panic situation and usually user is forced to fix file.
	Load(path string) error
}

type Saver interface {

	// Save is responsible for storing structure into file. It must always use Validator.Valid method
	// to ensure that saved file is not corrupted. It should (but not must) use Printer.Print method
	// to produce structure JSON.
	Save(path string) error
}

type Printer interface {

	// Print is responsible for producing JSON form of structure.
	Print() ([]byte, error)
}

type Validator interface {

	// Valid is responsible for checking that structure is correct with all custom validation rules.
	Valid() error
}

type Upgrader interface {

	// Upgrade is responsible for upgrading structure to current version. It is designed to be a fallback
	// method after Loader.Load wasn't able to load structure from file and returned with incorrect version
	// error.
	Upgrade(path string) error
}