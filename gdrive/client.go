package gdrive

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"io"
	"log"
	"net/http"
	"smcp/reps"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Client struct {
	Repository *reps.CloudRepository
	srv        *drive.Service
	pool       redsyncredis.Pool
}

// Retrieve a token, saves the token, then returns the generated client.
func (g *Client) getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := g.Repository.GetGdriveToken()
	if err != nil {
		tok = g.getTokenFromWeb(config)
		if tok == nil {
			log.Println("gdrive token is null, no action taken")
			return nil
		}
		err = g.Repository.SaveGdriveToken(tok)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
	}
	//todo: added for refresh, see it if it works.
	config.TokenSource(context.Background(), tok)
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func (g *Client) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authCode, _ := g.Repository.GetGdriveAuthCode()

	if len(authCode) == 0 {
		//then generate URL to get AuthCode
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		err := g.Repository.SaveGdriveUrl(authURL)
		if err != nil {
			log.Printf(err.Error())
			return nil
		}
		log.Printf("Auth Code has been saved : \n%v\n", authURL)
		return nil
	} else {
		tok, err := config.Exchange(context.TODO(), authCode)
		if err != nil {
			log.Printf("Unable to retrieve token from web %v", err)
			return nil
		}
		return tok
	}
}

func (g *Client) createService() (*drive.Service, error) {
	if g.srv != nil {
		return g.srv, nil
	}

	if g.pool == nil {
		g.pool = goredis.NewPool(g.Repository.Connection)
	}
	rs := redsync.New(g.pool)
	mutex := rs.NewMutex("mutex-gdrive")
	if e := mutex.Lock(); e != nil {
		log.Println("An error occurred on GdriveClient mutex lock: " + e.Error())
		return nil, e
	}
	defer func(mutex *redsync.Mutex) {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			log.Println("An error occurred on GdriveClient mutex unlock: " + err.Error())
		}
	}(mutex)

	if g.srv != nil {
		return g.srv, nil
	}

	b, err := g.Repository.GetGdriveCredentials()
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(b), drive.DriveScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}
	client := g.getClient(config)

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Drive client: %v", err)
		return nil, err
	}
	g.srv = srv

	return srv, nil
}

func (g *Client) query(q string) (*drive.FileList, error) {
	service, err := g.createService()
	if err != nil {
		log.Printf("Unable to create drive service: %v", err)
		return nil, err
	}
	r, err := service.Files.List().Q(q).Do()
	if err != nil {
		log.Printf("Unable to retrieve directories: %v", err)
		return nil, err
	}

	return r, nil
}

func (g *Client) createPermission(file *drive.File) (*drive.Permission, error) {
	service, err := g.createService()
	if err != nil {
		log.Printf("Unable to create drive service: %v", err)
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

func (g *Client) FindFolderByName(name string) (*drive.File, error) {
	r, err := g.query("mimeType='application/vnd.google-apps.folder' and name='" + name + "'")
	if err != nil {
		log.Printf("Unable to retrieve directories: %v", err)
		return nil, err
	}
	if len(r.Files) == 0 {
		return nil, nil
	}

	return r.Files[0], nil
}

func (g *Client) CreateFolder(name string) (*drive.File, error) {
	service, err := g.createService()
	if err != nil {
		log.Printf("Unable to create drive service: %v", err)
		return nil, err
	}

	folder := drive.File{Name: name, MimeType: "application/vnd.google-apps.folder"}
	do, err := service.Files.Create(&folder).Do()
	if err != nil {
		log.Printf("Unable to create drive service: %v", err)
		return nil, err
	}

	_, err = g.createPermission(do)
	if err != nil {
		return nil, err
	}

	return do, nil
}

func (g *Client) GetChildFolders(file *drive.File) (*drive.FileList, error) {
	return g.query("mimeType='application/vnd.google-apps.folder' and '" + file.Id + "' in parents")
}

func (g *Client) FindChildFolder(parentFolder *drive.File, childName string) (*drive.File, error) {
	r, err := g.query("mimeType='application/vnd.google-apps.folder' and '" + parentFolder.Id + "' in parents and name='" + childName + "'")
	if err != nil {
		log.Printf("Unable to retrieve directories: %v", err)
		return nil, err
	}

	if len(r.Files) == 0 {
		return nil, nil
	}

	return r.Files[0], nil
}

func (g *Client) CreateChildFolder(parentFolder *drive.File, childName string) (*drive.File, error) {
	service, err := g.createService()
	if err != nil {
		log.Printf("Unable to create drive service: %v", err)
		return nil, err
	}

	parents := make([]string, 1)
	parents[0] = parentFolder.Id
	childFolder := drive.File{Name: childName, MimeType: "application/vnd.google-apps.folder", Parents: parents}
	do, err := service.Files.Create(&childFolder).Do()
	if err != nil {
		log.Printf("Unable to child folder on drive service: %v", err)
		return nil, err
	}

	return do, nil
}

func (g *Client) CreateImageFile(parentId string, fileName string, imageBase64 *string) (*drive.File, error) {
	service, err := g.createService()
	if err != nil {
		log.Printf("Unable to create drive service: %v", err)
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
			log.Printf("Unable to create iame file: %v", err)
		}
	}(readerCloser)

	do, err := service.Files.Create(imageFile).Media(readerCloser).Do()
	if err != nil {
		log.Printf("Unable to create an image file on drive service: %v", err)
		return nil, err
	}

	return do, nil
}
