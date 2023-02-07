package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io/ioutil"
)
func validateSignature(signed []byte, signer *openpgp.Entity) error {
	block, err := armor.Decode(bytes.NewReader(signed))
	if err != nil {
		return err
	}

	message, err := openpgp.ReadMessage(block.Body, signer.PrimaryKey, nil, nil)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(message.UnverifiedBody)
	if err != nil {
		return err
	}

	if !message.IsSigned {
		return fmt.Errorf("message is not signed")
	}

	if _, err := message.Verify(); err != nil {
		return err
	}

	return nil
}


func signString(plaintext []byte, signer *openpgp.Entity) ([]byte, error) {
	var buf bytes.Buffer
	w, err := armor.Encode(&buf, openpgp.SignatureType, nil)
	if err != nil {
		return nil, err
	}
	if err := openpgp.ArmoredDetachSign(w, signer, bytes.NewReader(plaintext), nil); err != nil {
		return nil, err
	}
	w.Close()
	return buf.Bytes(), nil
}

func main() {
	// Load the signer's private key
	signerKey, err := ioutil.ReadFile("signer.key")
	if err != nil {
		fmt.Println(err)
		return
	}
	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(signerKey))
	if err != nil {
		fmt.Println(err)
		return
	}
	signer := entityList[0]

	// Sign the string "Hello, World!"
	plaintext := []byte("Hello, World!")
	signed, err := signString(plaintext, signer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(signed))
}

