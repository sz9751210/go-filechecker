package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type File struct {
	Path    string
	Name    string
	Size    int64
	ModTime time.Time
}

type BySize []File

func (s BySize) Len() int           { return len(s) }
func (s BySize) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySize) Less(i, j int) bool { return s[i].Size > s[j].Size }

type ByTime []File

func (t ByTime) Len() int           { return len(t) }
func (t ByTime) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTime) Less(i, j int) bool { return t[i].ModTime.After(t[j].ModTime) }

func formatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}

func main() {
	var targetPath string
	fmt.Print("請輸入要查看的路徑: ")
	fmt.Scan(&targetPath)

	var numFiles int
	fmt.Print("請輸入要顯示的檔案筆數: ")
	fmt.Scan(&numFiles)

	var sortOption int
	fmt.Println("請選擇排序方式:")
	fmt.Println("1. 按大小排序")
	fmt.Println("2. 按修改時間排序")
	fmt.Print("選擇: ")
	fmt.Scan(&sortOption)

	var from, to time.Time
	if sortOption == 2 {
		var fromStr, toStr string
		fmt.Print("請輸入起始日期（YYYY-MM-DD）: ")
		fmt.Scan(&fromStr)
		fmt.Print("請輸入結束日期（YYYY-MM-DD）: ")
		fmt.Scan(&toStr)

		from, _ = time.Parse("2006-01-02", fromStr)
		to, _ = time.Parse("2006-01-02", toStr)
	}

	var files []File
	err := filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			modTime := info.ModTime()
			if sortOption == 2 && (modTime.Before(from) || modTime.After(to)) {
				return nil
			}
			files = append(files, File{
				Path:    path,
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: modTime,
			})
		}
		return nil
	})
	if err != nil {
		fmt.Println("讀取檔案錯誤:", err)
		return
	}

	switch sortOption {
	case 1:
		sort.Sort(BySize(files))
		fmt.Println("按大小排序的檔案清單:")
	case 2:
		sort.Sort(ByTime(files))
		fmt.Println("按修改時間排序的檔案清單:")
	default:
		fmt.Println("無效的選擇")
		return
	}

	fmt.Printf("前 %d 筆檔案清單:\n", numFiles)
	for i, file := range files {
		if i >= numFiles {
			break
		}
		fmt.Printf("路徑: %s, 檔案名稱: %s, 大小/修改時間: %s, %s\n",
			file.Path, file.Name, formatSize(file.Size), file.ModTime)
	}
}
