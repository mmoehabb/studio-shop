package ancillaries

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	fmt.Printf("\nAuth Code: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	authCode := scanner.Text()
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return token
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens,
	// and is created automatically when the authorization flow completes
	// for the first time.
	tokFile := "token.json"
	token, err := tokenFromFile(tokFile)
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken(tokFile, token)
	}
	return config.Client(context.Background(), token)
}

// Get an instance of drive.Service struct.
// Usage Example: service.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
func GetDriveService() *drive.Service {
	ctx := context.Background()
	cred := Must(os.ReadFile("credentials.json")).([]byte)

	// If modifying these scopes, delete your previously saved token.json.
	config := Must(google.ConfigFromJSON(cred, drive.DriveMetadataReadonlyScope)).(*oauth2.Config)
	client := getClient(config)

	service := Must(drive.NewService(ctx, option.WithHTTPClient(client))).(*drive.Service)
	return service
}
