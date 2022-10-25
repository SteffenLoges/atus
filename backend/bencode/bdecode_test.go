package bencode

import (
	"testing"
)

var TestTorrent = []byte("d10:created by13:uTorrent/160013:creation datei1662738966e8:encoding5:UTF-84:infod5:filesld6:lengthi1e4:pathl14:testfile-1.txteed6:lengthi1e4:pathl14:testfile-2.txteee4:name4:test12:piece lengthi65536e6:pieces20:\xdd\xfe\x163E\xd38\x19:½\xc1\x83\xf8\xe9\xdc\xff\x90KCee")

func TestBDecode(t *testing.T) {

	dict, err := BDecode(TestTorrent)

	if err != nil {
		t.Error(err)
	}

	if dict.Info == nil {
		t.Error("missing 'info' in dict")
	}

}

func TestBDecode_GetFiles(t *testing.T) {

	dict, err := BDecode(TestTorrent)

	if err != nil {
		t.Error(err)
	}

	files := dict.GetFiles()

	t.Log(files)

	if len(files) != 2 {
		t.Error("expected 2 files")
	}

}

func TestBDecode_GenHash(t *testing.T) {

	dict, err := BDecode(TestTorrent)

	if err != nil {
		t.Error(err)
	}

	hash, err := dict.GenHash()
	if err != nil {
		t.Error(err)
	}

	t.Log(hash)

	expected := "da72195e6974c1a7c55957a625628491132965af"

	if hash != expected {
		t.Errorf("expected %s got %s", expected, hash)
	}

}
