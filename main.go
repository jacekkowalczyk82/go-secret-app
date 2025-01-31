package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
	"io"
	"bufio"
	"os"
	"strings"
	"encoding/hex"
	"encoding/base64"
	"github.com/atotto/clipboard"
)

func ShowUsage() {
	fmt.Println("\n")
	fmt.Println("Go secret app - Application encrypts text data. ")

	fmt.Println("Usage:")
	fmt.Println("    go-secret-app-0.2-windows-amd64.exe --add PATH_TO_SECRETS_FILE ENCODING_KEY SECRET_ID SECRET_VALUE")
	fmt.Println("    go-secret-app-0.2-windows-amd64.exe --get PATH_TO_SECRETS_FILE ENCODING_KEY SECRET_ID")

	fmt.Println("\nAll aguments you provided: ")
	fmt.Println(os.Args)

}

var debugEnabled bool = false;

func GetSecretDataHex(secretId string, secretLines []string) string {
	fmt.Println("Getting secret for id: ", secretId);

	secretValue := ""

	for i, line := range secretLines {
		if debugEnabled {
			fmt.Printf("DEBUG::Line %d: %s\n", i+1, line)
		}
		if strings.HasPrefix(line, secretId) {
			// Split the string by the separator
			parts := strings.Split(line, ":")
			secretValue = parts[1]

			return secretValue
		}
	}
	fmt.Println("Warning::No Secrets found ")
	return ""

} 

// saveListToFile function saves a list of strings to a specified file
func SaveListToFile(filePath string, list []string) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for _, line := range list {
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            return err
        }
    }
    
    err = writer.Flush()
    if err != nil {
        return err
    }

    return nil
}

func ReadSecretsDataFileLines(filePath string) ([]string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return lines, nil
}

func Decode(encryptionKey string, hexString string) string {

	if (len(encryptionKey) != 32) {
		fmt.Println("\nERROR::Encryption key must be 32 characters !!!");
		return "";
	}

	if hexString == "" || hexString == "null" {
        fmt.Println("\nERROR::Secret was not found!!")
		return "";
	}
	if debugEnabled {
		fmt.Println("DEBUG::Hex String: ", hexString)
	}
	dataToDecodeByteArray, err := hex.DecodeString(hexString)
	
	if err != nil {
		fmt.Println("\nError::Unable to convert hex to byte. ", err)
	}

	ciphertext := dataToDecodeByteArray;

    key := []byte(encryptionKey)
	

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
    if debugEnabled {
		fmt.Println("DEBUG::" + string(plaintext))
	}
	return string(plaintext)

}


func Encode(encryptionKey string, stext string) string {
	if (len(encryptionKey) != 32) {
		fmt.Println("\nWarning::Encryption key must be 32 characters !!!");
		return "";
	}

	text := []byte(stext)
    key := []byte(encryptionKey)

	// generate a new aes cipher using our 32 byte long key
    c, err := aes.NewCipher(key)
    // if there are any errors, handle them
    if err != nil {
        fmt.Println(err)
    }

    // gcm or Galois/Counter Mode, is a mode of operation
    // for symmetric key cryptographic block ciphers
    // - https://en.wikipedia.org/wiki/Galois/Counter_Mode
    gcm, err := cipher.NewGCM(c)
    // if any error generating new GCM
    // handle them
    if err != nil {
        fmt.Println(err)
    }

    // creates a new byte array the size of the nonce
    // which must be passed to Seal
    nonce := make([]byte, gcm.NonceSize())
    // populates our nonce with a cryptographically secure
    // random sequence
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        fmt.Println(err)
    }

    // here we encrypt our text using the Seal function
    // Seal encrypts and authenticates plaintext, authenticates the
    // additional data and appends the result to dst, returning the updated
    // slice. The nonce must be NonceSize() bytes long and unique for all
    // time, for a given key.
	encodedByteArray := gcm.Seal(nonce, nonce, text, nil)
    //fmt.Println(encodedByteArray)
	str := hex.EncodeToString(encodedByteArray)
	
	if debugEnabled {
		fmt.Println("DEBUG::" + str)
	}
	return str
}

// encodeToBase64 funkcja koduje podany tekst do formatu base64
func encodeToBase64(data string) string {
    encoded := base64.StdEncoding.EncodeToString([]byte(data))
    return encoded
}

// decodeFromBase64 funkcja dekoduje podany tekst z formatu base64 do zwykłego tekstu
func decodeFromBase64(encodedData string) (string, error) {
    decodedBytes, err := base64.StdEncoding.DecodeString(encodedData)
    if err != nil {
        return "", err
    }
    return string(decodedBytes), nil
}

func fileExists(filePath string) bool {
    _, err := os.Stat(filePath)
    if os.IsNotExist(err) {
        return false
    }
    return err == nil
}

func main() {
    fmt.Println("Go Secret App v0.1")

	argsWithProgramName := os.Args
	argsWithoutProgramName := os.Args[1:]
	fmt.Println("Application and arguments: ",argsWithProgramName)
	fmt.Println("only arguments: ",argsWithoutProgramName)
	fmt.Println("\n");
	fmt.Println("-------------------------------------------------------");
	if len(os.Args) == 5 || len(os.Args) == 6 { 
		// 4 or 5 arguments
		
		operation := os.Args[1];
		secretsFilePath := os.Args[2];
		encryptionKey := os.Args[3];
		secretId := os.Args[4];

		// Declare a slice of strings without initializing it
    	var lines []string
		var secretLines []string
		var err error
		var secretsMap = make(map[string]string)


		if fileExists(secretsFilePath) {
			fmt.Printf("\nReading secrets from file %s.\n", secretsFilePath)
			lines, err = ReadSecretsDataFileLines(secretsFilePath)
			secretLines = []string{} //clean secret lines in form: secret_id:hex_secreet_encoded_value 
			if err != nil {
				fmt.Println("\nError reading file:", err)
				return
			} 

			for i, line := range lines {
				if debugEnabled {
					fmt.Printf("DEBUG::Line %d: %s\n", i+1, line)
				}
				existingSecretLine, err := decodeFromBase64(line)
				if err != nil {
					fmt.Println("\nError decoding base64 line:", err)
					return
				}
				secretLines = append(secretLines, existingSecretLine)
				secretLineParts := strings.Split(existingSecretLine, ":")
				existingSecretId := secretLineParts[0]
				// existingSecretValue := secretLineParts[1]

				if debugEnabled {
					fmt.Printf("DEBUG::Adding secret %s: %s\n", existingSecretId, existingSecretLine)
				}
				secretsMap[existingSecretId] = existingSecretLine
			}

		} else {
			lines = []string{}
			secretLines = []string{}
		}

		

		if operation == "--get"  {

			secretDataHEX := GetSecretDataHex(secretId,secretLines);
		
			fmt.Println("\nDecoding: ", secretId)
			
			if debugEnabled {
				fmt.Println("\nDecoding: ", secretId, " with key: ", encryptionKey)
			}
			
			// encoded := Encode(os.Args[2], os.Args[3])
			decoded := Decode(encryptionKey, secretDataHEX)
			if decoded == "" || decoded == "null" {
				fmt.Println("ERROR::Secret was not found!!")
			} else {
				fmt.Println("\n");
				if debugEnabled {
					fmt.Println("DEBUG::Secret: ", secretId, decoded);
				}

					// Wstaw tekst do schowka
				err := clipboard.WriteAll(decoded)
				if err != nil {
					fmt.Println("\nError while inserting secret to the clipboard: ", err)
					return
				}

				fmt.Println("The secret was inserted to the clipboard")
				fmt.Println("\n");
			}
			
		} else if operation == "--add" && len(os.Args) == 6 {
			newBase64Lines := []string{}

			newSecretValue := os.Args[5]
			
			encodedSecretValue := Encode(encryptionKey, newSecretValue)
			secretLine := secretId + ":" + encodedSecretValue
			secretsMap[secretId] = secretLine

			newSecretLines := []string{}

			for _, value := range secretsMap {
				newSecretLines = append(newSecretLines, value)
			}

			for i, line := range newSecretLines {
				if debugEnabled {
					fmt.Printf("DEBUG::SecretLine %d: %s\n", i+1, line)
				}
				base64Line := encodeToBase64(line)
				newBase64Lines = append(newBase64Lines, base64Line)
			}

			SaveListToFile(secretsFilePath, newBase64Lines)
			fmt.Println("\nThe secret was added to the secrets file: " + secretsFilePath)

		}
		
		
	} else {
		ShowUsage();
	}
    
}
