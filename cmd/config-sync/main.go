package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	syncSrc := syncCmd.String("source", "", "Source directory")
	syncDst := syncCmd.String("destination", "", "Destination directory")

	hashCmd := flag.NewFlagSet("hash", flag.ExitOnError)
	hashFilePath := hashCmd.String("file", "", "File to hash")

	if len(os.Args) < 2 {
		fmt.Println("Usage: config-sync <command> [arguments]")
		fmt.Println("Commands: sync, hash")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "sync":
		syncCmd.Parse(os.Args[2:])
		if *syncSrc == "" || *syncDst == "" {
			fmt.Println("Error: --source and --destination are required")
			os.Exit(1)
		}
		syncDirs(*syncSrc, *syncDst)
	case "hash":
		hashCmd.Parse(os.Args[2:])
		if *hashFilePath == "" {
			fmt.Println("Error: --file is required")
			os.Exit(1)
		}
		computeFileHash(*hashFilePath)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func fileHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func computeFileHash(path string) {
	hash := fileHash(path)
	if hash == "" {
		fmt.Printf("Error hashing %s\n", path)
		os.Exit(1)
	}
	fmt.Printf("%s  %s\n", hash, path)
}

func syncDirs(src, dst string) {
	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, rel)

		os.MkdirAll(filepath.Dir(dstPath), 0755)

		srcHash := fileHash(path)
		dstHash := fileHash(dstPath)

		if srcHash != dstHash {
			srcFile, _ := os.Open(path)
			defer srcFile.Close()
			dstFile, _ := os.Create(dstPath)
			defer dstFile.Close()
			io.Copy(dstFile, srcFile)
			fmt.Printf("Synced: %s\n", rel)
		}

		return nil
	})
}
