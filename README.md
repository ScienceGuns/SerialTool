Mad Scientist Research Serial Tool
==================================
This tool can encode and decode the Universal Serial Number ("USN") format found on the products made by [Mad Scientist Research LLC](https://scienceguns.com/)

Settings File
-------------
To generate USNs, you need to point this tool at a MongoDB server which will track previously generated ones.
The MongoDB can be reset whenever you are done for the day. It's only purpose is to keep track of the day's counts.
**MongoDB MUST not be relied on for any A&D inventory purposes!**

The config file belongs in `~/.config/scienceguns/serialtool.json`
Currently supported values:
```json
{
  "data_url": "https://scienceguns.com/data/",
  "mongodb": {
    "uri": "mongodb://drkrieger:guest@localhost:27017",
    "database": "firearms",
    "collection": "Serials"
  }
}
```

You can make sure you have the latest product model list via `serialtool update` using the data URL in the example JSON above.

The "Universal Serial Number" Pattern
-------------------------------------
The pattern itself is pretty easy to remember. Serial Tool is meant to make things faster and avoid mistakes.

The format is `MMTYYYSSCCCC`

1. `MM` - Denotes the model of product using the `[0-Z][0-Z]` format (ex: `00`, `TJ`, or `1F`)
2. `T` - Denotes the type of firearm the weapon was when initially manufactured.
    1. `R` - Rifle
    2. `P` - Pistol
    3. `F` - Frame
    4. `G` - Shotgun
    5. `B` - Short Barreled Rifle
    6. `H` Short Barreled Shotgun
    7. `S` - Suppressor
    8. `A` - Any Other Weapon
    9. `M` - Machine Gun
3. `YYY` - Denotes the year.
    1. The first `Y` counts up sequentially from 1 through Z where `1=20` and `A=29`
    2. The remaining `YY` is a normal 2-digit number meaning `125=2025`, `201=2101`, and `A21=2921`
4. `SS` - Denotes an internal supply chain logistics code
5. `CCCC` - Denotes the sequentially increasing unit count for that weapon model manufactured that day (or month for marking variance weapons). This number starts at `0001` and is padded with zeroes.

Two examples:

1. `TIF125MS0023` would be a `The Isotope` manufactured as a Frame in 2025. It would be the 23rd one made that year.
2. `TKR125MS0023` would be a `The Kinetic` manufactured as a Rifle in 2025. It would be the 23rd one made that year.

Why Make The USN Pattern
------------------------
Mad Scientist Research is founded by a science loving engineer that sees value in thinking ahead on things like this. Having
a single serial number format across all product lines makes record keeping easier for inventory, warranty work, determining what a thing is or was at time of manufacture,
and worst case scenario, for tracing.
