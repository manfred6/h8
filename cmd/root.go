/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"io/fs"
	"path/filepath"
	"io"
	"math/bits"
	"math"
	"github.com/spf13/cobra"
)

var inputPaths []string
var files []string

type entropy struct {
	sliceStart int
	sliceEnd   int
	entropy    float64
}

func check(err error) {
	if (err != nil) {
		panic(err)	
	}
}

func getFiles(root string) []string {
	// https://medium.com/@andreiboar/listing-files-in-go-33f206c8dcd7				
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if ! d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	check(err)

	return files
}

func countOneBits(buf []byte) int {
	var total int = 0
	for _, i:= range(buf) {
		total += bits.OnesCount8(i)
	}
	return total
}

func byteEntropy(buf []byte) float64 {
	var counts []int = make([]int, 256)	
	var e float64
	for _, byte := range(buf) {
		counts[byte]++
	}

	for i := range(counts) {
		var s, p float64
		p = float64(counts[i]) / float64(len(buf))
		if p > 0 {
			s = math.Log2(1/p)
		}
		e += p*s
	}
	return e
}

func h(file string, chunkSize int) []entropy {
	// https://go.dev/tour/methods/21
	var entropyMap []entropy
	var currentSectionStart =  0
	f, err := os.Open(file)	
	check(err)
	defer f.Close()

	buf := make([]byte, chunkSize)

	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		//var pX1,sP1,pX0,sP0 float64
		//pX1 = float64(countOneBits(buf[:n])) / float64(n*8)
		//pX0 = float64(len(buf[:n])-countOneBits(buf[:n])) / float64(n*8)
		//if pX1 > 0 {
		//	sP1 = math.Log2(1/pX1)
		//}
		//if pX0 > 0 {
		//	sP0 = math.Log2(1/pX0)
		//}
		//e := entropy{currentSectionStart, currentSectionStart+n, pX1*sP1+pX0*sP0}
		e := entropy{currentSectionStart, currentSectionStart+n, byteEntropy(buf[:n])}
		entropyMap = append(entropyMap, e)
		currentSectionStart += n
	}
	return entropyMap
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "h8",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains`,

	RunE: func(cmd *cobra.Command, args []string) error {

		for i := range(inputPaths) {
			input, err := os.Open(inputPaths[i])
			check(err)

			info, err := input.Stat()
			check(err)

			if info.IsDir() {
				discoveredFiles := getFiles(inputPaths[i])
				for file := range(discoveredFiles) {
					files = append(files, discoveredFiles[file])
				}
			} else {
				files = append(files, inputPaths[i])
			}
		}

		for i:= range(files) {
			fmt.Printf("[i] -> Discovered %s\n", files[i])
			h := h(files[i], 1024*10000)
			for _, j := range h {
			    fmt.Printf(
			        "%8d - %8d | entropy: %.4f bits/byte\n",
			        j.sliceStart,
			        j.sliceEnd,
			        j.entropy,
			    )
			}
		}


		
		return nil
		
	},

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.h8.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringSliceVarP(&inputPaths, "input", "i", nil, "file(s) to extract from. Can be list of files or dir.")

}


