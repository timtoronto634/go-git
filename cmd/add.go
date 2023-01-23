/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const MAX_FILE_BYTES = 2000000

const (
	GO_GIT_DIR  = ".go-git/"
	OBJECTS_DIR = "objects/"
	INDEX_PATH  = "index"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Args:  cobra.MinimumNArgs(1),
	Short: "Stage files specified",
	Long: `Stage the given flies for commit
	only supports simle staging, and files must be specified by name
	'.' or '*' is not supported
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatal("currently supports only 1 file")
		}
		filename := args[0]

		// calc sha1 hash
		file_for_hash, err := os.Open(filename)
		if os.IsNotExist(err) {
			log.Fatalf("could not find the file `%v`; is the path correct? error: %v", filename, err)
		}
		if err != nil {
			log.Fatal(err)
		}
		defer file_for_hash.Close()

		sha := sha1.New()
		if _, err := io.Copy(sha, file_for_hash); err != nil {
			log.Fatal(err)
		}
		var bs []byte = sha.Sum(nil)
		hexhash := hex.EncodeToString(bs)
		object_dir := OBJECTS_DIR + hexhash[:2] + "/"
		object_name := hexhash[2:]

		// compress contents
		file_for_compress, err := os.Open(filename)
		if err != nil {
			fmt.Println("error on opening file")
			log.Fatal(err)
		}
		data := make([]byte, MAX_FILE_BYTES)
		count, err := file_for_compress.Read(data)
		if err != nil {
			log.Fatalf("error on reading contents of file. error: %v", err)
		}
		if count == MAX_FILE_BYTES {
			log.Fatalf("specified file is too big: %v", filename)
		}
		data = data[:count]

		var deflated bytes.Buffer
		zw := zlib.NewWriter(&deflated)
		_, err = zw.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		zw.Close()

		// cd
		err = os.Chdir(GO_GIT_DIR)
		if os.IsNotExist(err) {
			log.Fatalf("could not find the .go-git directory; call `go-git init`. error: %v", err)
		}
		if err != nil {
			log.Fatal(err)
		}

		// write object(compressed)
		if err := os.Mkdir(object_dir, os.ModePerm); err != nil && (!os.IsExist(err)) {
			log.Fatal(err)
		}
		blob, err := os.Create(object_dir + object_name)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := blob.WriteString(deflated.String()); err != nil {
			log.Fatal(err)
		}
		defer blob.Close()

		// write .go-git/index
		index, err := os.Create(INDEX_PATH)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := index.WriteString(fmt.Sprintf("%v %v\n", hexhash, filename)); err != nil {
			log.Fatal(err)
		}
		defer index.Close()

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
