package gdrive

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var parent = "/home/gokalp/Documents/shared_codes/smcp/gdrive"
var tokFile = parent + "/token.json"
var credentialsPath = parent + "/credentials.json"

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	//todo: added for refresh, see it if it works.
	config.TokenSource(context.Background(), tok)
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Unable to read token file %v", err)
		}
	}(f)
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Unable to save token file %v", err)
		}
	}(f)
	json.NewEncoder(f).Encode(token)
}

type GdriveClient struct {
	CredentialsPath string
}

func (g *GdriveClient) createService() (*drive.Service, error) {
	if len(g.CredentialsPath) == 0 {
		g.CredentialsPath = credentialsPath
	}
	b, err := ioutil.ReadFile(g.CredentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}
	client := getClient(config)

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
		return nil, err
	}

	return srv, nil
}

func (g *GdriveClient) query(q string) (*drive.FileList, error) {
	service, err := g.createService()
	if err != nil {
		log.Fatalf("Unable to create drive service: %v", err)
		return nil, err
	}
	r, err := service.Files.List().Q(q).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve directories: %v", err)
		return nil, err
	}

	return r, nil
}

func (g *GdriveClient) createPermission(file *drive.File) (*drive.Permission, error) {
	service, err := g.createService()
	if err != nil {
		log.Fatalf("Unable to create drive service: %v", err)
		return nil, err
	}

	perm := &drive.Permission{
		EmailAddress: "ioniangamer@gmail.com", //todo:read it from a config file.
		Type:         "user",
		Role:         "writer",
	}

	newPerm, err := service.Permissions.Create(file.Id, perm).Do()
	if err != nil {
		return nil, err
	}

	return newPerm, nil
}

func (g *GdriveClient) FindFolderByName(name string) (*drive.File, error) {
	r, err := g.query("mimeType='application/vnd.google-apps.folder' and name='" + name + "'")
	if err != nil {
		log.Fatalf("Unable to retrieve directories: %v", err)
		return nil, err
	}
	if len(r.Files) == 0 {
		return nil, nil
	}

	return r.Files[0], nil
}

func (g *GdriveClient) CreateFolder(name string) (*drive.File, error) {
	service, err := g.createService()
	if err != nil {
		log.Fatalf("Unable to create drive service: %v", err)
		return nil, err
	}

	folder := drive.File{Name: name, MimeType: "application/vnd.google-apps.folder"}
	do, err := service.Files.Create(&folder).Do()
	if err != nil {
		log.Fatalf("Unable to create drive service: %v", err)
		return nil, err
	}

	_, err = g.createPermission(do)
	if err != nil {
		return nil, err
	}

	return do, nil
}

func (g *GdriveClient) GetChildFolders(file *drive.File) (*drive.FileList, error) {
	return g.query("mimeType='application/vnd.google-apps.folder' and '" + file.Id + "' in parents")
}

func (g *GdriveClient) FindChildFolder(parentFolder *drive.File, childName string) (*drive.File, error) {
	r, err := g.query("mimeType='application/vnd.google-apps.folder' and '" + parentFolder.Id + "' in parents and name='" + childName + "'")
	if err != nil {
		log.Fatalf("Unable to retrieve directories: %v", err)
		return nil, err
	}

	if len(r.Files) == 0 {
		return nil, nil
	}

	return r.Files[0], nil
}

func (g *GdriveClient) CreateChildFolder(parentFolder *drive.File, childName string) (*drive.File, error) {
	service, err := g.createService()
	if err != nil {
		log.Fatalf("Unable to create drive service: %v", err)
		return nil, err
	}

	parents := make([]string, 1)
	parents[0] = parentFolder.Id
	childFolder := drive.File{Name: childName, MimeType: "application/vnd.google-apps.folder", Parents: parents}
	do, err := service.Files.Create(&childFolder).Do()
	if err != nil {
		log.Fatalf("Unable to child folder on drive service: %v", err)
		return nil, err
	}

	return do, nil
}

func (g *GdriveClient) CreateImageFile(parentId string, fileName string, imageBase64 *string) (*drive.File, error) {
	service, err := g.createService()
	if err != nil {
		log.Fatalf("Unable to create drive service: %v", err)
		return nil, err
	}

	imageFile := &drive.File{}
	imageFile.Name = fileName
	imageFile.MimeType = "image/jpeg"
	imageFile.Parents = []string{parentId}

	imageBytes, _ := base64.StdEncoding.DecodeString(*imageBase64)
	reader := bytes.NewReader(imageBytes)
	readerCloser := io.NopCloser(reader)
	defer func(readerCloser io.ReadCloser) {
		err := readerCloser.Close()
		if err != nil {
			log.Fatalf("Unable to create iame file: %v", err)
		}
	}(readerCloser)

	do, err := service.Files.Create(imageFile).Media(readerCloser).Do()
	if err != nil {
		log.Fatalf("Unable to create an image file on drive service: %v", err)
		return nil, err
	}

	return do, nil
}
