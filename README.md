# VRS Tracker

Fetches the latest [Counter-Strike Regional Standings](https://github.com/ValveSoftware/counter-strike_regional_standings) from GitHub, parses the markdown, and stores the results in MongoDB. Intended to run as a cron job. \
Was originally created for use within the [PickemsBot](https://github.com/zacharyab24/PickemsBot) project, so that user's could see the VRS ranking of a team when making decisions

## Configuration

Create a `.env` file in the project root:

```
MONGO_URI=mongodb://localhost:27017
DATABASE_NAME=your_db
COLLECTION_NAME=your_collection
```

## Running

```sh
go run .
```

Or build and run the binary:

```sh
go build -o vrs-tracker .
./vrs-tracker
```

## How it works

On each run the tool checks whether the standings in the database are already current by comparing the stored `standings_date` against the date in the fetched file. If they match, the run exits without touching the database. If the fetched data is newer (or the collection is empty or inconsistent), the collection is dropped and replaced wholesale. This avoids unnecessary writes since VRS standings update at most once a week.

## Known limitations

The GitHub API URL used to fetch standings has the year hardcoded (`2026`). This will break once 2027 data gets uploaded. (well it won't break, it just won't fetch new data anymore). When this happens, the hard coded year can be updated. We can't simply use `time.Now().Year()`, as if Valve don't upload data on Jan 1 2027, we will start getting 404's. Also this means you have to account for timezone dependency as well which is out of scope for something this small. The correct fix would actually be to make an extra api call to GitHub to find the latest year directory that actually contains data, which mirrors how the most recent standings file is located within a year. It's a trivial change but out of scope for now

## Cron example

```
0 * * * * /path/to/vrs-tracker >> /var/log/vrs-tracker.log 2>&1
```
