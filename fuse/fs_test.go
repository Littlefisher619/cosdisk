package fuse

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	testAbsRoot           = os.Getenv("FUSE_MOUNTPOINT")
	contentPoolSize       = 128
	dirCreateCount        = 32
	fileCreatePerDirCount = 32
	letterBytes           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type benchmark interface {
	Fatalf(format string, args ...interface{})
	ResetTimer()
	StopTimer()
	StartTimer()
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func createFile(b benchmark, dirname, file, content string) {
	f, err := os.Create(
		path.Join(testAbsRoot, dirname, file),
	)

	if err != nil {
		b.Fatalf("Failed to create file %s in %s: %s", file, dirname, err)
	}
	defer f.Close()

	f.Write([]byte(content))

}
func createDir(b benchmark, dirname string) {
	err := os.Mkdir(
		path.Join(testAbsRoot, dirname), os.ModePerm,
	)
	if err != nil {
		b.Fatalf("Failed to create dir %s: %s", dirname, err)
	}
}
func remove(b benchmark, dirOrFile string) {
	err := os.Remove(
		path.Join(testAbsRoot, dirOrFile),
	)
	if err != nil {
		b.Fatalf("Failed to remove %s: %s", dirOrFile, err)
	}
}

func createTestDirs(b benchmark, n int) []string {
	ret := make([]string, n)
	for i := 0; i < n; i++ {
		ret[i] = "dir" + strconv.Itoa(i)
		createDir(b, ret[i])
	}

	return ret
}

func prepareFileContent(n int) []string {
	ret := make([]string, n)
	for i := 0; i < n; i++ {
		ret[i] = randString(1024)
	}
	return ret
}

func createTestFiles(b benchmark, dirname string, files []string, contents []string) {
	n := len(files)
	for i := 0; i < n; i++ {
		createFile(b, dirname, files[i], contents[i])
	}
}

func removeDirAndFiles(b benchmark, dirname string, files []string) {
	for _, file := range files {
		remove(b, path.Join(dirname, file))
	}
	remove(b, dirname)
}

func runBenchmarkSingle(b benchmark) {
	b.StopTimer()
	contents := []string{""} //prepareFileContent(contentPoolSize)]
	dirs := createTestDirs(b, dirCreateCount)
	type dir struct {
		name     string
		files    []string
		contents []string
	}
	testcases := make([]dir, dirCreateCount)
	for i := 0; i < dirCreateCount; i++ {
		testcases[i].files = make([]string, fileCreatePerDirCount)
		testcases[i].contents = make([]string, fileCreatePerDirCount)
		testcases[i].name = dirs[i]
		for j := 0; j < fileCreatePerDirCount; j++ {
			testcases[i].files[j] = "file" + strconv.Itoa(j)
			testcases[i].contents[j] = contents[rand.Intn(len(contents))]
		}
	}

	b.StartTimer()
	wg := sync.WaitGroup{}
	for _, t := range testcases {
		wg.Add(1)
		go func(d dir) {
			defer wg.Done()
			createTestFiles(b, d.name, d.files, d.contents)
		}(t)
	}
	wg.Wait()

	b.StopTimer()

	wg = sync.WaitGroup{}
	for _, t := range testcases {
		wg.Add(1)
		go func(d dir) {
			defer wg.Done()
			removeDirAndFiles(b, d.name, d.files)
		}(t)
	}
	wg.Wait()
}

func cleanupFiles(b benchmark, files []string) {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			b.Fatalf("Failed to cleanup %s: %s", file, err)
		}
	}
}

func cleanup(b benchmark) {
	var (
		dirToRemove   []string
		filesToRemove []string
	)
	err := filepath.Walk(testAbsRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == testAbsRoot {
			return nil
		}
		if info.IsDir() {
			dirToRemove = append(dirToRemove, path)
		} else {
			filesToRemove = append(filesToRemove, path)
		}
		return nil
	})
	if err != nil {
		b.Fatalf("Failed to walk %s: %s", testAbsRoot, err)
	}

	cleanupFiles(b, filesToRemove)
	cleanupFiles(b, dirToRemove)
}

func BenchmarkFuse(b *testing.B) {
	cleanup(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runBenchmarkSingle(b)
	}
}
