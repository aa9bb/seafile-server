package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/haiwen/seafile-server/fileserver/fsmgr"
)

const (
	seafileConfPath = "/tmp/conf"
	seafileDataDir  = "/tmp/conf/seafile-data"
	repoID          = "0d18a711-c988-4f7b-960c-211b34705ce3"
	rootID1         = "93b65479170420a64770964304531785bdc4d964"
	rootID2         = "ecb18cd8271b2fe95f3745eb43770d494230f3ef"
	fileID1         = "bfa7b0e832f82a7b47ac4e37b1e0496394f2180f"
	fileID2         = "56b1c102c47e448b4c808474fec5eb8f530ea78c"
	dirID           = "9fa74437a5ce1550fbfb89f6ea811975c5132a65"
)

func TestMain(m *testing.M) {
	initSeafile()
	m.Run()
	delFiles()
}

func TestDiffTrees(t *testing.T) {
	t.Run("test1", testDiffTrees1)
	t.Run("test2", testDiffTrees2)
	t.Run("test3", testDiffTrees3)

	delFiles()

	t.Run("test4", testDiffTrees4)
}

func initSeafile() {
	fsmgr.Init(seafileConfPath, seafileDataDir)
}

func testDiffTrees1(t *testing.T) {
	dents1 := []fsmgr.SeafDirent{}
	saveFS(dents1, rootID1)
	dent2 := fsmgr.SeafDirent{ID: dirID, Name: "dir", Mode: 0x4000}
	dents2 := []fsmgr.SeafDirent{dent2}
	saveFS(dents2, rootID2)

	var results []interface{}
	opt := &diffOptions{
		fileCB:  collectFileIDs,
		dirCB:   collectDirIDs,
		repoID:  repoID,
		results: &results}
	diffTrees([]string{rootID2, rootID1}, opt)
	ret := make([]string, len(results))
	for i, v := range results {
		ret[i] = fmt.Sprint(v)
	}
	if ret[0] != dirID {
		t.Errorf("diff error:%s!=%s", ret[0], dirID)
	}
}

func testDiffTrees2(t *testing.T) {
	dent1 := fsmgr.SeafDirent{ID: fileID1, Name: "file1", Mode: 0x81a4}
	dents1 := []fsmgr.SeafDirent{dent1}
	saveFS(dents1, rootID1)
	dent2 := fsmgr.SeafDirent{ID: dirID, Name: "dir", Mode: 0x4000}
	dents2 := []fsmgr.SeafDirent{dent2, dent1}
	saveFS(dents2, rootID2)

	var results []interface{}
	opt := &diffOptions{
		fileCB:  collectFileIDs,
		dirCB:   collectDirIDs,
		repoID:  repoID,
		results: &results}
	diffTrees([]string{rootID2, rootID1}, opt)
	ret := make([]string, len(results))
	for i, v := range results {
		ret[i] = fmt.Sprint(v)
	}
	if ret[0] != dirID {
		t.Errorf("diff error:%s!=%s", ret[0], dirID)
	}
}

func testDiffTrees3(t *testing.T) {
	dent3 := fsmgr.SeafDirent{ID: fileID2, Name: "file2", Mode: 0x81a4}
	dents3 := []fsmgr.SeafDirent{dent3}
	saveFS(dents3, dirID)

	var results []interface{}
	opt := &diffOptions{
		fileCB:  collectFileIDs,
		dirCB:   collectDirIDs,
		repoID:  repoID,
		results: &results}
	diffTrees([]string{rootID2, rootID1}, opt)
	ret := make([]string, len(results))
	for i, v := range results {
		ret[i] = fmt.Sprint(v)
	}
	if ret[0] != dirID {
		t.Errorf("diff error:%s!=%s", ret[0], dirID)
	}
	if ret[1] != fileID2 {
		t.Errorf("diff error:%s!=%s", ret[0], dirID)
	}
}

func testDiffTrees4(t *testing.T) {
	dent1 := fsmgr.SeafDirent{ID: dirID, Name: "dir", Mode: 0x4000}
	dents1 := []fsmgr.SeafDirent{dent1}
	saveFS(dents1, rootID1)
	dents2 := []fsmgr.SeafDirent{}
	saveFS(dents2, rootID2)

	var results []interface{}
	opt := &diffOptions{
		fileCB:  collectFileIDs,
		dirCB:   collectDirIDs,
		repoID:  repoID,
		results: &results}
	diffTrees([]string{rootID2, rootID1}, opt)
	if len(results) != 0 {
		t.Errorf("diff error: %v", results)
	}
}

func saveFS(dents []fsmgr.SeafDirent, dirID string) {
	seafdir := new(fsmgr.SeafDir)
	seafdir.Entries = dents
	seafdir.Version = 1
	fsmgr.SaveSeafdir(repoID, dirID, seafdir)
}

func delFiles() error {
	err := os.RemoveAll(seafileConfPath)
	if err != nil {
		return err
	}

	return nil
}
