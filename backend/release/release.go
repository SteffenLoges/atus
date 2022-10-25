package release

import (
	"atus/backend/bencode"
	"atus/backend/config"
	"atus/backend/predb"
	"atus/backend/source"
	"atus/backend/sqlite"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ReleaseState string

const (
	StateNew          ReleaseState = "NEW"
	StateDownloadInit ReleaseState = "DOWNLOAD_INIT"
	StateDownloading  ReleaseState = "DOWNLOADING"
	StateDownloaded   ReleaseState = "DOWNLOADED"
	StateUploaded     ReleaseState = "UPLOADED"
	StateUploadError  ReleaseState = "UPLOAD_ERROR"
	StateGeneralError ReleaseState = "GENERAL_ERROR"
)

type Release struct {
	UID       string
	Hash      string
	Name      string
	NameRaw   string
	Added     time.Time
	MetaFiles []*MetaFile
	Size      int64
	Source    *source.Source
	State     ReleaseState
}

func New(ctx context.Context, s *source.Source, nameRaw string, torrentURL, imageURL *url.URL) (*Release, error) {

	rls := &Release{
		UID:     sqlite.GenerateUID("releases"),
		State:   StateNew,
		Added:   time.Now(),
		NameRaw: nameRaw,
		Source:  s,
	}

	// --------------------------------------------

	torrentFile, err := s.GetTorrentFile(ctx, torrentURL.String())
	if err != nil {
		return nil, err
	}

	dictRaw, err := bencode.BDecodeRaw(torrentFile)
	if err != nil {
		return nil, err
	}

	rls.Hash, err = dictRaw.GenHash()
	if err != nil {
		return nil, err
	}

	// ----------------------

	dict, err := bencode.BDecode(torrentFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode torrent: %s\n\nresponse: %s", err, torrentFile)
	}

	rls.Size = dict.GetSize()

	rls.Name = strings.TrimSpace(dict.Info.Name)
	if rls.Name == "" {
		return nil, errors.New("no name in meta file")
	}

	// --------------------------------------------

	metaFiles := []*MetaFile{
		NewMetaFile(rls.UID, "release.torrent", -1, MetafileTypeTorrent, MetafileStateProcessed, torrentFile, MetaInfo{
			// "url": torrentURL.String(),
		}),
	}

	// -- search torrent file for nfo, images and other metadata
	hasNFO := false
	hasSample := false
	for i, f := range dict.GetFiles() {
		file := strings.Join(f.Path, "/")
		lowerPath := strings.ToLower(file)
		ext := filepath.Ext(lowerPath)

		info := MetaInfo{
			"releasePath": file,
		}

		if ext == "" {
			continue
		}

		// nfo
		if ext == ".nfo" {
			// ignore if we already have one
			if hasNFO {
				continue
			}
			hasNFO = true

			metaFiles = append(metaFiles,
				NewMetaFile(rls.UID, fmt.Sprintf("%d_%s%s", i, rls.UID, ext), i, MetafileTypeNFO, MetafileStateUnknown, nil, info),
			)
			continue
		}

		// images
		for _, imageType := range []string{"jpg", "jpeg", "png"} {
			if ext != "."+imageType {
				continue
			}

			theType := MetafileTypeImage
			if strings.Contains(lowerPath, "proof") {
				theType = MetafileTypeProofImage
			} else if strings.Contains(lowerPath, "screen") {
				theType = MetafileTypeScreenImage
			}

			metaFiles = append(metaFiles,
				NewMetaFile(rls.UID, fmt.Sprintf("%d_%s%s", i, rls.UID, ext), i, theType, MetafileStateUnknown, nil, info),
			)
			break
		}

		// sample file.
		if config.GetBool("SAMPLES__ENABLED") {
			if strings.Contains(lowerPath, "sample") && f.Length > config.GetInt64("SAMPLES__MIN_SIZE") && f.Length < config.GetInt64("SAMPLES__MAX_SIZE") {
				// ignore if we already have one
				if hasSample {
					continue
				}
				hasSample = true

				metaFiles = append(metaFiles,
					NewMetaFile(rls.UID, fmt.Sprintf("%d_%s%s", i, rls.UID, ext), i, MetafileTypeSampleVideo, MetafileStateUnknown, nil, info),
				)
			}
		}
	}

	// -- source image
	if imageURL != nil {
		if buf, ext, err := s.GetImage(ctx, imageURL.String()); err == nil {
			metaFiles = append(metaFiles,
				NewMetaFile(rls.UID, fmt.Sprintf("source_%s%s", rls.UID, ext), -1, MetafileTypeSourceImage, MetafileStateProcessed, buf, MetaInfo{
					// "url": imageURL.String(),
				}),
			)
		}
	}

	rls.MetaFiles = metaFiles

	// --------------------------------------------

	if !hasNFO {
		return nil, errors.New("no nfo file in release")
	}

	return rls, nil

}

func (r *Release) IsKnown() bool {

	row := sqlite.Conn.QueryRow(`SELECT 1 FROM releases WHERE name = ? LIMIT 1`, r.Name)

	var isKnown bool
	if err := row.Scan(&isKnown); err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	return isKnown
}

// Save new release to database
// Do NOT call this function directly, use atus.saveNewRelease instead
func (r *Release) Save(p *predb.Pre) error {

	_, err := sqlite.Conn.Exec(
		`INSERT INTO releases (
				uid,
				hash,
				name,
				name_raw,
				state,
				pre,
				category,
				category_raw,
				size,
				added,
				source_uid
			) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.UID,
		r.Hash,
		r.Name,
		r.NameRaw,
		r.State,
		p.At.Format(time.RFC3339),
		p.Category.Name,
		p.CategoryRaw,
		r.Size,
		time.Now().Format(time.RFC3339),
		r.Source.UID,
	)

	if err != nil {
		return err
	}

	// create folder for meta files
	if _, err := os.Stat(path.Join(config.Base.Folders.Data, r.UID)); os.IsNotExist(err) {
		os.Mkdir(path.Join(config.Base.Folders.Data, r.UID), 0755)
	}

	// save meta files
	for _, meta := range r.MetaFiles {
		meta.Save()
	}

	return nil

}
