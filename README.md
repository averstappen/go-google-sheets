# go-google-sheets

A very crude library based on the google sheets go quickstart guide, for importing into my other projects.

Note that credential file is always named the same, and is looked for in .credentials directory in the current working directory. That's probably not best practice, but helps me to keep things contained in one place.

you will need to:

```
import "github.com/averstappen/go-google-sheets"
```

which will give you

```
resp, err := gsheets.GetValues(sheetId string, readRange string)
// you can use &gsheets.ValueRange{} where you'd normally use &sheets.ValueRange{}
resp, err := gsheets.UpdateValues(sheetId string, writeRange string, values *sheets.ValueRange, valueinputoption string)
```

