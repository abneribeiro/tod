package cmd

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
)

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
    Use:   "compress [source]",
    Short: "Compress a file or directory",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        source := args[0]
        destination := source + ".zip"
        err := compress(source, destination)
        if err != nil {
            fmt.Printf("Failed to compress: %v\n", err)
        }
    },
}

func init() {
    rootCmd.AddCommand(compressCmd)
}

func compress(source, destination string) error {
    zipfile, err := os.Create(destination)
    if err != nil {
        return err
    }
    defer zipfile.Close()

    archive := zip.NewWriter(zipfile)
    defer archive.Close()

    info, err := os.Stat(source)
    if err != nil {
        return err
    }

    var baseDir string
    if info.IsDir() {
        baseDir = filepath.Base(source)
    }

    return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        header, err := zip.FileInfoHeader(info)
        if err != nil {
            return err
        }

        if baseDir != "" {
            header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
        }

        if info.IsDir() {
            header.Name += "/"
        } else {
            header.Method = zip.Deflate
        }

        writer, err := archive.CreateHeader(header)
        if err != nil {
            return err
        }

        if !info.IsDir() {
            file, err := os.Open(path)
            if err != nil {
                return err
            }
            defer file.Close()
            _, err = io.Copy(writer, file)
            if err != nil {
                return err
            }
        }

        return nil
    })
}