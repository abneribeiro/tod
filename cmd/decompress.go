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

// decompressCmd represents the decompress command
var decompressCmd = &cobra.Command{
    Use:   "decompress [source]",
    Short: "Decompress a zip file",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        source := args[0]
        destination := strings.TrimSuffix(source, filepath.Ext(source))
        err := decompress(source, destination)
        if err != nil {
            fmt.Printf("Failed to decompress: %v\n", err)
        }
    },
}

func init() {
    rootCmd.AddCommand(decompressCmd)
}

func decompress(source, destination string) error {
    reader, err := zip.OpenReader(source)
    if err != nil {
        return err
    }
    defer reader.Close()

    for _, file := range reader.File {
        path := filepath.Join(destination, file.Name)
        if file.FileInfo().IsDir() {
            os.MkdirAll(path, os.ModePerm)
            continue
        }

        if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
            return err
        }

        outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
        if err != nil {
            return err
        }
        defer outFile.Close()

        rc, err := file.Open()
        if err != nil {
            return err
        }
        defer rc.Close()

        _, err = io.Copy(outFile, rc)
        if err != nil {
            return err
        }
    }

    return nil
}