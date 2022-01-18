package handler

import (
	"fmt"
	"testing"
)

func TestColor(t *testing.T) {
	fmt.Printf(Bold + "bold" + Reset + "\n")
	fmt.Printf(Faint + "faint" + Reset + "\n")
	fmt.Printf(Underlined + "underlined" + Reset + "\n")
	fmt.Printf(Blink + "blink" + Reset + "\n")

	fmt.Printf(Black + "black" + Reset + "\n")
	fmt.Printf(Red + "red" + Reset + "\n")
	fmt.Printf(Green + "green" + Reset + "\n")
	fmt.Printf(Yellow + "yellow" + Reset + "\n")
	fmt.Printf(Blue + "blue" + Reset + "\n")
	fmt.Printf(Purple + "purple" + Reset + "\n")
	fmt.Printf(Cyan + "cyan" + Reset + "\n")
	fmt.Printf(Gray + "gray" + Reset + "\n")

	fmt.Printf(White + "white logging" + Reset + "\n")

	fmt.Printf(GrayBlack + "gray background with black foreground" + Reset + "\n")
	fmt.Printf(CyanRed + "cyan background with red foreground" + Reset + "\n")
	fmt.Printf(PurpleGreen + "purple background with green foreground" + Reset + "\n")
	fmt.Printf(BlueYellow + "blue background with yellow foreground" + Reset + "\n")
	fmt.Printf(YellowBlue + "yellow background with blue foreground" + Reset + "\n")
	fmt.Printf(GreenPurple + "green background with purple foreground" + Reset + "\n")
	fmt.Printf(RedCyan + "red background with cyan foreground" + Reset + "\n")
	fmt.Printf(BlackGray + "black background with gray foreground" + Reset + "\n")

	fmt.Printf(BlackWhite + "black background with white foreground" + Reset + "\n")

}
