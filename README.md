# SatIntel

```
 .       .                   .       .      .     .      .                      .              .
    .           .            .     ________
                  .               /////////                .         .      .       .       .          .
        .   ________   .  .      /////////     .    .
           |.____.  /\         /////////    .                      .               .               .
  .       //      \/  |\     /////////
         //          \ |  \ /////////         _______ _______ _______ _____ __   _ _______ _______       .
        ||           | |  ///////// .     .   |______ |_____|    |      |   | \  |    |    |______ |
   .    ||           | |//  /////             ______| |     |    |    __|__ |  \_|    |    |______ |_____  .
        \\         / //     \/   .
          \\.___./ //\      ,_\     .     .                                                            .
  .       .    \ //////\   /    \                 .    .      Satellite OSINT CLI Tool          .            .
          .    ///////// \|      |    .
       .      ///////// .  \ __ /          .               Made by Angelina Tsuboi (G4LXY)              .
 .           /////////                              .               .                   .
   .   .    /////////     .     .                           .                   .                   .     .
           --------   .                  ..             .               .                .
    .        .         .                       .                                 .                .
```

#### SatIntel is a OSINT tool for satellite reconnaissance made with Golang. The tool can extract satellite telemetry, receive orbital predictions, and parse TLEs.

### Features
- Satellite Catalog Retrieval from NORAD ID or Selection Menu
- Display Satellite Telemetry
- Visual and Radio Orbital Predictions
- Parse Two Line Elements (TLE)

### Preview
<img src="https://github.com/ANG13T/SatIntel/blob/main/assets/image.png" alt="SatIntel Image" width="600"/>

### Usage
Make an account at [**Space Track**](https://space-track.org) save username and password.

Create an account at [**N2YO**](https://n2yo.com) and save API key.

The CLI will prompt for both Space Track and N2YO credentials if none are present in your environmental variables. To export your credentials enter the following commands:
```bash
$ export SPACE_TRACK_USERNAME="YOUR_USER_NAME"
$ export SPACE_TRACK_PASSWORD="YOUR_PASSWORD"
$ export N2YO_API_KEY="YOUR_API_KEY"
```

To build from source, you will need Go installed.

```bash
$ export GO111MODULE=on
$ go get ./...
$ go run main.go
```

### APIs Used
- [Space Track](https://space-track.org): Retrieve Satellite Catalog and TLE Information
- [N2YO](https://n2yo.com/api): Retrieve Passes Predictions

### Satellite OSINT Explained
To get a general overview of the underlying concepts of this tool, [read this article](https://medium.com/@angelinatsuboi/satellite-osint-space-based-intelligence-in-cybersecurity-e87f9dca4d81).

### Upcoming Features
+ [ ] Map Layout of Satellite Positioning
+ [ ] Including the [SGP4 Algorithm](joshuaferrara/go-satellite) for Position Prediction

