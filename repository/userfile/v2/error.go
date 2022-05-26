package userfile

var (
	ErrNotFound           = FilesystemError{"Your request file is not found"}
	ErrPathExists         = FilesystemError{"Your requested path is already exists a file or directory"}
	ErrIsFile             = FilesystemError{"Operation can not be performed on a file"}
	ErrIsDir              = FilesystemError{"Operation can not be performed on a directory"}
	ErrShareExpired       = FilesystemError{"The share id is expired"}
	ErrDirectoryNotEmpty  = FilesystemError{"Operation can not be performed on a non-empty directory"}
	ErrDeleteRootDir      = FilesystemError{"Cannot delete root directory"}
	ErrAddFileToFile      = FilesystemError{"Cannot add file to file"}
	ErrOverwriteFileToDir = FilesystemError{"There is a file with the same name"}
	ErrParentDirNotExist  = FilesystemError{"Parent directory does not exist"}
)

type FilesystemError struct {
	message string
}

func (e FilesystemError) Error() string {
	return e.message
}
