package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func rootFiles(root string) ([]string, error) {
	var files []string

	skipExt := map[string]struct{}{
		".json": {},
		".md":   {},
		".txt":  {},
		".mod":  {},
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if _, skip := skipExt[filepath.Ext(path)]; skip {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func main() {
	root := flag.String("dir", "dir", "Add full directory to file")
	useCwd := flag.Bool("cwd", false, "Use current working directory as root")
	flag.Parse()

	dir := *root

	if *useCwd {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get working directory:", err)
			return
		}
		dir = wd
	}

	fmt.Println("---------------------------------")

	files, err := rootFiles(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Could not open file %s: %v\n", file, err)
			continue
		}

		lines := 0

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				lines++
			}
		}

		f.Close()

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		fmt.Printf("%-20s: %d lines\n", filepath.Base(file), lines)
	}

	fmt.Println("---------------------------------")
}
