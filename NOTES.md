
## Basic features

- Given random paragraph, type as fast as possible
    - Start new race
    - End race (Stops current race and go to command mode)
    - Command mode
        - When race is not happening (after a race or when race is stopped)
        - See stats (recent races)
    - stats:
        - time to finish typing
        - count of mistakes
    
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
- need module to strip new lines in paragraph files
- need to arrange words properly in para view. algorithm to adjust spacing? just look for library
