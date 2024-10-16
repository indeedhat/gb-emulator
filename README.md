# Gameboy emulator written in go

Well, it runs pokemon red

## Keymap (Sorry its setup for dvorak)
```
Up     -> Comma  
Right  -> E  
Down   -> O  
Left   -> A  
A      -> Enter  
B      -> J  
Start  -> Period  
Select -> Quote  
```

## Usage
```
make build
./build/gb-emu ./roms/tetris.gb
```


## TODO (Emulation)
- [x] Bank switching
- [ ] allow multiple state slots
- [ ] auto save state
- [ ] fix screen flicker
- [ ] Sound
- [ ] Window scroll (still work to do)

## TODO (User Interface)
- [ ] file menu
    - [x] load rom
    - [x] save state
    - [x] load state
    - [ ] remap key binds
- [ ] debug windows
    - [ ] tile maps
    - [ ] Cpu registers
    - [ ] Input

