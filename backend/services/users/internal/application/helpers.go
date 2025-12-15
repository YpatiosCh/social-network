package application

import "fmt"

// HashPassword hashes a password using bcrypt.
// func hashPassword(password string) (string, error) {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(hash), err
// }

func checkPassword(storedPassword, newHashedPassword string) bool {
	// err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	// return err == nil
	fmt.Println("Comparing passwords:", storedPassword, newHashedPassword)
	return storedPassword == newHashedPassword
}
