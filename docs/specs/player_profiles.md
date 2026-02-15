## Feature Description

Ability to view basic player info on dedicated player profile pages in the TUI.
The profiles are loaded dynamically (data queried from NBA API on request), data is not guaranteed persisted across sessions.
Player profiles can be accessed from box scores and league leaders table.

## Implementation plan
-Enable users to navigate to a selected player's profile page from BoxScore and League Leaders view
- Display basic bio info, awards info, current season averages, stat splits by season
- Scrollable view

### Functional scope of Player Profile
- Player profile is divided into three sections - top, mid, bottom
- Top section contains
  - Player Name
  - Player Team
  - Player Position
  - OPTIONAL (needs discussion): player picture. not sure if possible to render in terminal
  - Current season averages on PPG, RPG, APG, SPG, BPG
  - Height
  - Weight
  - Country
  - Draft info (year, selected at)
  - Years of experience
- Mid section contains
  - Awards and honors
- Bottom section contains
  - Last 5 games statline
  - By year regular season splits

### Application
- Change league leaders to also have selectable rows
- When a player row is selected and enter pressed, the corresponding profile is loaded
- Data is queried, deserialized and stored in a <player_id> JSON file, enabling soft caching (until next fs scrub)
- Data is loaded into corresponding structs and presented (BubbleTea) and styled (Lipgloss)

### Data
- player.go structs already contain some of the relevant info (current season stats)
- additional endpoints have to be queried nad new structs created for awards, per season splits
- each player info is stored in <player_id>.json file

### TUI
- player profile page is styled/colored based on their current team colors
- profile page's UI design is proposed by the /frontend-design skill/agent
- QUESTION: is it possible to render player pictures in terminal?
- top section requires different sized fonts and stylings, similar to how player pages are on NBA.com
  - like so https://www.nba.com/player/1628389/bam-adebayo/profile
- mid section (awards) can be a list with some small icons representing the awards
- bottom section can be a simple table

## Implementation strategy
- Start with UI design draft
- Implement section by section, go even smaller if needed (especially on the top section with a lot of elements)
- Work in small increments, frequent small commits
- Write tests to cover new functionality
- ASK when something is ambiguous
- Offer multiple solutions when several feasible options are available
- Use relevant subagents for specific tasks (research, plan, implement, test)

### Testing and validations
- Run full test suite after each committed change
- Validate by running the application and using it

## Out of scope
- Search