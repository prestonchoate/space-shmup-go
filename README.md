# Untitled Space Shooter

This game is a passion project I'm using to learn Go and game dev more intimately. It's a classic space shoot em up style game with deep inspiration from games like Galaga, and some roguelike elements thrown in to keep things interesting. This project utilizes the incredible efforts of Raysan5's [Raylib](https://github.com/raysan5/raylib) and the Go bindings by [gen2brain](https://github.com/raysan5/raylib)

## (Eventual) Features

- **Classic Arcade Style Gameplay**: Fast and difficult to put down gameplay that is simple to grasp, but challenging to master :white_check_mark: 
- **Diverse Enemies**: Multiple enemy types with varied movement and attack patterns :black_square_button:
- **Power-Ups**: Randomized power ups and upgrades that will allow you to create powerful and game changing combinations :black_square_button:
- **Challenging Levels**: Better get those upgrades in place because as the waves keep coming the number of enemies increases along with their difficulty :black_square_button:
- **Online Leaderboard**: Get an incredibly high score one run? Upload it to the online leaderboard and re-live those glory days at the arcade and chasing the top spot :black_square_button:
- **Configurable Controls**: Don't like the default controls? Rebind them to any key you wish. Custom settings for music and SFX volume as well :white_check_mark:

## Installation

Grab the latest release for your OS (only Windows and Linux are currently supported) and run the binary. I may implement an installation wizard later but for now it's as simple as download and run.

## Building From Source

1. **Clone the repository**:
    ```bash
        git clone https://github.com/prestonchoate/space-shmup-go.git
    ```
2. **Navigate to the Project Directory**:
    ```bash
        cd space-shmup-go
    ```
3. **Install Dependencies (Ensure Go is installed)**:
    ```bash
        go mod tidy
    ```
4. **Build the project**:
    
    Linux
    ```bash
        CGO_ENABLED=1 go build -o space-shmup-go
    ```
    Windows
    ```bash
        CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o build/space-shmup-windows-amd64.exe
    ```
5. **Run the game**:
    ```bash
        ./space-shmup-go
    ```

## Controls
- **Movement**: Defaults to `WASD` but configurable through the in game settings menu
- **Fire Weapon**: Defaults to `SPACE` but is also configurable
- **Pause**: `ESC` to pause or resume

## Contributing
Contributions are welcome! If you'd like to contribute, please fork the repository and use a feature branch. Pull requests are warmly welcomed.
