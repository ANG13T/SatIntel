# SatIntel

```
 .       .                   .       .      .     .      .                      .              .
  .    .         .    .            .     ______ 
          .             .               ////////                .         .      .       .       .          .
        .    .   ________   .  .      /////////     .    .
   .            |.____.  /\         /////////    .                      .               .               .
               //      \/  |\     /////////
       .     //          \ |  \ /////////         _______ _______ _______ _____ __   _ _______ _______       .
            ||           | |  ///////// .     .   |______ |_____|    |      |   | \  |    |    |______ |     
  .         ||           | |//  /////             ______| |     |    |    __|__ |  \_|    |    |______ |_____  .       
     .       \\         / //     \/   .                    
               \\.___./ //\      ,_\     .     .                                                            .
  .           .    \ //////\   /    \                 .    .      Satellite OSINT CLI Tool          .            .
               .   ///////// \|      |    .                    
      .           ///////// .  \ __ /          .               Made by Angelina Tsuboi (G4LXY)              .
                 /////////                              .               .                   .
        .   .   /////////     .     .                           .                   .                   .     .
 .             --------   .                  ..             .               .                .
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
Make an account at [**Space Track**](space-track.org) save username and password.

Create an account at [**N2YO**](n2yo.com) and save API key.

Update `main.go` to have proper credentials
```
os.Setenv("SPACE_TRACK_USERNAME", "username")
os.Setenv("SPACE_TRACK_PASSWORD", "password")
os.Setenv("N2YO_API_KEY", "api-key")
```

To build from source, you will need Go installed.

```bash
$ export GO111MODULE=on 
$ go get ./...
$ go run main.go
```

### APIs Used
- [Space Track](space-track.org): Retrieve Satellite Catalog and TLE Information
- [N2YO](n2yo.com/api): Retrieve Passes Predictions

### Upcoming Features
+ [ ] Map Layout of Satellite Positioning
+ [ ] Including the [SGP4 Algorithm](joshuaferrara/go-satellite) for Position Prediction 

