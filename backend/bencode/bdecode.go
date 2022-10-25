package bencode

import (
	"errors"
	"sort"
	"strings"

	"github.com/zeebo/bencode"
)

type Dict struct {
	CreatedAt int64  `bencode:"creation date" json:"creationDate"`
	CreatedBy string `bencode:"created by" json:"createdBy"`
	Announce  string `bencode:"announce" json:"announce"`
	// AnnounceList [][]string `bencode:"announce-list"`
	Comment string `bencode:"comment" json:"comment"`
	Info    *Info  `bencode:"info" json:"info"`
	Private int    `bencode:"private" json:"private"`
	Bot     string `bencode:"bot" json:"bot"`
}

type Info struct {
	Pieces      []byte `bencode:"pieces" json:"pieces"`
	PieceLength int64  `bencode:"piece length" json:"pieceLength"`
	Name        string `bencode:"name" json:"name"`

	// Length is always 0 unless this is a single file torrent.
	Length int64 `bencode:"length,omitempty" json:"length"`

	// these are only present if the torrent is a multi file torrent
	Files    []*File `bencode:"files" json:"files"`
	UniqueID string  `bencode:"unique id" json:"uniqueId"`
}

func BDecode(data []byte) (*Dict, error) {

	dict := &Dict{}
	err := bencode.DecodeBytes(data, dict)
	if err != nil {
		return nil, err
	}

	if dict.Info == nil {
		return nil, errors.New("missing 'info' in dict")
	}

	return dict, nil

}

type File struct {
	Path   []string `bencode:"path" json:"path"`
	Length int64    `bencode:"length" json:"length"`
}

// GetSize returns the size of the torrent in bytes.
func (d *Dict) GetSize() int64 {
	if d.Info.Length > 0 {
		return d.Info.Length
	}

	var size int64
	for _, file := range d.GetFiles() {
		size += file.Length
	}
	return size
}

// GetFiles returns a list of files in the torrent.
// single file torrents will return a list with a length of 1.
func (d *Dict) GetFiles() []*File {

	if len(d.Info.Files) > 0 {
		return d.Info.Files
	}

	return []*File{
		{
			Path:   []string{d.Info.Name},
			Length: d.Info.Length,
		},
	}
}

func SortFiles(files []*File) []*File {
	sortedFiles := make([]*File, len(files))
	copy(sortedFiles, files)

	sort.Slice(sortedFiles, func(i, j int) bool {
		// return files[i].Path[len(files[i].Path)-1] < files[j].Path[len(files[j].Path)-1]
		return strings.Join(files[i].Path, "/") < strings.Join(files[j].Path, "/")
	})

	return sortedFiles
}

type HierarchicalFile struct {
	Name   string `json:"name"`
	Length int64  `json:"length"`
}

type HierarchicalFolder struct {
	Path       string                 `json:"path"`
	SubFolders *[]*HierarchicalFolder `json:"subFolders,omitempty"`
	Files      *[]*HierarchicalFile   `json:"files,omitempty"`
}

func progressFileRecursively(fh *HierarchicalFolder, f *File) {

	// is this a subfolder?
	if len(f.Path) > 1 {

		if fh.SubFolders == nil {
			fh.SubFolders = &[]*HierarchicalFolder{}
		}

		// find the subfolder
		for _, subfolder := range *fh.SubFolders {
			if subfolder.Path == f.Path[0] {
				progressFileRecursively(subfolder, &File{
					Path:   f.Path[1:],
					Length: f.Length,
				})
				return
			}
		}

		// not found, create it
		subfolder := &HierarchicalFolder{
			Path: f.Path[0],
		}
		progressFileRecursively(subfolder, &File{
			Path:   f.Path[1:],
			Length: f.Length,
		})
		*fh.SubFolders = append(*fh.SubFolders, subfolder)

		return

	}

	if fh.Files == nil {
		fh.Files = &[]*HierarchicalFile{}
	}

	// it's a file
	*fh.Files = append(*fh.Files, &HierarchicalFile{
		Name:   f.Path[0],
		Length: f.Length,
	})

}

// GetFilesHierarchical returns a list of all files in a hierarchical structure.
func (d *Dict) GetFilesHierarchical() *HierarchicalFolder {

	hs := &HierarchicalFolder{
		Path: d.Info.Name,
	}

	for _, file := range d.GetFiles() {
		progressFileRecursively(hs, file)
	}

	return hs

}

func SortHierarchicalFiles(fh *HierarchicalFolder) {

	if fh.SubFolders != nil {
		sort.Slice(*fh.SubFolders, func(i, j int) bool {
			return (*fh.SubFolders)[i].Path < (*fh.SubFolders)[j].Path
		})
		for _, subfolder := range *fh.SubFolders {
			SortHierarchicalFiles(subfolder)
		}
	}

	if fh.Files != nil {
		sort.Slice(*fh.Files, func(i, j int) bool {
			return (*fh.Files)[i].Name < (*fh.Files)[j].Name
		})
	}

}
