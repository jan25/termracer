
## Basic features

- Goal
    - Given random paragraph, type as fast as possible
    - Improved typing skills can be infered from WPM, Accuracy metrics of each race

- Functions:
    - Start new race
        - Ctrl+s starts new race if ui is in stats mode
        - Start a countdown in word view before race begins
        - When race begins, view is focused to word view to allow typing words
        - There is a default maximum time a race can be in progress (2:30), after this time race automatically expires/ends
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
        - We can collapse or expand this view
    - Exit game
        - Ctrl+c exits from racer ui to the command line
        - Asks for confirmation Y/N before exiting
    
```                                                                     
+-------------------------------------------------------------------------+
|+--------------------------------------------------++-------------------+|
||                                                  ||                   ||
||                                                  ||                   ||
||                                                  ||  Stats View       ||
||     Paragraph View                               ||                   ||
||                                                  ||                   ||
||                                                  ||                   ||
||                                                  ||                   ||
||                                                  |+-------------------+|
||                                                  |+-------------------+|
|+--------------------------------------------------+|                   ||
|+--------------------------------------------------+|  Controls View    ||
||     Word View                                    ||                   ||
|+--------------------------------------------------++-------------------+|
+-------------------------------------------------------------------------+         
```   

- UI
    - UI will always have 3 Views/Panes- paragraph, word, stats and controls
    - Paragraph View
        - Shows a paragraph to type if race is in progress
        - Highlights next word to type during a race
            - White bg to indicate next word to type
            - Red bg for the word indicates mistyping in word view
        - Grey out this view if no race is in progress
    - Word View
        - This is a input box to type next target word in the given paragraph
        - Greyed out if no race is in progress
        - Red bg for typed word indicates mistyped, otherwise default bg is applied
        - Greyed out if no race in progress
        - Show countdown before a race begins (1..2..3.. GO!)
    - Stats View
        - This view shows the historical race stats in default mode
        - Shows stats in recent first order. And we will be able to scroll through historical stats
        - Highlights a stat row to show we can scroll
        - In race mode, this view shows auto updating stat(WPM, Accuracy) for race in progress. In addition also shows a timer that counts down from 2:30
    - Controls View
        - Used to show help text that displays the ui controls
        - Ctrl+x can expand/collapse this view
        
- Data directory
    - All data related to termracer is put under $HOME/termracer
    - This directory contains
        - samples
        - log file
        - racehistory file
        
TODO
- Known issues
    - Only one paragraph support on master build
    - Need to get rid of server and generate paragraphs on the fly    

- need to arrange words properly in para view. algorithm to adjust spacing? just look for library
- How about a dashboard to submit book urls to generate samples?



> Credits for the above ASCII art to https://github.com/astashov/tixi