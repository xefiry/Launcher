package main

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func GUI_Start(config *Config) {
	const (
		TARGET_FPS   = 60
		WINDOW_WIDTH = 600
	)

	var (
		WINDOW_HEIGHT = config.UI.TitleFontSize +
			2*config.UI.MainFontSize +
			config.Search.MaxResults*config.UI.MainFontSize

		// Colors
		color_main          = rl.SkyBlue
		color_box           = rl.DarkGray
		color_text_area     = rl.RayWhite
		color_font_active   = rl.Black
		color_font_match    = rl.Red
		color_font_inactive = rl.Beige
		color_row_even      = rl.LightGray
		color_row_odd       = rl.Gray
		color_selected      = rl.Green

		// Coordinates
		coord_main rl.Vector2
		coord_text rl.Vector2
		rect_main  rl.Rectangle
		rect_text  rl.Rectangle

		// Working vars
		tmp_color rl.Color
		tmp_text  string

		// Fonts
		font_text  rl.Font
		font_title rl.Font

		// Elements
		input_text         string = ""
		rules_filtered     []*Rule
		strings_filtered   [][]string
		rules_needs_filter bool = true
		active_element     int  = -1 // -1 = typing field, 0 to n = element in list
		nb_displayed_rules int  = 0

		// Misc.
		is_running bool = true
	)

	// Only show warnings and above
	rl.SetTraceLogLevel(rl.LogWarning)

	// Set config and flags
	//rl.SetConfigFlags(rl.FlagWindowTransparent)
	rl.SetWindowState(rl.FlagWindowUndecorated)
	rl.SetTargetFPS(TARGET_FPS)

	// Create new window
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, APP_TITLE)
	defer rl.CloseWindow()

	// Load fonts with right size to avoid blurry text
	// See https://github.com/raysan5/raylib/wiki/Frequently-Asked-Questions#why-is-my-font-blurry
	// Note: I put 256 in last param to make accentuated characters work ... but I don't know why it works
	font_text = rl.LoadFontEx(config.UI.MainFontFile, config.UI.MainFontSize, nil, 256)
	font_title = rl.LoadFontEx(config.UI.TitleFontFile, config.UI.TitleFontSize, nil, 256)

	// Defer the unloading
	defer rl.UnloadFont(font_text)
	defer rl.UnloadFont(font_title)

	for is_running {

		is_running = !rl.WindowShouldClose()

		//---------- Input ----------//

		// Manage adding text
		if key := rl.GetCharPressed(); key != 0 {
			input_char := string(key)

			// add the character to input
			input_text += input_char
			rules_needs_filter = true
		}

		// Manage deleting text
		if rl.IsKeyPressed(rl.KeyBackspace) && len(input_text) > 0 {
			// manage word deletion with Ctrl+Backspace
			if rl.IsKeyDown(rl.KeyLeftControl) {
				// Remove trailing spaces
				input_text = strings.TrimRight(input_text, " ")

				// we search the last space and delete up to it (while keeping it)
				if idx := strings.LastIndex(input_text, " "); idx > 0 {
					input_text = input_text[:idx+1]
				} else { // if none is found, delete everything
					input_text = ""
				}
			} else {
				// delete last character

				// convert string to rune list
				tmp := []rune(input_text)

				// remove the last rune
				tmp = tmp[:len(tmp)-1]

				// put back into input_text as string
				input_text = string(tmp)
			}
			rules_needs_filter = true
		}

		// Get filtered rules (only if it needs to)
		if rules_needs_filter {
			rules_filtered = FilterRules(config.Rules, input_text, config.Search.SearchDescription)
			SortRules(rules_filtered)
			rules_needs_filter = false

			// update number of rules to display
			// the number of rules in the list, if it does not exceed the max number of lines to display
			nb_displayed_rules = min(len(rules_filtered), int(config.Search.MaxResults))

			// reset the active element
			active_element = -1

			// redo the list of display strings
			strings_filtered = [][]string{}
			for i := 0; i < nb_displayed_rules; i++ {
				tmp := rules_filtered[i].GetDisplayStrings(input_text, config.Search.SearchDescription)
				strings_filtered = append(strings_filtered, tmp)
			}
		}

		// Manage navigation
		if rl.IsKeyPressed(rl.KeyUp) {
			if active_element > -1 {
				active_element--
			}
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			if active_element < nb_displayed_rules-1 {
				active_element++
			}
		}
		if rl.IsKeyPressed(rl.KeyHome) {
			// if we are already on the first rule, or if there are no rules
			if active_element == 0 || len(rules_filtered) == 0 {
				// we make the text field as active
				active_element = -1
			} else {
				// otherwise, we make the first rule active
				active_element = 0
			}
		}
		if rl.IsKeyPressed(rl.KeyEnd) {
			// if there are no rules, it will be active_element = 0 - 1 = -1
			active_element = nb_displayed_rules - 1
		}

		// Validation
		if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter) {
			var nb_exec int

			// Only execute if there is at least a rule displayed
			if len(rules_filtered) > 0 {
				// If no rule is selected, use the first one
				if active_element == -1 {
					nb_exec = 0
				} else {
					nb_exec = active_element
				}

				// Execute the rule
				rules_filtered[nb_exec].Execute()

				// Flag the program to exit
				is_running = false
			}
		}

		//---------- Drawing ----------//

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		coord_main = rl.NewVector2(10, 0)
		rect_main = rl.NewRectangle(0, 0, WINDOW_WIDTH, float32(config.UI.TitleFontSize))

		// Title with background
		rl.DrawRectangleRec(rect_main, color_main)
		rl.DrawTextEx(font_title, APP_TITLE, coord_main, float32(config.UI.TitleFontSize), 0, color_font_active)

		// Add version number next to the title
		coord_text = coord_main
		coord_text.X += rl.MeasureTextEx(font_title, APP_TITLE, float32(config.UI.TitleFontSize), 0).X + 20 // Title width + some margin
		coord_text.Y += float32(config.UI.TitleFontSize) / 2                                                // half the height of the title
		rl.DrawTextEx(font_text, APP_VERSION, coord_text, float32(config.UI.MainFontSize), 0, color_font_active)

		// Increase Y for next usages
		coord_main.Y += rect_main.Height
		rect_main.Y += rect_main.Height

		// Text input area background
		rect_main.Height = float32(config.UI.MainFontSize) * 2
		rl.DrawRectangleRec(rect_main, color_main)

		// Text input area coord (with margins each side)
		rect_text = rl.NewRectangle(10, rect_main.Y, rect_main.Width-20, float32(config.UI.MainFontSize)*1.5)
		coord_text = coord_main
		coord_text.X += 10
		coord_text.Y = rect_text.Y + float32(config.UI.MainFontSize)/3
		if len(input_text) == 0 {
			tmp_text = "Enter text here ..."
			tmp_color = color_font_inactive
		} else {
			tmp_text = fmt.Sprintf("%v|", input_text)
			if active_element == -1 {
				tmp_color = color_font_active
			} else {
				tmp_color = color_font_inactive
			}
		}

		rl.DrawRectangleRec(rect_text, color_text_area)
		rl.DrawRectangleLinesEx(rect_text, 2, color_box)
		rl.DrawTextEx(font_text, tmp_text, coord_text, float32(config.UI.MainFontSize), 0, tmp_color)

		// Increase Y for next usages
		coord_main.Y += rect_main.Height
		rect_main.Y += rect_main.Height

		rect_main.Height = float32(config.UI.MainFontSize)
		for i, texts := range strings_filtered {
			if i == active_element {
				tmp_color = color_selected
			} else if i%2 == 0 {
				tmp_color = color_row_even
			} else {
				tmp_color = color_row_odd
			}

			rl.DrawRectangleRec(rect_main, tmp_color)

			coord_text = coord_main
			for j, tmp_text := range texts {
				switch j % 2 {
				case 0:
					tmp_color = color_font_match
				case 1:
					tmp_color = color_font_active
				}
				rl.DrawTextEx(font_text, tmp_text, coord_text, float32(config.UI.MainFontSize), 0, tmp_color)
				coord_text.X += rl.MeasureTextEx(font_text, tmp_text, float32(config.UI.MainFontSize), 0).X
			}

			// Increase Y for next usages
			coord_main.Y += rect_main.Height
			rect_main.Y += rect_main.Height
		}

		// Outline rect
		rect_text = rl.NewRectangle(0, 0, WINDOW_WIDTH, float32(WINDOW_HEIGHT))
		rl.DrawRectangleLinesEx(rect_text, 2, color_box)

		rl.EndDrawing()

		// Uncomment this line to be able to measure performances
		// It will close the program after the first loop
		//is_running = false
	}
}
