package a2022

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	cmdCd   = regexp.MustCompile(`^\$ cd (?P<dir>.+)\s?`)
	cmdLs   = regexp.MustCompile(`^\$ ls`)
	objDir  = regexp.MustCompile(`^dir (?P<name>\w+)`)
	objFile = regexp.MustCompile(`^(?P<size>\d+) (?P<name>.+)`)
)

type directory struct {
	path      string
	parent    *directory
	childDirs map[string]*directory
	files     map[string]int
	totalSize int
}

func (d *directory) addDir(name string) *directory {
	if dir, exists := d.childDirs[name]; exists {
		return dir
	}
	newDir := makeDir(strings.Join([]string{d.path, name}, "/"), d)
	d.childDirs[name] = &newDir
	return &newDir
}

func (d *directory) addSize(size int) {
	d.totalSize += size

	if d.parent != nil {
		d.parent.addSize(size)
	}
}

func (d *directory) sumDirsUnderLimit(limit int) int {
	total := 0
	for _, dir := range d.childDirs {
		total += dir.sumDirsUnderLimit(limit)
	}
	if d.totalSize <= limit {
		total += d.totalSize
	}

	return total
}

func (d *directory) smallestDirSatisfying(min, max int) int {
	if d.totalSize > max || d.totalSize < min {
		return max
	}
	bestChild := d.totalSize
	for _, child := range d.childDirs {
		// search children against max, rather than bestChild, to avoid
		// actualBest < currentBest < child.totalSize < max
		bestChild = minInt(bestChild, child.smallestDirSatisfying(min, max))
	}

	return bestChild
}

func (d *directory) string(pad int) string {
	dirName := fmt.Sprintf("%s- %s (dir, total=%d)", strings.Repeat(" ", pad), d.path, d.totalSize)
	var files []string
	for _, dir := range d.childDirs {
		files = append(files, dir.string(pad+2))
	}
	for file, size := range d.files {
		files = append(files, fmt.Sprintf("%s- %s (file, size=%d)", strings.Repeat(" ", pad+2), file, size))
	}
	return strings.Join(append([]string{dirName}, files...), "\n")
}

func makeDir(path string, parent *directory) directory {
	return directory{
		strings.ReplaceAll(path, "//", "/"),
		parent,
		map[string]*directory{},
		map[string]int{},
		0,
	}
}

const (
	minNeeded     int = 30000000
	totalDisk     int = 70000000
	smallDirLimit int = 100000
)

func day7(in io.Reader) (int, int) {
	rootDir := makeDir("/", nil)
	currentDir := &rootDir
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case cmdCd.MatchString(line):
			dest := cmdCd.FindStringSubmatch(line)[cmdCd.SubexpIndex("dir")]

			switch dest {
			case "..":
				if currentDir.parent != nil {
					currentDir = currentDir.parent
				}
			case "/":
				currentDir = &rootDir
			default:
				if dir, ok := currentDir.childDirs[dest]; ok {
					currentDir = dir
				} else {
					currentDir = currentDir.addDir(dest)
				}
			}
		case cmdLs.MatchString(line):
			// do nothing
		case objDir.MatchString(line):
			name := objDir.FindStringSubmatch(line)[objDir.SubexpIndex("name")]
			currentDir.addDir(name)
		case objFile.MatchString(line):
			size, name := mustInt(objFile.FindStringSubmatch(line)[objFile.SubexpIndex("size")]), objFile.FindStringSubmatch(line)[objFile.SubexpIndex("name")]

			if _, exists := currentDir.files[name]; !exists {
				currentDir.files[name] = size
				currentDir.addSize(size)
			}
		}
	}
	if shouldLog {
		fmt.Println(rootDir.string(0))
	}
	needed := minNeeded - (totalDisk - rootDir.totalSize)
	return rootDir.sumDirsUnderLimit(smallDirLimit), rootDir.smallestDirSatisfying(needed, rootDir.totalSize)
}
