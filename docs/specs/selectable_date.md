# Feature Description

- Change the application so that it does not take any arguments
- The default date used when loading is "today" with TZ as US Eastern
- When the application is running, the date can be changed/selected, but ONLY in the Daily Scores view
- Date selector needs /research and decisions based on feasibility (calendar modal or three validated number fields etc)

## Data
- No major changes in data handling, if another date is selected, the corresponding json file is downloaded to the local FS and deserialized into respective structs

## Application
- When a date is selected, then a new API request is made towards NBA to download the respective JSON

## TUI
- Date selector should be at the top of the page in Daily Scores
- Similarly to how "NBA on YYYY-MM-DD" is displayed on the opening page

# Acceptance Criteria
Given a user wants to see daily nba results
When the user runs the nba-now binary without the "-d YYYY-MM-DD"
Then the application launches with today's dataset (where today is in TZ US East)

Given the application is running
When the user navigates to Daily Scores
and the user selects a date from the date selector
Then the game results (as gamecards) for that date are displayed

# Implementation strategy
- Work in small increments, frequent small commits
- Write tests to cover new functionality
- ASK when something is ambiguous
- Offer multiple solutions when several feasible options are available
- Use relevant subagents for specific tasks (research, plan, implement, test)

Testing and Validation
- Run full test suite after each committed change
- Validate by running the application and using it
