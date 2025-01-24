package main

import (
    "crypto/aes"
    "crypto/cipher"
    // "crypto/rand"
    "fmt"
    // "io"
	// "io/ioutil"
	"os"
	"encoding/hex"
)

func ShowUsage() {
	fmt.Println("\n")
	fmt.Println("Go secret app - Application encrypts text data. ")

	fmt.Println("Usage:")
	fmt.Println("    go-secret-app-0.1-windows-amd64.exe KEY SECRET_ID")

	fmt.Println("\nAll aguments you provided: ")
	fmt.Println(os.Args)

}

func GetSecretDataHex(secretKey string) string {
	fmt.Println("Getting secret for id: ", secretKey);

	secrets := map[string]string {
		"SECRET": "f2925690aeb411443bb2fe3756118b2acf6c00c96c6bccc26cc5738edd96935414f4fd302ce8a4e8768e226787d7eeffda4daf20",
    	"PASS":   "a49f58c6dce0536e54f4cf5d064932dfe5702d2f280f7b292861206aae75f70b76a51ceb9982e63e5950",
	}

	secretValue := secrets[secretKey];
	if secretValue == "" || secretValue == "null" {
        fmt.Println("ERROR::Secret for key: ", secretKey, " was not found!!")
		return "";
    } else {
		return secretValue;
	}
	// return "f2925690aeb411443bb2fe3756118b2acf6c00c96c6bccc26cc5738edd96935414f4fd302ce8a4e8768e226787d7eeffda4daf20";
} 


func Decode(skey string, hexString string) string {

	if (len(skey) != 32) {
		fmt.Println("    Encryption key must be 32 characters !!!");
		return "";
	}

	if hexString == "" || hexString == "null" {
        fmt.Println("ERROR::Secret was not found!!")
		return "";
	}

	fmt.Println("Hex String: ", hexString)

	dataToDecodeByteArray, err := hex.DecodeString(hexString)
	
	if err != nil {
		fmt.Println("Unable to convert hex to byte. ", err)
	}

	ciphertext := dataToDecodeByteArray;

    key := []byte(skey)
	

	c, err := aes.NewCipher(key)
    if err != nil {
        fmt.Println(err)
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        fmt.Println(err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        fmt.Println(err)
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(plaintext))

	return string(plaintext)

}


// func Encode(skey string, stext string) string {
// 	if (len(skey) != 32) {
// 		fmt.Println("    Encryption key must be 32 characters !!!");
// 		return "";
// 	}

// 	text := []byte(stext)
//     key := []byte(skey)

// 	// generate a new aes cipher using our 32 byte long key
//     c, err := aes.NewCipher(key)
//     // if there are any errors, handle them
//     if err != nil {
//         fmt.Println(err)
//     }

//     // gcm or Galois/Counter Mode, is a mode of operation
//     // for symmetric key cryptographic block ciphers
//     // - https://en.wikipedia.org/wiki/Galois/Counter_Mode
//     gcm, err := cipher.NewGCM(c)
//     // if any error generating new GCM
//     // handle them
//     if err != nil {
//         fmt.Println(err)
//     }

//     // creates a new byte array the size of the nonce
//     // which must be passed to Seal
//     nonce := make([]byte, gcm.NonceSize())
//     // populates our nonce with a cryptographically secure
//     // random sequence
//     if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
//         fmt.Println(err)
//     }

//     // here we encrypt our text using the Seal function
//     // Seal encrypts and authenticates plaintext, authenticates the
//     // additional data and appends the result to dst, returning the updated
//     // slice. The nonce must be NonceSize() bytes long and unique for all
//     // time, for a given key.
// 	encodedByteArray := gcm.Seal(nonce, nonce, text, nil)
//     //fmt.Println(encodedByteArray)
// 	str := hex.EncodeToString(encodedByteArray)
// 	fmt.Println(str)
// 	return str
// }

func main() {
    fmt.Println("Go Secret App v0.1")

	argsWithProgramName := os.Args
	argsWithoutProgramName := os.Args[1:]
	fmt.Println("Application and arguments: ",argsWithProgramName)
	fmt.Println("only arguments: ",argsWithoutProgramName)
	fmt.Println("\n");
	fmt.Println("-------------------------------------------------------");
	if len(os.Args) > 2 {
		secretId := os.Args[2];
		secretDataHEX := GetSecretDataHex(secretId);
		
		encryptionKey := os.Args[1];
		fmt.Println("Decoding: ", secretId, " with key: ", encryptionKey)
		// encoded := Encode(os.Args[2], os.Args[3])
		decoded := Decode(encryptionKey, secretDataHEX)
		if decoded == "" || decoded == "null" {
			fmt.Println("ERROR::Secret was not found!!")
		} else {
			fmt.Println("\n");
			fmt.Println("Secret: ", secretId, decoded);
			fmt.Println("\n");
		}
		
		
	} else {
		ShowUsage();
	}

    // text := []byte("My Super Secret Code Stuff")
    // key := []byte("passphrasewhichneedstobe32bytes!")

    
}
