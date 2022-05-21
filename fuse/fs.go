package fuse

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/Littlefisher619/cosdisk/service"
	"github.com/sirupsen/logrus"
)

type FS struct {
	logger  *logrus.Entry
	cosdisk *service.CosDisk
	userId  string
	root    *Node
	nodes   map[string]*Node
}

func NewFS(username string, password string,
	cosdisk *service.CosDisk) (FS, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	gotUser, err := cosdisk.UserLogin(context.Background(), username, password)
	if err != nil {
		return FS{}, err
	}

	fs1 := FS{
		cosdisk: cosdisk,
		nodes:   make(map[string]*Node),
		userId:  fmt.Sprint(gotUser.Id),
		logger:  logger.WithField("fuse", username),
	}
	fs1.root = &Node{
		fs:    &fs1,
		name:  "/",
		mode:  os.ModeDir | 0o555,
		entry: &fuse.Dirent{Inode: 0, Name: "root", Type: fuse.DT_Dir},
		filedata: &fileData{
			content: make([]byte, 0),
		},
		IsDir: true,
	}
	return fs1, nil
}

func (f FS) Root() (fs.Node, error) {
	return f.root, nil
}

// Node implements both Node and Handle for the directory and file.
type Node struct {
	fs    *FS
	name  string
	mode  os.FileMode
	IsDir bool

	// dir
	entry *fuse.Dirent

	filedata *fileData
}

type fileData struct {
	// file
	content  []byte
	isWrited bool
	isCached bool
}

func (node Node) Attr(ctx context.Context, a *fuse.Attr) error {
	node.fs.logger.Debug("Attr", node.name)
	a.Mode = node.mode
	a.Size = uint64(len(node.filedata.content))
	return nil
}

func (node Node) Lookup(ctx context.Context, name string) (fs.Node, error) {
	path := path.Join(node.name, name)
	node.fs.logger.Debug("Lookup " + path)

	v, ok := node.fs.nodes[path]
	if ok {
		return v, nil
	}
	info, err := node.fs.cosdisk.GetFileInfo(node.fs.userId, path)
	if err != nil {
		return nil, syscall.ENOENT
	}
	new_node := Node{
		fs:       node.fs,
		name:     path,
		filedata: &fileData{},
	}
	if info.IsDir() {
		new_node.entry = &fuse.Dirent{Inode: 0, Name: path, Type: fuse.DT_Dir}
		new_node.mode = os.ModeDir | 0o555
		new_node.IsDir = true
	} else {
		new_node.entry = &fuse.Dirent{Inode: 0, Name: path, Type: fuse.DT_File}
		new_node.IsDir = false
	}
	node.fs.nodes[path] = &new_node
	return &new_node, nil
}

func (node Node) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	node.fs.logger.Debug("ReadDirAll", node.name)

	var entries []fuse.Dirent
	info, err := node.fs.cosdisk.ListFiles(node.fs.userId, node.name)
	if err != nil {
		return entries, err
	}
	for _, v := range info {
		if v.IsDir() {
			entries = append(entries, fuse.Dirent{Name: v.Name(), Type: fuse.DT_Dir})
		} else {
			entries = append(entries, fuse.Dirent{Name: v.Name(), Type: fuse.DT_File})
		}
	}
	return entries, nil
}

func (d Node) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	new_path := path.Join(d.name, req.Name)
	d.fs.logger.Debug("Mkdir", new_path)

	node, err := d.Lookup(ctx, req.Name)
	if err == nil {
		return node, nil
	}
	newDir := Node{
		fs:    d.fs,
		name:  path.Join(d.name, req.Name),
		mode:  req.Mode,
		entry: &fuse.Dirent{Inode: 0, Name: req.Name, Type: fuse.DT_Dir},
		filedata: &fileData{
			content: make([]byte, 0),
		},
		IsDir: true,
	}
	err = d.fs.cosdisk.CreateDir(d.fs.userId, new_path)
	if err != nil {
		return nil, err
	}
	d.fs.nodes[req.Name] = &newDir
	return newDir, nil
}

func (d Node) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	new_path := path.Join(d.name, req.Name)
	d.fs.logger.Debug("Remove", new_path)

	node, err := d.Lookup(ctx, req.Name)
	if err != nil {
		return err
	}
	if node.(*Node).IsDir {
		err = d.fs.cosdisk.DeleteDir(d.fs.userId, new_path)
		if err != nil {
			return err
		}
	} else {
		err = d.fs.cosdisk.DeleteFIle(d.fs.userId, new_path)
		if err != nil {
			return err
		}
	}
	delete(d.fs.nodes, new_path)
	return nil
}

func (f Node) ReadAll(ctx context.Context) ([]byte, error) {
	f.fs.logger.Debug("ReadAll", f.name)
	if !f.filedata.isCached {
		reader, err := f.fs.cosdisk.DownloadUserFileByReader(context.Background(), f.fs.userId, f.name)
		if err != nil {
			return []byte{}, err
		}
		f.filedata.content, err = ioutil.ReadAll(reader)
		if err != nil {
			return []byte{}, err
		}
		f.filedata.isCached = true
		f.filedata.isWrited = false
	}
	f.fs.logger.Debug(string(f.filedata.content))
	return f.filedata.content, nil
}

func (f Node) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	f.fs.logger.Debug("Read", f.name)
	if !f.filedata.isCached {
		reader, err := f.fs.cosdisk.DownloadUserFileByReader(context.Background(), f.fs.userId, f.name)
		if err != nil {
			return err
		}
		f.filedata.content, err = ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		f.filedata.isCached = true
		f.filedata.isWrited = false
	}
	n := 0
	if req.Size > 0 && int(req.Offset) == len(f.filedata.content) {
		return io.EOF
	}
	if len(f.filedata.content)-int(req.Offset) >= req.Size {
		n = req.Size
	} else {
		n = len(f.filedata.content) - int(req.Offset)
	}
	resp.Data = f.filedata.content[req.Offset : req.Offset+int64(n)]
	return nil
}

func (f Node) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	f.fs.logger.Debug("Write", f.name)
	if !f.filedata.isWrited {
		f.filedata.content = make([]byte, 0)
		f.filedata.isWrited = true
	}
	if int(req.Offset) > len(f.filedata.content) {
		return io.EOF
	}
	f.filedata.content = append(f.filedata.content[:req.Offset], req.Data...)
	resp.Size = len(req.Data)
	f.filedata.isWrited = true
	f.filedata.isCached = true
	return nil
}

func (f Node) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	f.fs.logger.Debug("Flush", f.name)
	if f.filedata.isWrited {
		err := f.fs.cosdisk.UploadUserFileByReader(context.Background(), f.fs.userId, f.name, bytes.NewReader(f.filedata.content))
		if err != nil {
			return err
		}
	}
	f.filedata.isCached = false
	f.filedata.isWrited = false
	return nil
}

func (f Node) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	new_path := path.Join(f.name, req.Name)
	f.fs.logger.Debug("Create", new_path)

	node, err := f.Lookup(ctx, req.Name)
	if err == nil {
		return node, node, nil
	}
	newfile := Node{
		fs:    f.fs,
		name:  new_path,
		mode:  req.Mode,
		entry: &fuse.Dirent{Inode: 0, Name: req.Name, Type: fuse.DT_File},
		filedata: &fileData{
			content: make([]byte, 0),
		},
		IsDir: false,
	}
	err = f.fs.cosdisk.UploadUserFileByReader(context.Background(), f.fs.userId, new_path, bytes.NewReader(newfile.filedata.content))
	if err != nil {
		return nil, nil, err
	}
	f.fs.nodes[new_path] = &newfile
	return newfile, newfile, nil
}
