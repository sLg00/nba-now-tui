<h1>NBA Now</h1>

NBA NOW is a simple terminal application built with Go. 
The TUI is built using [Charm](https://charm.sh)'s suite of libraries and leveraging the Bubble Tea framework.

This is a personal hobby project not to be taken seriously. The current version is lacking a lot of features i want to ultimately
have - soon:tm:

Big thanks to [Swar's NBA API](https://github.com/swar/nba_api) project's 
extensive documentation of the NBA APIs. This project would not be possible
without all that work!


<h3>Build</h3>

Just 'make' in the root directory runs all the test and builds executables for Linux, Mac and Win.

<h4>Run</h4>
Execute **./<binary> -d YYYY-MM-DD** to launch the app. Due to the timezone differences between the States
and the rest of the world, i changed the app logic to take in a specific date. The date
denotes local time when the games occurred.

<h3>Available Features</h3>

* Daily game results
* Box scores
* League leaders
* Season standings


<h3>Shit still missing, yo</h3>

* _(feature)_ Team profile pages (these will be persistent, and updated with stats)
* _(feature)_ Player profile pages (on-demand, when a player is selected from an existing view)
* _(feature)_ Playoff bracket
* _(feature)_ News headlines
* _(feature)_ Support for PreSeason & Playoffs
* _(basics)_ Sorting on tables
* _(FE)_ Consistent & distinct styling of elements
* _(QA)_ Bunch of tests still missing, especially on TUI logic
* _(FE)_ Error modals and loading bar


<h4>Not gonna happen</h4>

* Search


<h3>Tech details</h3> 
Concurrently querying NBA APIs and displaying the results in the terminal (revolutionary, i know).
Results are stored in json files in a designated folder **(~/.config/nba-tui)**. Only the necessary files are downloaded
and parsed. On app initiation the daily scores, league leaders and season standings are queried. Once daily view
is opened, the files for the box scores are downloaded and parsed.

Logs are written to a dedicated log file **(~/.config/nba-tui/logs/appLog.log)**. All downloaded json files, older than 48 hours
are deleted on app launch to avoid cluttering the filesystem.

Why filesystem and not an sqlite db? The database already exists on NBA's side, so this is just about the terminal client and not
persisting a ton of data.

The app is also not fully leveraging Bubble Tea, since i wanted to keep the deep backend logic
isolated from frontend in the early days, to make sure it's easy to jump off of Bubble Tea, if i don't like it. 
But i do like it, so i have bit of refactoring to do at some point.

**Currently only tested on Linux, because that's where i use it.** /shrug