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

		file, err := os.Open(filename)
		if os.IsNotExist(err) {
			log.Fatalf("could not find the file `%v`; is the path correct? error: %v", filename, err)
		}
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// calc sha1 hash
		sha := sha1.New()
		if _, err := io.Copy(sha, file); err != nil {
			log.Fatal(err)
		}
		var bs []byte = sha.Sum(nil)
		hexhash := hex.EncodeToString(bs)
		object_dir := OBJECTS_DIR + hexhash[:2] + "/"
		object_name := hexhash[2:]

		file2, err := os.Open(filename)
		if err != nil {
			fmt.Println("error on opening file")
			log.Fatal(err)
		}
		data := make([]byte, 2000000)
		count, err := file2.Read(data)
		if err != nil {
			fmt.Println("error on reading contents of file")
			log.Fatal(err)
		}
		data = data[:count]
		fmt.Printf("read %v bytes\ngot data len: %v\n", count, len(data))

		deflated := new(bytes.Buffer)
		zw := zlib.NewWriter(deflated)
		defer zw.Close()

		n, err := zw.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(deflated)
		fmt.Println(n)
		bb := deflated.Bytes()
		fmt.Printf("%d bytes: %v\n", len(bb), bb)
		fmt.Println(deflated.String())

		err = os.Chdir(GO_GIT_DIR)
		if os.IsNotExist(err) {
			log.Fatalf("could not find the .go-git directory; call `go-git init`. error: %v", err)
		}
		if err != nil {
			log.Fatal(err)
		}

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

		index, err := os.Create(INDEX_PATH)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := index.WriteString(fmt.Sprintf("%v %v\n", hexhash, filename)); err != nil {
			log.Fatal(err)
		}
		defer index.Close()

		// zr, err := zlib.NewReader(deflated)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// io.Copy(os.Stdout, zr)

		// var b bytes.Buffer
		// w := zlib.NewWriter(&b)
		// w.Write([]byte("hello, world\n"))
		// w.Close()
		// r, err := zlib.NewReader(&b)
		// io.Copy(os.Stdout, r)
		// r.Close()

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
