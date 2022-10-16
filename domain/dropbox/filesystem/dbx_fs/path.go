package dbx_fs

import (
	"github.com/watermint/toolbox/domain/dropbox/model/mo_path"
	"github.com/watermint/toolbox/essentials/file/es_filesystem"
	"path/filepath"
)

func NewPath(namespaceId string, path mo_path.DropboxPath) es_filesystem.Path {
	return &dbxPath{
		namespaceId: namespaceId,
		path:        path,
	}
}

type dbxPath struct {
	namespaceId string
	path        mo_path.DropboxPath
}

func (z dbxPath) Base() string {
	return filepath.Base(z.path.LogicalPath())
}

func (z dbxPath) Path() string {
	if z.path.IsRoot() || z.path.Path() == "" {
		return "/"
	} else {
		return z.path.Path()
	}
}

func (z dbxPath) Shard() es_filesystem.Shard {
	return es_filesystem.ShardData{
		FileSystemType: FileSystemTypeDropbox,
		ShardId:        z.namespaceId,
		Attributes:     map[string]interface{}{},
	}
}

func (z dbxPath) Ancestor() es_filesystem.Path {
	return NewPath(z.namespaceId, z.path.ParentPath())
}

func (z dbxPath) Descendant(pathFragment ...string) es_filesystem.Path {
	return NewPath(z.namespaceId, z.path.ChildPath(pathFragment...))
}

func (z dbxPath) IsRoot() bool {
	return z.path.IsRoot()
}

func (z dbxPath) AsData() es_filesystem.PathData {
	return es_filesystem.PathData{
		FileSystemType: FileSystemTypeDropbox,
		EntryPath:      z.Path(),
		EntryShard:     z.Shard().AsData(),
		Attributes:     map[string]interface{}{},
	}
}

func ToDropboxPath(path es_filesystem.Path) (dbxPath mo_path.DropboxPath, err es_filesystem.FileSystemError) {
	if path.AsData().FileSystemType != FileSystemTypeDropbox {
		return nil, NewError(ErrorInvalidEntryDataFormat)
	}
	return mo_path.NewDropboxPath(path.Path()), nil
}
