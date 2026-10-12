package button

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/IridiumNan/link-manager/internal/models"
)

const (
	buildTestDir      = "build_test"
	buildTestSrcDir   = "build_test/src"
	buildTestDstDir   = "build_test/dst"
	buildTestInnerDir = "build_test/dst/inner"
	fileSuffix        = ".sh"
)

var srcTestFile = []string{
	"test1.sh",
	"test2.sh",
	"test3.sh",
	"test4.sh",
	"test5.sh",
	"test6.sh",
}

var buildTestNameMap = map[string]string{
	"test4.sh": "inner/test4",
	"test5.sh": "inner/test5",
	"test6.sh": "inner/test6",
}

func dirPrepare() error {
	err := os.MkdirAll(buildTestSrcDir, 0o755)
	if err != nil {
		return fmt.Errorf("error while creating build test src dir, dir: %s, err: %s", buildTestSrcDir, err.Error())
	}

	err = os.MkdirAll(buildTestInnerDir, 0o755)
	if err != nil {
		return fmt.Errorf("error while creating build test dst dir, dir: %s, err: %s", buildTestDstDir, err.Error())
	}

	return nil
}

func srcFilePrepare() error {
	for _, p := range srcTestFile {
		filePath := path.Join(buildTestSrcDir, p)
		_, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("error while creating a new file, path: %s, err: %s", filePath, err.Error())
		}
	}

	return nil
}

func dstFilePrepare() error {
	// test1.sh (no created)

	// test2.sh (file)
	filePath := path.Join(buildTestDstDir, srcTestFile[1])
	_, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error while creating %s file, err: %s", srcTestFile[1], err.Error())
	}

	// test3.sh (symlink)
	oldName := path.Join(buildTestSrcDir, srcTestFile[2])
	newName := path.Join(buildTestDstDir, srcTestFile[2])
	err = os.Symlink(oldName, newName)
	if err != nil {
		return fmt.Errorf("error while creating new symlink %s -> %s, err: %s", oldName, newName, err.Error())
	}

	// inner/test4 (no created)

	// inner/test5 (file)
	base, found := strings.CutSuffix(srcTestFile[4], fileSuffix)
	if !found {
		return fmt.Errorf("file name error: suffix not found, fix test pls. name: %s, expected suffix: %s", srcTestFile[4], fileSuffix)
	}
	filePath = path.Join(buildTestInnerDir, base)
	_, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error while creating file: %s, err: %s", filePath, err.Error())
	}

	// inner/test6 (link)
	oldName = path.Join(buildTestSrcDir, srcTestFile[5])
	base, found = strings.CutSuffix(srcTestFile[5], fileSuffix)
	if !found {
		return fmt.Errorf("file name error: suffix not found, fix test pls. name: %s, expected suffix: %s", srcTestFile[5], fileSuffix)
	}
	newName = path.Join(buildTestInnerDir, base)

	err = os.Symlink(oldName, newName)
	if err != nil {
		return fmt.Errorf("error while creating symlink: %s -> %s, err: %s", oldName, newName, err.Error())
	}

	return nil
}

func buildTestPrepare() error {
	if err := dirPrepare(); err != nil {
		return err
	}

	if err := srcFilePrepare(); err != nil {
		return err
	}
	if err := dstFilePrepare(); err != nil {
		return err
	}

	return nil
}

func TestButtonBuild(t *testing.T) {
	if _, err := os.Stat(buildTestDir); err == nil {
		t.Log("build_test dir has already existed, try to removing it")
		err = os.RemoveAll(buildTestDir)
		if err != nil {
			t.Fatalf("error while removing legacy build test dir, path: %s, err: %s", buildTestDir, err.Error())
		}
	}

	defer os.RemoveAll(buildTestDir)

	t.Log("preparing for builder test")
	err := buildTestPrepare()
	if err != nil {
		t.Fatalf("prepare: %s", err.Error())
	}

	b := NewBuilder(buildTestSrcDir, buildTestDstDir)

	b.SetNameMap(buildTestNameMap)

	buttons, err := b.Build()
	if err != nil {
		t.Fatalf("build error: %s", err.Error())
	}

	if f1 := buttons[0]; !test1(f1) {
		t.Fatalf("test1 failed, no file no inner, exp: %s, got: %s", models.StatusStr(models.StatusOff), models.StatusStr(f1.Status()))
	}
	if f2 := buttons[1]; !test2(f2) {
		t.Fatalf("test2 failed, file no inner, exp: %s, got: %s", models.StatusStr(models.StatusNone), models.StatusStr(f2.Status()))
	}
	if f3 := buttons[2]; !test3(f3) {
		t.Fatalf("test3 failed, link no inner, exp: %s, got: %s", models.StatusStr(models.StatusOn), models.StatusStr(f3.Status()))
	}
	if f4 := buttons[3]; !test4(f4) {
		t.Fatalf("test4 failed, no file inner, exp: %s, got: %s", models.StatusStr(models.StatusOff), models.StatusStr(f4.Status()))
	}
	if f5 := buttons[4]; !test5(f5) {
		t.Fatalf("test5 failed, file inner, exp: %s, got: %s", models.StatusStr(models.StatusNone), models.StatusStr(f5.Status()))
	}
	if f6 := buttons[5]; !test6(f6) {
		t.Fatalf("test6 failed, link inner, exp: %s, got: %s", models.StatusStr(models.StatusOn), models.StatusStr(f6.Status()))
	}
}

func test1(b models.NameButtonOperator) bool {
	return b.Status() == models.StatusOff
}

func test2(b models.NameButtonOperator) bool {
	return b.Status() == models.StatusNone
}

func test3(b models.NameButtonOperator) bool {
	return b.Status() == models.StatusOn
}

func test4(b models.NameButtonOperator) bool {
	return b.Status() == models.StatusOff
}

func test5(b models.NameButtonOperator) bool {
	return b.Status() == models.StatusNone
}

func test6(b models.NameButtonOperator) bool {
	return b.Status() == models.StatusOn
}
