/*
Copyright © 2026 manfred6
*/
package cmd

import (
	"fmt"
	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"io"
	"io/fs"
	"math"
	"math/bits"
	"os"
	"path/filepath"
)

var (
	useBit     bool
	useByte    bool
	etype      string
	files      []string
	chunkSize  int64
	fsize      int64
	inputPaths []string
	graphUnit  int
	graphXDesc string
	sectorSize int64
)

const (
	KiB = 1024
	MiB = 1024 * 1024
	GiB = 1024 * 1024 * 1024
	TiB = 1024 * 1024 * 1024 * 1024
)

type entropy struct {
	sliceStart int64
	sliceEnd   int64
	entropy    float64
}

type byteUnit struct {
	size int64
	label string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getFiles(root string) []string {
	// https://medium.com/@andreiboar/listing-files-in-go-33f206c8dcd7
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	check(err)

	return files
}

func bitEntropy(buf []byte) float64 {
	var pX1, sP1, pX0, sP0 float64
	var totalOneBits, bufSize int

	for _, i := range buf {
		totalOneBits += bits.OnesCount8(i)
	}

	pX1 = float64(totalOneBits) / float64(bufSize)
	pX0 = float64(bufSize-totalOneBits) / float64(bufSize)

	if pX1 > 0 {
		sP1 = math.Log2(1 / pX1)
	}
	if pX0 > 0 {
		sP0 = math.Log2(1 / pX0)
	}
	return pX1*sP1 + pX0*sP0
}

func byteEntropy(buf []byte) float64 {
	var counts []int64 = make([]int64, 256)
	var e float64
	for _, byte := range buf {
		counts[byte]++
	}

	for i := range counts {
		var s, p float64
		p = float64(counts[i]) / float64(len(buf))
		if p > 0 {
			s = math.Log2(1 / p)
			e += p * s
		}
	}
	return e
}

func h(file string, etype string, chunkSize int64) []entropy {
	// https://go.dev/tour/methods/21
	var entropyMap []entropy
	var currentSectionStart int64
	var entropyFunc func([]byte) float64

	f, err := os.Open(file)
	check(err)
	defer f.Close()

	buf := make([]byte, chunkSize)

	if etype == "bit" {
		entropyFunc = bitEntropy
	} else {
		entropyFunc = byteEntropy
	}

	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		e := entropy{currentSectionStart, currentSectionStart + int64(n), entropyFunc(buf[:n])}
		entropyMap = append(entropyMap, e)
		currentSectionStart += int64(n)
	}
	return entropyMap
}

func translateBytes(bytes int64) byteUnit {
	var unit byteUnit
	if (bytes / MiB) < 1 {
		unit.size = KiB
		unit.label = "KiB"
	} else if (bytes / GiB) < 1 {
		unit.size = MiB
		unit.label = "MiB"
	} else if (bytes / TiB) < 1 {
		unit.size = GiB
		unit.label = "GiB"
	} else {
		unit.size = TiB
		unit.label = "TiB"
	}
	return unit
}

func graph(h []entropy, fileName string) {
	var values []float64 = make([]float64, len(h))
	for i, _ := range h {
		values[i] = h[i].entropy
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 10
	}

	var graphUnit byteUnit = translateBytes(h[len(h)-1].sliceEnd)
	xStart := float64(h[0].sliceEnd) / float64(graphUnit.size)
	xEnd := float64(h[len(h)-1].sliceEnd) / float64(graphUnit.size)
	graphXDesc = graphUnit.label

	graph := asciigraph.Plot(
		values,
		asciigraph.Height(10),
		asciigraph.Width(width-12),
		asciigraph.XAxisRange(xStart, xEnd),
		asciigraph.XAxisTickCount(len(h)),
		asciigraph.Caption(fmt.Sprintf("Shannon Entropy Graph for %s (X-Units in %s)\n", fileName, graphXDesc)),
	)

	fmt.Println(graph)
}

var rootCmd = &cobra.Command{
	Use:   "h8",
	Short: "Entropy analysis for files and disk images",
	Long: `
██╗  ██╗ █████╗
██║  ██║██╔══██╗
███████║╚█████╔╝
██╔══██║██╔══██╗
██║  ██║╚█████╔╝
╚═╝  ╚═╝ ╚════╝

h8 is a Shannon entropy analysis tool for files and raw disk images.

It maps entropy across configurable regions of data, helping identify
encrypted, compressed, sparse, or otherwise unusual areas.

Designed for disk analysis, ransomware recovery, digital forensics,
and low-level data inspection.
`,

	RunE: func(cmd *cobra.Command, args []string) error {

		if !cmd.Flags().Changed("byte") && !cmd.Flags().Changed("bit") {
			useByte = true
		}

		if useByte {
			etype = "byte"
		} else {
			etype = "bit"
		}

		for i := range inputPaths {
			input, err := os.Open(inputPaths[i])
			check(err)

			info, err := input.Stat()
			check(err)

			if info.IsDir() {
				discoveredFiles := getFiles(inputPaths[i])
				for file := range discoveredFiles {
					files = append(files, discoveredFiles[file])
				}
			} else {
				files = append(files, inputPaths[i])
			}
		}

		fmt.Printf("[i] -> Settings:\n")
		fmt.Printf("  \\__	Mode: %s\n", etype)
		fmt.Printf("  \\__	Logical Sector Size: %d-bytes\n", sectorSize)
		fmt.Printf("  \\__	Entropy Window: %d %s (%d-bytes)\n", chunkSize / translateBytes(chunkSize).size, translateBytes(chunkSize).label, chunkSize)


		for i := range files {
			info, err := os.Stat(files[i])
			check(err)
			fsize = info.Size()
			if fsize == 0 {
				fmt.Printf("[w] -> Skipping %s with size 0\n", files[i])
				continue
			}
			h := h(files[i], etype, chunkSize)
			graph(h, files[i])
			fmt.Printf("[i] -> Total size of %s: %d %s\n", files[i], fsize / translateBytes(fsize).size, translateBytes(chunkSize).label)
			for _, j := range h {
			    fmt.Printf(
					"   \\__ H(X): %.4f bits/byte |> Sector: %d - %d\n",
			        j.entropy,
			        j.sliceStart / sectorSize,
			        j.sliceEnd / sectorSize - 1,
			    )
			}
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringSliceVarP(&inputPaths, "input", "i", nil, "file(s) to extract from. Can be list of files or dir.")
	rootCmd.Flags().Int64VarP(&chunkSize, "size", "s", 1048576, "chunk size (1MiB default)")
	rootCmd.Flags().BoolVar(&useByte, "byte", false, "bytewise shannon entropy")
	rootCmd.Flags().BoolVar(&useBit, "bit", false, "bitwise shannon entropy")
	rootCmd.MarkFlagsMutuallyExclusive("byte", "bit")
	rootCmd.Flags().Int64Var(&sectorSize,"sector-size", 512, "logical sector size in bytes")
}
