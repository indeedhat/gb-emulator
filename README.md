# Gameboy emulator written in go

Currently technically functional although it only works with games that don't require window scroling or cartridge bank switching

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
- [ ] Bank switching
- [ ] Sound
- [ ] Window scroll (still work to do)

## TODO (User Interface)
- [ ] file menu
    - [ ] load rom
    - [ ] save state
    - [ ] load state
- [ ] debug windows
    - [ ] tile maps
    - [ ] Cpu registers
    - [ ] Input

