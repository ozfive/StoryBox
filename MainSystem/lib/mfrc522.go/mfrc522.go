package mfrc522

import (
    "errors"
    "fmt"
    "os"
    "time"

    "github.com/stianeikeland/go-rpio/v4"
)

const (
    idle     = 0x00
    auth     = 0x0E
    receive  = 0x08
    transmit = 0x04
    translen = 16
)

type SimpleMFRC522 struct {
    spi   rpio.SpiDev
    csPin rpio.Pin
    rstPin rpio.Pin
}

func NewSimpleMFRC522() (*SimpleMFRC522, error) {
    spi := rpio.SpiDev{
        Bus:     0,
        ChipSelect: 0,
        Speed:   500000,
        Mode:    0,
    }
    if err := spi.Open(); err != nil {
        return nil, err
    }

    csPin := rpio.Pin(22)
    if err := rpio.Open(); err != nil {
        return nil, err
    }
    csPin.Output()
    csPin.High()

    rstPin := rpio.Pin(25)
    if err := rpio.Open(); err != nil {
        return nil, err
    }
    rstPin.Output()
    rstPin.High()

    return &SimpleMFRC522{
        spi: spi,
        csPin: csPin,
        rstPin: rstPin,
    }, nil
}

func (reader *SimpleMFRC522) Read() (int, error) {
    err := reader.init()
    if err != nil {
        return 0, err
    }

    id, err := reader.readCardId()
    if err != nil {
        return 0, err
    }

    return id, nil
}

func (reader *SimpleMFRC522) writeRegister(address, value byte) {
    reader.csPin.Low()
    reader.spi.Write([]byte{address, value})
    reader.csPin.High()
}

func (reader *SimpleMFRC522) readRegister(address byte) byte {
    reader.csPin.Low()
    reader.spi.Write([]byte{address | 0x80, 0x00})
    data := make([]byte, 2)
    reader.spi.Transfer(nil, data)
    reader.csPin.High()
    return data[1]
}

func (reader *SimpleMFRC522) init() error {
    reader.rstPin.Low()
    time.Sleep(100 * time.Millisecond)
    reader.rstPin.High()
    reader.writeRegister(0x01, 0x0f)
    reader.writeRegister(0x2a, 0x8d)
    reader.writeRegister(0x2b, 0x3e)
    reader.writeRegister(0x2c, 0x00)
    reader.writeRegister(0x2d, 0x00)
    reader.writeRegister(0x15, 0x40)
    reader.writeRegister(0x11, 0x3d)
    return nil
}

func (reader *SimpleMFRC522) readCardId() (int, error) {
    buffer := []byte{0, 0, 0, 0, 0}
    buffer[0] = receive

    uid := make([]byte, 5)
    var uidIndex int

    for {
        length := byte(translen)
        data := append(buffer, []byte{length}...)

        reader.writeRegister(0x01, auth)
        reader.writeRegister(0x0a, 0x01)
