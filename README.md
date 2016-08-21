# Golang library to interact with Nuimo devices

This library uses the [currentlabs/ble](https://github.com/currentlabs/ble) library and implements an interaction layer for [Senic Nuimo devices](https://www.senic.com/). Similar to [nathankunicki/nuimojs](https://github.com/nathankunicki/nuimojs) for NodeJS, it was a good inspiration for the library.
 
## Disclaimer
 
At the moment this is a weekend project for me to learn Golang programming. Feel free to suggest changes which change code and interaction to be more #Golang style.

## Example usage*

    go get github.com/tolleiv/nuimo
    # Check out the inputs:
    sudo go run src/github.com/tolleiv/nuimo/examples/inputs/main.go
    # Use the display
    sudo go run src/github.com/tolleiv/nuimo/examples/display/main.go

*this has been tested successfully on Linux only

## License 
 
 MIT License
