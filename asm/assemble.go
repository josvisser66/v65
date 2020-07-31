// Package asm contains all the code required to assemble a source file.
package asm

// Assemble assembles a source file.
func Assemble(filename string) error {
	_, err := newSource(filename)
	if err != nil {
		return err
	}
	return nil
}
