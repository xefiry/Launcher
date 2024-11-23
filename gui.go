package main

import (
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

		// Sizes
		title_size = float32(config.UI.TitleFontSize)
		main_size  = float32(config.UI.MainFontSize)

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
		coord_main  rl.Vector2
		coord_text  rl.Vector2
		rect_main   rl.Rectangle
		rect_text   rl.Rectangle
		rect_scroll rl.Rectangle

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
		nb_rules           int
		first_display_rule int
		last_display_rule  int

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

			// update number of rules$
			nb_rules = len(rules_filtered)

			// reset the active element
			active_element = -1

			// redo the list of display strings
			strings_filtered = [][]string{}
			for i := 0; i < nb_rules; i++ {
				tmp := rules_filtered[i].GetDisplayStrings(input_text, config.Search.SearchDescription)
				strings_filtered = append(strings_filtered, tmp)
			}

			// mark as filtered
			rules_needs_filter = false

			first_display_rule = 0

			// The last rule to display is the minimum between the last rule in the list
			// and the first rule to display + the number of rules to display
			// -1 to get the id of the last rule to display
			last_display_rule = min(nb_rules, first_display_rule+int(config.Search.MaxResults)) - 1
		}

		// Manage navigation
		if rl.IsKeyPressed(rl.KeyUp) {
			if active_element > 0 {
				active_element--

				// if we are not on the first element of the list
				// and we are on the first element displayed
				if active_element > 0 && active_element == first_display_rule {
					// we increase the first rule displayed
					first_display_rule--

					// the last is the first + the number of rules to display, or the number of rules in the list (if not enough)
					last_display_rule = min(first_display_rule+int(config.Search.MaxResults), nb_rules) - 1
				}
			}
		}

		if rl.IsKeyPressed(rl.KeyDown) {
			if active_element < nb_rules-1 {
				active_element++

				// if we are not on the last element of the list
				// and we are on the last element displayed
				if active_element < nb_rules-1 && active_element == last_display_rule {
					// we increase the first rule displayed
					first_display_rule++

					// the last is the first + the number of rules to display, or the number of rules in the list (if not enough)
					last_display_rule = min(first_display_rule+int(config.Search.MaxResults), nb_rules) - 1
				}
			}
		}

		if rl.IsKeyPressed(rl.KeyHome) {
			active_element = 0

			// reset the first rule to display as the first rule of the list
			first_display_rule = 0
			// the last is the first + the number of rules to display, or the number of rules in the list (if not enough)
			last_display_rule = min(first_display_rule+int(config.Search.MaxResults), nb_rules) - 1
		}

		if rl.IsKeyPressed(rl.KeyEnd) {
			// if there are no rules, it will be active_element = 0 - 1 = -1
			active_element = nb_rules - 1

			// the last rule to display becomes the last rule
			last_display_rule = nb_rules - 1

			// and the first is last - nb of displayed rules, but cannot be negative
			first_display_rule = max(0, last_display_rule-int(config.Search.MaxResults)+1)
		}

		// Validation
		if rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter) {
			var nb_exec int

			// Only execute if there is at least a rule displayed
			if nb_rules > 0 {
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
		rect_main = rl.NewRectangle(0, 0, WINDOW_WIDTH, title_size)

		// Title with background
		rl.DrawRectangleRec(rect_main, color_main)
		rl.DrawTextEx(font_title, APP_TITLE, coord_main, title_size, 0, color_font_active)

		// Add version number next to the title
		coord_text = coord_main
		coord_text.X += rl.MeasureTextEx(font_title, APP_TITLE, title_size, 0).X + 20 // Title width + some margin
		coord_text.Y += title_size / 2                                                // half the height of the title
		rl.DrawTextEx(font_text, APP_VERSION, coord_text, main_size, 0, color_font_active)

		// Increase Y for next usages
		coord_main.Y += rect_main.Height
		rect_main.Y += rect_main.Height

		// Text input area background
		rect_main.Height = main_size * 2
		rl.DrawRectangleRec(rect_main, color_main)

		// Text input area coord (with margins each side)
		rect_text = rl.NewRectangle(10, rect_main.Y, rect_main.Width-20, main_size*1.5)
		coord_text = coord_main
		coord_text.X += 10
		coord_text.Y = rect_text.Y + main_size/3
		if len(input_text) == 0 {
			tmp_text = "Enter text here ..."
			tmp_color = color_font_inactive
		} else {
			tmp_text = input_text
			tmp_color = color_font_active
		}

		rl.DrawRectangleRec(rect_text, color_text_area)
		rl.DrawRectangleLinesEx(rect_text, 1, color_box)
		rl.DrawTextEx(font_text, tmp_text, coord_text, main_size, 0, tmp_color)

		// Increase Y for next usages
		coord_main.Y += rect_main.Height
		rect_main.Y += rect_main.Height

		// Scroll bar management
		if nb_rules > int(config.Search.MaxResults) {
			// bar width is relative to font size
			bar_width := main_size / 2

			// reduce the size of the background for rules to leave space for the scrolling bar
			rect_main.Width -= bar_width

			// height of the whole scroll bar space
			height := main_size * float32(config.Search.MaxResults)

			// height of the actual bar, proportionnal with the number of rules in the list
			bar_height := height * float32(config.Search.MaxResults) / float32(nb_rules)

			// calculate the space the bar can move
			vertical_space := height - bar_height

			// calculate the first rule displayed when at the bottom of the rules
			bottom := nb_rules - int(config.Search.MaxResults)

			// calculate the vertical offset for the scroll bar
			//              start position + availlable space for the bar to move proportioned
			vertical_offset := rect_main.Y + vertical_space*float32(first_display_rule)/float32(bottom)

			// create the scrolling bar to fill the left side
			rect_scroll = rl.NewRectangle(rect_main.Width, vertical_offset, bar_width, bar_height)

			rl.DrawRectangleRounded(rect_scroll, main_size/2, 5, color_box)
		}

		rect_main.Height = main_size
		for i, texts := range strings_filtered[first_display_rule : last_display_rule+1] {
			i += first_display_rule
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
				rl.DrawTextEx(font_text, tmp_text, coord_text, main_size, 0, tmp_color)
				coord_text.X += rl.MeasureTextEx(font_text, tmp_text, main_size, 0).X
			}

			// Increase Y for next usages
			coord_main.Y += rect_main.Height
			rect_main.Y += rect_main.Height
		}

		// Outline rect
		rect_text = rl.NewRectangle(0, 0, WINDOW_WIDTH, float32(WINDOW_HEIGHT))
		rl.DrawRectangleLinesEx(rect_text, 1, color_box)

		rl.EndDrawing()

		// Uncomment this line to be able to measure performances
		// It will close the program after the first loop
		//is_running = false
	}
}
