package wallet

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

const (
	addressFileMode fs.FileMode = 0664
)

// Simple class to wrap the node's address file
type addressManager struct {
	path     string
	address  common.Address
	isLoaded bool
}

// Creates a new address manager
func newAddressManager(path string) *addressManager {
	return &addressManager{
		path: path,
	}
}

// Gets the address saved on disk. Returns false if the address file doesn't exist.
func (m *addressManager) LoadAddress() (common.Address, bool, error) {
	m.address = common.Address{}
	m.isLoaded = false

	_, err := os.Stat(m.path)
	if errors.Is(err, fs.ErrNotExist) {
		return common.Address{}, false, nil
	} else if err != nil {
		return common.Address{}, false, fmt.Errorf("error checking if address file exists: %w", err)
	}

	bytes, err := os.ReadFile(m.path)
	if err != nil {
		return common.Address{}, false, fmt.Errorf("error loading address file [%s]: %w", m.path, err)
	}
	m.address = common.HexToAddress(string(bytes))
	m.isLoaded = true
	return m.address, true, nil
}

// Get the cached address
func (m *addressManager) GetAddress() (common.Address, bool) {
	return m.address, m.isLoaded
}

// Sets the node address without saving it to disk
func (m *addressManager) SetAddress(newAddress common.Address) {
	m.address = newAddress
	m.isLoaded = true
}

// Sets the node address and saves it to disk
func (m *addressManager) SetAndSaveAddress(newAddress common.Address) error {
	m.address = newAddress
	m.isLoaded = true
	bytes := []byte(newAddress.Hex())
	err := os.WriteFile(m.path, bytes, addressFileMode)
	if err != nil {
		return fmt.Errorf("error writing address file [%s] to disk: %w", m.path, err)
	}
	return nil
}
