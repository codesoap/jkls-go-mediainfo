// allparams prints all available parameters with their values
// for all streams. Parameters without values will be hidden.
//
// This tool can be used to identify useful parameters for a given
// application.
//
// Usage:
//
//	allparams <path-to-media-file>
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/codesoap/jkls-go-mediainfo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Give exactly one filename as an argument.")
		os.Exit(1)
	}
	mi := mediainfo.New()
	if err := mi.Open(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, "Error: Could not open file:", err)
		os.Exit(1)
	}
	defer mi.Close()
	genParams := make(map[string]string)   // Map from parameter name to full description.
	vidParams := make(map[string]string)   // Map from parameter name to full description.
	audioParams := make(map[string]string) // Map from parameter name to full description.
	txtParams := make(map[string]string)   // Map from parameter name to full description.
	otherParams := make(map[string]string) // Map from parameter name to full description.
	imgParams := make(map[string]string)   // Map from parameter name to full description.
	menuParams := make(map[string]string)  // Map from parameter name to full description.
	newCat := true
	cat := ""
	for line := range strings.Lines(mi.Option("Info_Parameters")) {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			// A new category will start after an empty line.
			newCat = true
			continue
		}
		if newCat {
			cat = fields[0]
		} else {
			switch cat {
			case "General":
				genParams[fields[0]] = strings.Join(fields, " ")
			case "Video":
				vidParams[fields[0]] = strings.Join(fields, " ")
			case "Audio":
				audioParams[fields[0]] = strings.Join(fields, " ")
			case "Text":
				txtParams[fields[0]] = strings.Join(fields, " ")
			case "Other":
				otherParams[fields[0]] = strings.Join(fields, " ")
			case "Image":
				audioParams[fields[0]] = strings.Join(fields, " ")
			case "Menu":
				menuParams[fields[0]] = strings.Join(fields, " ")
			default:
				fmt.Fprintf(os.Stderr, "Error: Unknown category '%s'.", cat)
			}
		}
		newCat = false
	}
	printAllParamsOfStreams(mi, mediainfo.StreamGeneral, genParams)
	printAllParamsOfStreams(mi, mediainfo.StreamVideo, vidParams)
	printAllParamsOfStreams(mi, mediainfo.StreamAudio, audioParams)
	printAllParamsOfStreams(mi, mediainfo.StreamText, txtParams)
	printAllParamsOfStreams(mi, mediainfo.StreamOther, otherParams)
	printAllParamsOfStreams(mi, mediainfo.StreamImage, imgParams)
	printAllParamsOfStreams(mi, mediainfo.StreamMenu, menuParams)
}

func printAllParamsOfStreams(mi *mediainfo.MediaInfo, stream mediainfo.StreamKind, params map[string]string) {
	for i := range mi.Count(stream) {
		switch stream {
		case mediainfo.StreamGeneral:
			fmt.Printf("[General #%d]\n", i+1)
		case mediainfo.StreamVideo:
			fmt.Printf("[Video #%d]\n", i+1)
		case mediainfo.StreamAudio:
			fmt.Printf("[Audio #%d]\n", i+1)
		case mediainfo.StreamText:
			fmt.Printf("[Text #%d]\n", i+1)
		case mediainfo.StreamOther:
			fmt.Printf("[Other #%d]\n", i+1)
		case mediainfo.StreamImage:
			fmt.Printf("[Image #%d]\n", i+1)
		case mediainfo.StreamMenu:
			fmt.Printf("[Menu #%d]\n", i+1)
		}
		keys := make([]string, 0, len(params))
		for key := range params {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, param := range keys {
			if val := mi.Get(stream, i, param); val != "" {
				if param == "Inform" {
					printInform(val)
				} else {
					fmt.Printf("\t%-28s: %s\n", params[param], val)
				}
			}
		}
	}
}

// printInform prints the contents of the special "Inform" value,
// which contains a list of keys and values.
func printInform(val string) {
	fmt.Printf("\t%-28s:\n", "Inform")
	for line := range strings.Lines(val) {
		s := strings.SplitN(line, ":", 2)
		if len(s) != 2 {
			panic("Unexpected format in the 'Inform' parameter.")
		}
		fmt.Printf("\t\t%-28s: %s\n", strings.TrimSpace(s[0]), strings.TrimSpace(s[1]))
	}
}
