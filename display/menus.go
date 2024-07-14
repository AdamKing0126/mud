package display

import (
	"fmt"
	"strings"
)

type MenuContents struct {
	Title     *string
	Delimiter *string
	LineStyle *string
	MaxWidth  int
	Elements  []string
}

func splitString(inputString string, maxWidth int, left string, right string, delimiter string) []string {
	words := strings.Fields(inputString)
	var lines []string
	var line string

	innerWidth := maxWidth - len(delimiter) - 2

	if len(inputString) < innerWidth {
		extraSpace := strings.Repeat(" ", innerWidth-len(inputString)+1)
		lines = append(lines, fmt.Sprintf("%s %s %s%s%s", left, delimiter, inputString, extraSpace, right))
		return lines
	}

	for idx, word := range words {
		if len(word) >= innerWidth {
			// If the word is longer than the maximum width, split it into substrings
			for start := 0; start < len(word); start += (len(line) + 1) {
				// for start := 0; start < len(word); start += innerWidth {
				end := start + innerWidth
				if len(line) != 0 {
					line += " "
					end = end - len(line) + 5 + len(delimiter)
				} else {
					line = fmt.Sprintf("%s %s ", left, delimiter)
				}
				if end > len(word) {
					end = len(word)
				}
				if len(lines) == 0 {
					line += fmt.Sprintf("%s- %s", word[start:end-2], right)
					lines = append(lines, line)
					leftPadding := strings.Repeat(" ", len(delimiter))
					line = fmt.Sprintf("%s  %s%s", left, leftPadding, word[end-2:])
					start = len(word)
					if idx == len(words)-1 {
						repeatCount := innerWidth - len(word[end-2:])
						if repeatCount < 0 {
							// TODO: edge case where a very long word must be broken up into 3+ lines
							// ex: "antidisestablishmentarianism" with maxWidth 10
							// "antidises-"
							// "tablishme-"
							// "ntarianis-"
							// "m"
							repeatCount = 0
						}
						rightPadding := strings.Repeat(" ", repeatCount)
						line += fmt.Sprintf("%s%s", rightPadding, right)
					}
				} else {
					lines = append(lines, fmt.Sprintf("%s %s %s", left, word[start:end], right))
				}
			}
		} else if len(line)+len(word) > maxWidth {
			repeatCount := maxWidth - len(line) + 3
			if repeatCount < 0 {
				repeatCount = 0
			}
			endPadding := strings.Repeat(" ", repeatCount)
			line += fmt.Sprintf("%s%s", endPadding, right)
			lines = append(lines, line)
			frontPadding := strings.Repeat(" ", len(delimiter)+2)

			line = fmt.Sprintf("%s%s%s", left, frontPadding, word)
			if idx == len(words)-1 {
				repeatCount := maxWidth - len(line) + 3
				if repeatCount < 0 {
					repeatCount = 0
				}

				endPadding = strings.Repeat(" ", repeatCount)
				line += fmt.Sprintf("%s%s", endPadding, right)
			}
		} else {
			if len(lines) == 0 && len(line) == 0 {
				line += fmt.Sprintf("%s %s %s", left, delimiter, word)
			} else {
				line += fmt.Sprintf(" %s", word)
			}
			if idx == len(words)-1 {
				repeatCount := maxWidth - len(line) + 3
				if repeatCount < 0 {
					repeatCount = 0
				}
				padding := strings.Repeat(" ", repeatCount)
				line += fmt.Sprintf("%s%s", padding, right)
			}
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
}

func PrintMenu(player ProfilePlayer, contents MenuContents) {
	Newline(player)
	topLeftDbl := string('╔')
	topLeftSingle := string('┌')
	topLeft := topLeftDbl

	topRightDbl := string('╗')
	topRightSingle := string('┐')
	topRight := topRightDbl

	horizontalDbl := string('═')
	horizontalSingle := string('─')
	horizontal := horizontalDbl

	verticalDbl := string('║')
	verticalSingle := string('│')
	vertical := verticalDbl

	bottomLeftDbl := string('╚')
	bottomLeftSingle := string('└')
	bottomLeft := bottomLeftDbl

	bottomRightDbl := string('╝')
	bottomRightSingle := string('┘')
	bottomRight := bottomRightDbl

	if *contents.LineStyle == "single" {
		topLeft = topLeftSingle
		topRight = topRightSingle
		horizontal = horizontalSingle
		vertical = verticalSingle
		bottomLeft = bottomLeftSingle
		bottomRight = bottomRightSingle
	}

	// first, find the longest line of text
	width := 0
	if contents.Title != nil {
		width = len(*contents.Title) + 2
	}

	for _, element := range contents.Elements {
		lineLength := len(element) + 2
		if lineLength > width {
			width = lineLength
		}
	}
	if width > contents.MaxWidth {
		width = contents.MaxWidth
	}
	horizontalLine := strings.Repeat(horizontal, width+1)

	if contents.Title != nil {
		titleLength := len(*contents.Title) + 2
		repeatCount := width - titleLength
		if repeatCount < 0 {
			repeatCount = 0
		}
		titleLinePadding := strings.Repeat(horizontal, repeatCount)
		msg := fmt.Sprintf("%s%s %s %s%s\n", topLeft, horizontal, *contents.Title, titleLinePadding, topRight)
		PrintWithColor(player, msg, "primary")
	} else {
		PrintWithColor(player, topLeft, "primary")
		for range width {
			PrintWithColor(player, horizontal, "primary")
		}
		PrintWithColor(player, fmt.Sprintf("%s\n", topRight), "primary")
	}

	var elementsToPrint []string
	for idx, element := range contents.Elements {
		var elementsToAdd []string
		if contents.Delimiter == nil || *contents.Delimiter == "none" {
			elementsToAdd = splitString(element, width, vertical, vertical, "")
		} else if *contents.Delimiter == "number" {
			spacer := ""
			if len(contents.Elements) > 9 && idx < 9 {
				spacer = " "
			}
			elementsToAdd = splitString(element, width, vertical, vertical, fmt.Sprintf("%d.%s", idx+1, spacer))
		} else if *contents.Delimiter == "letter" {
			letter := rune('a' + idx)
			elementsToAdd = splitString(element, width, vertical, vertical, fmt.Sprintf("%c.", letter))
		} else {
			elementsToAdd = splitString(element, width, vertical, vertical, *contents.Delimiter)
		}
		elementsToPrint = append(elementsToPrint, elementsToAdd...)
	}

	for _, element := range elementsToPrint {
		PrintWithColor(player, element+"\n", "primary")
	}

	PrintWithColor(player, fmt.Sprintf("%s%s%s\n", bottomLeft, horizontalLine, bottomRight), "primary")
}

// func main() {
// 	title := "Hello world!"
// 	delimiter := "number"
// 	// delimiter := "letter"
// 	// delimiter := "*"
// 	lineStyle := "double"
// 	// lineStyle := "single"
// 	contents := MenuContents{
// 		title:     &title,
// 		delimiter: &delimiter,
// 		lineStyle: &lineStyle,
// 		maxWidth:  60,
// 		// elements:  []string{"antidisestablishmentarianism",},
// 		// elements:  []string{"Yo. antidisestablishmentarianism is a pretty long word",},
// 		elements: []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "this one is really long but we want the test to wrap around.  Sound good?", "antidisestablishmentarianism", "Yo! antidisestablishmentarianism is a pretty long word, don't you think?", "Yo?"},
// 	}
// 	PrintMenu(contents)
// }
