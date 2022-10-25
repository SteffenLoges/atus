package release

import (
	"atus/backend/config"
	"atus/backend/sqlite"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
	"strings"
)

type MetaFileType string

const (
	MetafileTypeTorrent               MetaFileType = "TORRENT"
	MetafileTypeNFO                   MetaFileType = "NFO"
	MetafileTypeSourceImage           MetaFileType = "SOURCE_IMAGE"
	MetafileTypeImage                 MetaFileType = "IMAGE"
	MetafileTypeProofImage            MetaFileType = "PROOF_IMAGE"
	MetafileTypeScreenImage           MetaFileType = "SCREEN_IMAGE"
	MetafileTypeScreenImageFromSample MetaFileType = "SCREEN_IMAGE__FROM_SAMPLE"
	MetafileTypeSampleVideo           MetaFileType = "SAMPLE_VIDEO"
)

type MetaFileState string

const (
	MetafileStateUnknown    MetaFileState = "UNKNOWN"    // not yet downloaded
	MetafileStateDownloaded MetaFileState = "DOWNLOADED" // downloaded but not yet processed
	MetafileStateProcessed  MetaFileState = "PROCESSED"  // ready to be uploaded
	MetafileStateError      MetaFileState = "ERROR"      // generic error
)

type MetaInfo map[string]string

type MetaFile struct {
	uid        string
	ReleaseUID string        `json:"releaseUID"`
	Index      int           `json:"-"`
	Type       MetaFileType  `json:"type"`
	State      MetaFileState `json:"state"`
	FileName   string        `json:"fileName"`
	buffer     []byte
	Info       MetaInfo `json:"info"`
}

// filename is realtive to the data folder (e.g. "uid/123.torrent")
func NewMetaFile(rlsUID, fileName string, index int, theType MetaFileType, state MetaFileState, buffer []byte, info MetaInfo) *MetaFile {
	return &MetaFile{
		ReleaseUID: rlsUID,
		Index:      index,
		Type:       theType,
		State:      state,
		FileName:   fileName,
		buffer:     buffer,
		Info:       info,
	}
}

// pass empty string as state to get all files
func GetMetaFiles(rlsUID string, state MetaFileState, theType ...MetaFileType) ([]*MetaFile, error) {

	query :=
		`SELECT
			uid,
			` + "`index`" + `,
			type,
			state,
			file_name,
			info
		FROM release_metafiles
		WHERE release_uid = ?`

	binds := []interface{}{rlsUID}

	if state != "" {
		query += " AND state = ?"
		binds = append(binds, state)
	}

	if len(theType) > 0 {
		query += " AND type IN (?" + strings.Repeat(",?", len(theType)-1) + ")"
		for _, t := range theType {
			binds = append(binds, t)
		}
	}

	rows, err := sqlite.Conn.Query(query, binds...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var metaFiles []*MetaFile
	var info []byte
	for rows.Next() {
		m := &MetaFile{
			ReleaseUID: rlsUID,
		}
		if err := rows.Scan(
			&m.uid,
			&m.Index,
			&m.Type,
			&m.State,
			&m.FileName,
			&info,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(info, &m.Info); err != nil {
			return nil, err
		}

		metaFiles = append(metaFiles, m)
	}

	return metaFiles, nil

}

func (m *MetaFile) GetFile() ([]byte, error) {
	if m.State == MetafileStateUnknown {
		if m.buffer == nil {
			return nil, errors.New("file not found")
		}

		return m.buffer, nil
	}

	return ioutil.ReadFile(path.Join(config.Base.Folders.Data, m.ReleaseUID, m.FileName))
}

func (m *MetaFile) Save() error {

	// write file to disk
	if m.buffer != nil {
		if err := ioutil.WriteFile(path.Join(config.Base.Folders.Data, m.ReleaseUID, m.FileName), m.buffer, 0644); err != nil {
			return err
		}
	}

	// uid is empty if this is a new file
	if m.uid == "" {
		m.uid = sqlite.GenerateUID("release_metafiles")
	}

	info, err := json.Marshal(m.Info)
	if err != nil {
		return err
	}

	_, err = sqlite.Conn.Exec(
		`INSERT INTO release_metafiles (
					uid,
					release_uid,
					`+"`index`"+`,
					type,
					state,
					file_name,
					info
				) 
				VALUES (?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(uid) DO UPDATE SET
					state = ?,
					file_name = ?,
					info = ?`,
		m.uid,
		m.ReleaseUID,
		m.Index,
		m.Type,
		m.State,
		m.FileName,
		info,
		m.State,
		m.FileName,
		info,
	)

	return err

}
