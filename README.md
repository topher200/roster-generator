# Roster generator

Roster generator for making balanced teams based on weighted critera.

## Inputs

Takes two csv files as inputs; one with player data, and one of a list of baggages.

  1. The player data is expected to have the following headings:
  - First Name
  - Last Name
  - Gender
  - Balanced Rating

  2. The baggages is a list of "firstname1,lastname1,firstname2,lastname2" baggage
  pairs.

Examples of these can be seen at `sample_players.csv` and `sample_baggages.csv`.

## How it works

roster_generator.go takes in list of ranked players and a list of baggages as
input. It prints to STDOUT the most balanced rosters it can make.

Teams are balanced in the following dimensions:
 - number of baggages satisfied
 - number of players per team
 - number of men/women per team
 - average team rating
 - average male/female rating
 - the standard deviation of each team's ratings (so each team has a balanced "spread")
 - the standard deviation of each team's male/female ratings
 - the standard deviation of each team's top male/female players' ratings

We "balance" a team against the rest in a given category by scoring each team
and trying to minimize the standard deviation (the "distanance apart") of all
those scores. Each dimension is weighted, so some count more or less.

For the implemenation and actual weights used, check out `scoring.go`.

### The genetic algorithm

We have a function that scores a given solution based on the above dimensions.
We make a solution set randomly. take the best solutions as parents for the next
generation. We repeatedly recombine two random solutions to create each new
generation of solutions. We repeat this process a set number of times.

### Development notes

Development can be followed here:
https://trello.com/b/VsN3co1C/smulti-roster-generator

# License

Project copyright topher200@gmail.com. Released under the MIT license.
