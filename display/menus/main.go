package main

import (
	"fmt"
	"strings"
)

type menuContents struct {
	title     *string
	delimiter *string
	lineStyle *string
	maxWidth int 
	elements  []string
}

func splitString(inputString string, maxWidth int, left string, right string, delimiter string) []string {
    words := strings.Fields(inputString)
    var lines []string
    var line string

	innerWidth := maxWidth - len(delimiter) - 2

	if len(inputString) < innerWidth {
		extraSpace := strings.Repeat(" ", innerWidth - len(inputString))
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
					line += fmt.Sprintf("%s- %s", word[start:end - 2], right)
					lines = append(lines, line)
					leftPadding := strings.Repeat(" ", len(delimiter))
					line = fmt.Sprintf("%s  %s%s", left, leftPadding, word[end - 2:])
					start = len(word)
					if idx == len(words) - 1 {
						repeatCount := innerWidth - len(word[end - 2:])
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
			frontPadding := strings.Repeat(" ", len(delimiter) + 2)

            line = fmt.Sprintf("%s%s%s", left, frontPadding, word)
			if idx == len(words) - 1 {
				repeatCount := maxWidth - len(line) + 3
				if repeatCount < 0 {
					repeatCount = 0
				}

				endPadding = strings.Repeat(" ", repeatCount)
				line += fmt.Sprintf("%s%s", endPadding, right)
			}
        } else {
			if len(lines) == 0 && len(line) == 0{
				line += fmt.Sprintf("%s %s %s", left, delimiter, word)
			} else {
				line += fmt.Sprintf(" %s",  word)
			}
			if idx == len(words) - 1 {
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

func printBox(contents menuContents) {
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

	if *contents.lineStyle == "single" {
		topLeft = topLeftSingle
		topRight = topRightSingle
		horizontal = horizontalSingle
		vertical = verticalSingle
		bottomLeft = bottomLeftSingle
		bottomRight = bottomRightSingle
	}

	// first, find the longest line of text
	width := 0
	if contents.title != nil {
		width = len(*contents.title) + 2
	}

	for _, element := range contents.elements {
		lineLength := len(element) + 2
		if lineLength > width {
			width = lineLength
		}
	}
	if width > contents.maxWidth {
		width = contents.maxWidth
	}
	horizontalLine := strings.Repeat(horizontal, width)


	if contents.title != nil {
		titleLength := len(*contents.title) + 1
		repeatCount := width - titleLength
		if repeatCount < 0 {
			repeatCount = 0
		}
		titleLinePadding := strings.Repeat(horizontal, repeatCount - 2)
		fmt.Printf("%s%s %s %s%s\n", topLeft, horizontal, *contents.title, titleLinePadding, topRight)
	} else {
		fmt.Printf("%s", topLeft)
		for range width {
			fmt.Printf("%s", horizontal)
		}
		fmt.Printf("%s\n", topRight)
	}

	var elementsToPrint []string
	for idx, element := range contents.elements {
		var elementsToAdd []string
		if *contents.delimiter == "number" {
			spacer := ""
			if len(contents.elements) > 9 && idx < 9 {
				spacer = " "
			}
			elementsToAdd = splitString(element, width, vertical, vertical, fmt.Sprintf("%d.%s", idx + 1, spacer))
		} else if *contents.delimiter == "letter" {
			letter := rune('a' + idx)
			elementsToAdd = splitString(element, width, vertical, vertical, fmt.Sprintf("%c.", letter))
		} else if contents.delimiter != nil {
			elementsToAdd = splitString(element, width, vertical, vertical, *contents.delimiter)
		}
		elementsToPrint = append(elementsToPrint, elementsToAdd...)
	}

	for _, element := range elementsToPrint {
		fmt.Println(element)
	}

	fmt.Printf("%s%s%s\n", bottomLeft, horizontalLine, bottomRight)
}

func main() {
	title := "Hello world!"
	delimiter := "number"
	// delimiter := "letter"
	// delimiter := "*"
	lineStyle := "double"
	// lineStyle := "single"
	contents := menuContents{
		title:     &title,
		delimiter: &delimiter,
		lineStyle: &lineStyle,
		maxWidth: 60,
		// elements:  []string{"antidisestablishmentarianism",},
		// elements:  []string{"Yo. antidisestablishmentarianism is a pretty long word",},
		elements:  []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "this one is really long but we want the test to wrap around.  Sound good?", "antidisestablishmentarianism", "Yo! antidisestablishmentarianism is a pretty long word, don't you think?", "Yo?"},

	}
	printBox(contents)
}
