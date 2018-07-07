package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/gin-gonic/gin"
	"github.com/xrdcode/aws-hack/api"
)

var (
	IgnoredWords = []string{
		"search",
		"macbook air",
		"dell",
		"home",
		"kompas.com",
		"news",
		"megapolitan",
		"berita",
		"travel news",
		"github",
		"gihux",
		"detiknews",
	}
)

type DetectedResponse struct {
	MaxConfidence float64
	MinConfidence float64
	Detected      []DetectedText
}

type DetectedText struct {
	Confidence   float64 `json:"Confidence"`
	DetectedText string  `json:"DetectedText"`
	Geometry     struct {
		BoundingBox struct {
			Height float64 `json:"Height"`
			Left   float64 `json:"Left"`
			Top    float64 `json:"Top"`
			Width  float64 `json:"Width"`
		} `json:"BoundingBox"`
		Polygon []struct {
			X float64 `json:"X"`
			Y float64 `json:"Y"`
		} `json:"Polygon"`
	} `json:"Geometry"`
	ID       int         `json:"Id"`
	ParentID interface{} `json:"ParentId"`
	Type     string      `json:"Type"`
}

type GlobalResponse struct {
	Data interface{} `json:"data"`
}

type HoaxDetectionResponse struct {
	Text   string      `json:"text"`
	Found  []TextFound `json:"text_found"`
	Detail HoaxDetail  `json:"detail"`
}

type HoaxDetail struct {
	FinalScore    float64 `json:"final_score"`
	SimiliarTitle string  `json:"similiar_title"`
}

type TextFound struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
	Link  string  `json:"link"`
}

func Uploadimg(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload")
	fmt.Println(header.Filename)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {

	}

	byt := buf.Bytes()

	reko := rekognition.New(session.New(), &aws.Config{Region: aws.String("us-east-2")})
	input := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			Bytes: byt,
		},
	}

	detect, err := reko.DetectText(input)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := getText(detect.TextDetections)

	concated := concat(data)

	test, _ := api.CalculateHoax(concated.Text)
	log.Println(test)

	resp := GlobalResponse{
		Data: concated,
	}

	c.JSON(200, resp)
}

func getText(td []*rekognition.TextDetection) (DetectedResponse, error) {
	tmp := DetectedResponse{}
	tmp.MaxConfidence = -0.1
	tmp.MinConfidence = 100.1
	for _, d := range td {
		tmpD := DetectedText{}
		m, err := json.Marshal(d)
		if err != nil {
			return tmp, errors.New(err.Error())
		}

		err = json.Unmarshal(m, &tmpD)
		if err != nil {
			return tmp, errors.New(err.Error())
		}

		if tmpD.Type == "LINE" {

			matched, _ := regexp.MatchString("([/\\_!]|[0-9]{2}(-|/)[0-9]{2}(-|/)[0-9]{4}|\\.[a-zA-Z]|^[A-Z\\s]*$|[0-9]{1,}\\s$|^[A-Za-z0-9]+$)|([@#$%])|^[A-z]{1,}\\s[A-z]{1,}$", tmpD.DetectedText)
			if !matched && !filter(tmpD.DetectedText) {
				if tmpD.Confidence > tmp.MaxConfidence {
					tmp.MaxConfidence = tmpD.Confidence
				}

				if tmpD.Confidence < tmp.MinConfidence {
					tmp.MinConfidence = tmpD.Confidence
				}
				tmp.Detected = append(tmp.Detected, tmpD)
			} else {
				log.Println("Confidence : ", tmpD.Confidence)
				log.Println("Filtered : ", tmpD.DetectedText)
			}
		}
	}

	log.Println("Max Conf : ", tmp.MaxConfidence)

	return tmp, nil
}

func filter(text string) bool {
	for _, f := range IgnoredWords {
		if strings.Contains(strings.ToLower(text), f) {
			return true
		}
	}
	return false
}

func concat(resp DetectedResponse) HoaxDetectionResponse {
	hdr := HoaxDetectionResponse{}
	txt := ""
	counter := 0

	for _, t := range resp.Detected {
		log.Println("Confidence : ", t.Confidence)
		log.Println("Filtered : ", t.DetectedText)
		if t.Confidence >= resp.MaxConfidence-3.0 {
			txt += " " + t.DetectedText
			log.Println(t.DetectedText)
			counter++
		}

		if counter == 5 {
			break
		}
	}

	log.Println(txt)
	hdr.Text = txt

	detect, _ := api.CalculateHoax(txt)

	for _, h := range detect.Data {
		found := TextFound{
			Text:  h.Text,
			Link:  h.Link,
			Score: h.Score,
		}
		hdr.Found = append(hdr.Found, found)
	}

	if len(detect.Data) == 0 {
		hdr.Detail = HoaxDetail{
			SimiliarTitle: "",
			FinalScore:    0,
		}
		return hdr
	}

	finalScore := math.Floor(hdr.Found[0].Score*100) / 100

	hdr.Detail = HoaxDetail{
		SimiliarTitle: hdr.Found[0].Text,
		FinalScore:    finalScore,
	}

	return hdr
}
