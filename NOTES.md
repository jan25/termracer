
## Basic features

- Goal
    - Given random paragraph, type as fast as possible
    - Improved typing skills can be infered from WPM, Accuracy metrics of each race
- Functions:
    - Start new race
        - Ctrl+s starts new race if ui is in stats mode
        - Start a countdown in word view before race begins
        - When race begins, view is focused to word view to allow typing words
    - End race
        - Ctrl+e ends a current race and switches ui to stats mode
        - Race auto ends when all words are finished, and ui moves to stats mode
        - Only if race was finished by typing all the words, will the stats for the race get added to stats view
    - Stats mode
        - When no race is in progress, the stats view will show stats from recent races
        - Stats are ordered in recent first order
        - Scroll through historical stats
        - We keep track of Words per minute and Accuracy for a race
    - Controls view
        - Displays help text for ui controls
        - Ctrl+x can collapse or expand this view
    - Exit game
        - Ctrl+c exits from racer ui to the command line
        - Asks for confirmation Y/N before exiting
    
```    
+---------------------------------------------------------------------+             
|+-------------------------------------------++---------------------+ |             
||Paragraph                                  ||Timer                | |             
||                                           ||                     | |             
||                                           ||Stats                | |             
||                                           ||                     | |             
||                                           ||                     | |             
||                                           ||                     | |             
||                                           ||                     | |             
||                                           ||Controls             | |             
||                                           ||                     | |             
|+-------------------------------------------+|                     | |             
|+-------------------------------------------+|                     | |             
|| nextword                                  ||                     | |             
||                                           ||                     | |             
|+-------------------------------------------++---------------------+ |             
+---------------------------------------------------------------------+             
```   

- UI
    - 3 panes: paragraph, type box, Info area
    - Paragraph pane:
        - Shows paragraph to type if race is in progress
        - Shows next word to type during a race
            - Green bg
            - Red bg if word is mistyped
        - Show Lorem lepsum if no race is in progress
    - Type box:
        - Input box to type the next word highlighted in paragraph
        - Green word
        - Red word if word is mistyped
        - Greyed if no race in progress
    - Info box
        - Timer if race in progress
        - Stats when race is not in progress
        - Controls
            - Always visible
            - TODO define controls
        

        
TODO
- need to arrange words properly in para view. algorithm to adjust spacing? just look for library
