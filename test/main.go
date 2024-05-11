package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/koffeinsource/go-imgur"
	// "external/koffeinsource/go-imgur"
)

func main() {

	// f, _ := os.Open("golang.png")
	// fmt.Println(f.Name())
	// upload(f, "")
	// fmt.Println()

	v, _ := os.ReadFile("golang.png")
	// upload2(*bytes.NewReader([]byte("s")))

	fmt.Println(len(v))

	clientID := "01e493a6685580c"
	client, err := imgur.NewClient(new(http.Client), clientID, "")
	if err != nil {
		fmt.Printf("failed during imgur client creation. %+v\n", err)
		return
	}

	info, st, err := client.UploadImage(v, "", "URL", "test title", "test description") // info
	// info, st, err := client.UploadImage("/home/keril/Curs 3.2/url-shortener-main/test/golang.png", "", "test title", "test desc")
	if st != 200 || err != nil {
		fmt.Printf("Status: %v\n", st)
		fmt.Printf("Err: %v\n", err)
	}

	fmt.Println(info.Link)
	fmt.Println("///////////////////")

	info2, st, err := client.GetImageInfo(info.ID)
	if st != 200 || err != nil {
		fmt.Printf("Status: %v\n", st)
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Println(info.ID)
	fmt.Println(info2.ID)
	fmt.Println(info2.Link)
}

// type Client struct {
// 	imgur_cli *imgur.Client
// 	http_cli       *http.Client
// }

// func (client *Client) Upload(image []byte, album string, dtype string, title string, description string) (*imgur.ImageInfo, int, error) {
// 	if image == nil {
// 		return nil, -1, errors.New("Invalid image")
// 	}
// 	if dtype != "file" && dtype != "base64" && dtype != "URL" {
// 		return nil, -1, errors.New("Passed invalid dtype: " + dtype + ". Please use file/base64/URL.")
// 	}

// 	URL := "https://api.imgur.com/3/image"
// 	req, err := createRequest(image, album, dtype, title, description) //http.NewRequest("POST", URL, bytes.NewBufferString(form.Encode()))
// 	client.Log.Debugf("Posting to URL %v\n", URL)
// 	if err != nil {
// 		return nil, -1, err //errors.New("Could create request for " + URL + " - " + err.Error())
// 	}

// 	req.Header.Add("Authorization", "Client-ID "+os.Getenv("IMGURID"))
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

// 	res, err := client.http_cli.Do(req)
// 	if err != nil {
// 		return nil, -1, errors.New("Could not post " + URL + " - " + err.Error())
// 	}
// 	defer res.Body.Close()

// 	// Read the whole body
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, -1, errors.New("Problem reading the body of " + URL + " - " + err.Error())
// 	}

// 	// client.Log.Debugf("%v\n", string(body[:]))

// 	dec := json.NewDecoder(bytes.NewReader(body))
// 	var img imageInfoDataWrapper
// 	if err = dec.Decode(&img); err != nil {
// 		return nil, -1, errors.New("Problem decoding json result from image upload - " + err.Error() + ". JSON(?): " + string(body))
// 	}

// 	if !img.Success {
// 		return nil, img.Status, errors.New("Upload to imgur failed with status: " + strconv.Itoa(img.Status))
// 	}

// 	img.Ii.Limit, _ = extractRateLimits(res.Header)

// 	return img.Ii, img.Status, nil
// }

// func createRequest(image []byte, album string, dtype string, title string, description string) (*http.Request, error) {
// 	form := url.Values{}

// 	form.Add("image", string(image[:]))
// 	form.Add("type", dtype)

// 	var buf = new(bytes.Buffer)
// 	writer := multipart.NewWriter(buf)

// 	part, _ := writer.CreateFormFile("image", "dont care about name")
// 	io.Copy(part, image)

// 	if album != "" {
// 		writer.WriteField("album", album)
// 	}
// 	if title != "" {
// 		writer.WriteField("title", title)
// 	}
// 	if description != "" {
// 		writer.WriteField("description", description)
// 	}
// 	writer.Close()
// 	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", buf)
// 	if err != nil {
// 		return nil, errors.New("Could not post " + "https://api.imgur.com/3/image" + " - " + err.Error())
// 	}
// 	req.Header.Set("Content-Type", writer.FormDataContentType())
// 	// req.Header.Set("Authorization", "Bearer "+token)
// 	// req.Header.Add("Authorization", "Client-ID "+"01e493a6685580c")

// 	return req, nil
// }

// func upload2(image []byte) {
// 	form := url.Values{}

// 	form.Add("image", string(image[:]))
// 	form.Add("type", "base64")

// 	form.Add("title", "simple title")
// 	form.Add("description", "simple descript")

// 	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", bytes.NewBufferString(form.Encode())) //bytes.NewBufferString(form.Encode()))

// 	req.Header.Add("Authorization", "Client-ID "+"01e493a6685580c")
// 	req.Header.Add("Content-Type", "multipart/form-data") //"application/x-www-form-urlencoded")

// 	client := &http.Client{}
// 	res, _ := client.Do(req)
// 	defer res.Body.Close()
// 	body, _ := ioutil.ReadAll(res.Body)

// 	fmt.Println(string(body))

// }

// func upload(image io.Reader, token string) {
// 	var buf = new(bytes.Buffer)
// 	writer := multipart.NewWriter(buf)

// 	part, _ := writer.CreateFormFile("image", "dont care about name")
// 	io.Copy(part, image)

// 	writer.Close()
// 	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", buf)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())
// 	// req.Header.Set("Authorization", "Bearer "+token)
// 	req.Header.Add("Authorization", "Client-ID "+"01e493a6685580c")

// 	fmt.Println()
// 	fmt.Println()
// 	fmt.Println(req)

// 	// client := &http.Client{}
// 	// res, _ := client.Do(req)
// 	// defer res.Body.Close()
// 	// body, _ := ioutil.ReadAll(res.Body)

// 	// dec := json.NewDecoder(bytes.NewReader(body))
// 	// var img imageInfoDataWrapper
// 	// if err := dec.Decode(&img); err != nil {
// 	// 	errors.New("Problem decoding json result from image upload - " + err.Error() + ". JSON(?): " + string(body))
// 	// }

// 	// if !img.Success {
// 	// 	errors.New("Upload to imgur failed with status: " + strconv.Itoa(img.Status))
// 	// }

// 	// fmt.Println(string(body))
// }
