package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsHiddenDirectory 路径是否是隐藏路径
func IsHiddenDirectory(path string) bool {
	return len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")
}

// SubDir 输出所有子目录，目录名
func SubDir(folder string) ([]string, error) {
	subs, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, sub := range subs {
		if sub.IsDir() {
			ret = append(ret, sub.Name())
		}
	}
	return ret, nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// CopyFolder 将一个目录复制到另外一个目录中
func CopyFolder(source, destination string) error {
	var err error = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), 0755)
		} else {
			var data, err1 = os.ReadFile(filepath.Join(source, relPath))
			if err1 != nil {
				return err1
			}
			return os.WriteFile(filepath.Join(destination, relPath), data, 0777)
		}
	})
	return err
}

// CopyFile 将一个目录复制到另外一个目录中
func CopyFile(source, destination string) error {
	var data, err1 = os.ReadFile(source)
	if err1 != nil {
		return err1
	}
	return os.WriteFile(destination, data, 0777)
}

// Close wraps an io.Closer and logs an error if one is returned. When
// manipulating the file, prefer util.WithReadFile over util.Close, as
// it handles closing automatically.
func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		// Don't start the stacktrace here, but at the caller's location
		log.Printf("Failed to Close(): %v", err)

	}
}

// RemoveAll removes the specified path and logs an error if one is returned.
func RemoveAll(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Printf("Failed to RemoveAll(%s): %v", path, err)
	}
}

// Remove removes the specified file and logs an error if one is returned.
func Remove(name string) {
	if err := os.Remove(name); err != nil {
		log.Printf("Failed to Remove(%s): %v", name, err)
	}
}

// IsDirEmpty checks to see if the specified directory has any contents.
func IsDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer Close(f)

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// FileExists returns true if the given path exists.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	return false, err
}

// DirExists returns true if the given path exists and is a directory.
func DirExists(path string) (bool, error) {
	exists, _ := FileExists(path)
	fileInfo, _ := os.Stat(path)
	if !exists || !fileInfo.IsDir() {
		return false, fmt.Errorf("path either doesn't exist, or is not a directory <%s>", path)
	}
	return true, nil
}

// Touch creates an empty file at the given path if it doesn't already exist.
func Touch(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

// EnsureDir will create a directory at the given path if it doesn't already exist.
func EnsureDir(path string) error {
	exists, err := FileExists(path)
	if !exists {
		err = os.Mkdir(path, 0o755)
		return err
	}
	return err
}

// EnsureDirAll will create a directory at the given path along with any necessary parents if they don't already exist.
func EnsureDirAll(path string) error {
	return os.MkdirAll(path, 0o755)
}

// RemoveDir removes the given dir (if it exists) along with all of its contents.
func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

// EmptyDir will recursively remove the contents of a directory at the given path.
func EmptyDir(path string) error {
	d, err := os.Open(path)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))
		if err != nil {
			return err
		}
	}

	return nil
}

// ListDir will return the contents of a given directory path as a string slice.
func ListDir(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		path = filepath.Dir(path)
		files, _ = ioutil.ReadDir(path)
	}

	//nolint: prealloc
	var dirPaths []string
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		dirPaths = append(dirPaths, filepath.Join(path, file.Name()))
	}
	return dirPaths
}

// GetHomeDirectory returns the path of the user's home directory.  ~ on Unix and C:\Users\UserName on Windows.
func GetHomeDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

// SafeMove move src to dst in safe mode.
func SafeMove(src, dst string) error {
	err := os.Rename(src, dst)
	//nolint: nestif
	if err != nil {
		fmt.Printf("[fileutil] unable to rename: \"%s\" due to %s. Falling back to copying.", src, err.Error())

		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}

		err = out.Close()
		if err != nil {
			return err
		}

		err = os.Remove(src)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsZipFileUncompressed returns true if zip file in path is using 0 compression level.
func IsZipFileUncompressed(path string) (bool, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		fmt.Printf("Error reading zip file %s: %s\n", path, err)
		return false, err
	}
	defer r.Close()
	for _, f := range r.File {
		if f.FileInfo().IsDir() { // skip dirs, they always get store level compression
			continue
		}
		return f.Method == 0, nil // check compression level of first actual  file
	}
	return false, nil
}

// WriteFile writes file to path creating parent directories if needed.
func WriteFile(path string, file []byte) error {
	pathErr := EnsureDirAll(filepath.Dir(path))
	if pathErr != nil {
		return fmt.Errorf("cannot ensure path %s", pathErr)
	}

	err := ioutil.WriteFile(path, file, 0o600)
	if err != nil {
		return fmt.Errorf("write error for thumbnail %s: %s ", path, err)
	}
	return nil
}

// GetIntraDir returns a string that can be added to filepath.Join to implement directory depth, "" on error
// eg given a pattern of 0af63ce3c99162e9df23a997f62621c5 and a depth of 2 length of 3
// returns 0af/63c or 0af\63c ( dependin on os)  that can be later used like this  filepath.Join(directory, intradir,
// basename).
func GetIntraDir(pattern string, depth, length int) string {
	if depth < 1 || length < 1 || (depth*length > len(pattern)) {
		return ""
	}
	intraDir := pattern[0:length] // depth 1 , get length number of characters from pattern
	for i := 1; i < depth; i++ {  // for every extra depth: move to the right of the pattern length positions, get length number of chars
		intraDir = filepath.Join(
			intraDir,
			pattern[length*i:length*(i+1)],
		) //  adding each time to intradir the extra characters with a filepath join
	}
	return intraDir
}

// GetParent returns the parent directory of the given path.
func GetParent(path string) *string {
	isRoot := path[len(path)-1:] == "/"
	if isRoot {
		return nil
	}
	parentPath := filepath.Clean(path + "/..")
	return &parentPath
}

// ServeFileNoCache serves the provided file, ensuring that the response
// contains headers to prevent caching.
func ServeFileNoCache(w http.ResponseWriter, r *http.Request, filepath string) {
	w.Header().Add("Cache-Control", "no-cache")

	http.ServeFile(w, r, filepath)
}

// MatchEntries returns a string slice of the entries in directory dir which
// match the regexp pattern. On error an empty slice is returned
// MatchEntries isn't recursive, only the specific 'dir' is searched
// without being expanded.
func MatchEntries(dir, pattern string) ([]string, error) {
	var res []string
	var err error

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if re.Match([]byte(file)) {
			res = append(res, filepath.Join(dir, file))
		}
	}
	return res, err
}

// OutDir creates the absolute path name from path and checks path exists.
// Returns absolute path including trailing '/' or error if path does not exist.
func OutDir(path string) (string, error) {
	outDir, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	stat, err := os.Stat(outDir)
	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("output directory %s is not a directory", outDir)
	}
	outDir = outDir + "/"
	return outDir, nil
}
