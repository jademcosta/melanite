package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jademcosta/melanite/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

const testImagesFolder = "./test/images"
const fakeTestServerPort = ":8081"

type FakeExternalServerTestSuite struct {
	suite.Suite
	fakeServer    *http.Server
	subjectServer *httptest.Server
}

func (suite *FakeExternalServerTestSuite) SetupSuite() {
	suite.fakeServer = &http.Server{Addr: fakeTestServerPort,
		Handler: http.FileServer(http.Dir(testImagesFolder))}
	go func() {
		suite.fakeServer.ListenAndServe()
	}()
}

func TestFakeExternalServerTestSuite(t *testing.T) {
	suite.Run(t, new(FakeExternalServerTestSuite))
}

func (suite *FakeExternalServerTestSuite) SetupTest() {
	logger := log.New()
	logger.SetLevel(log.FatalLevel)
	suite.subjectServer = httptest.NewServer(
		GetApp(config.Config{ImageSource: "http://localhost:8081"}, logger))
}

func (suite *FakeExternalServerTestSuite) TestAnswers404WhenImageNotFound() {
	res, err := http.Get(fmt.Sprintf("%s/%s", suite.subjectServer.URL, "oi.png"))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	suite.Equal(404, res.StatusCode, "status code should be 404")
	suite.Equal(string(body), "", "The body should be empty")
}

func (suite *FakeExternalServerTestSuite) TestAnswers200WhenImageExists() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "park-view-XS.jpg"))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XS.jpg"))
	if err != nil {
		panic(err)
	}

	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/jpeg", res.Header.Get("Content-Type"),
		"Content-Type header should be image/jpg")
	suite.Equal(strconv.Itoa(len(body)), res.Header.Get("Content-Length"),
		"Content-Length header should be sent")
	suite.Equal(len(body), len(dat),
		"The image should be the one we asked for")
}

//
func (suite *FakeExternalServerTestSuite) TestAnswers400WheNoUrlIsGiven() {
	res, err := http.Get(fmt.Sprintf("%s",
		suite.subjectServer.URL))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	suite.Equal(400, res.StatusCode, "status code should be 400")
	suite.Equal(string(body), "", "The body should be empty")
}

func (suite *FakeExternalServerTestSuite) TestAnswers400WhenAskForUnsupportedConversionFormat() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "park-view-XS.jpg?out=pngdas"))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	suite.Equal(400, res.StatusCode, "status code should be 400")
	suite.Equal(string(body), "", "The body should be empty")
}

func (suite *FakeExternalServerTestSuite) TestConvertJpgToPng() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "park-view-XS.jpg?out=png"))
	if err != nil {
		panic(err)
	}

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/png", res.Header.Get("Content-Type"),
		"Content-Type header should be image/png")
}

func (suite *FakeExternalServerTestSuite) TestConvertPngToJpg() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "park-view-XS.png?out=jpg"))
	if err != nil {
		panic(err)
	}

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/jpeg", res.Header.Get("Content-Type"),
		"Content-Type header should be image/jpeg")
}

func (suite *FakeExternalServerTestSuite) TestConvertPngToWebp() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "park-view-XS.png?out=webp"))
	if err != nil {
		panic(err)
	}

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/webp", res.Header.Get("Content-Type"),
		"Content-Type header should be image/webp")
}

func (suite *FakeExternalServerTestSuite) TestConvertPngToPng() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "park-view-XS.png?out=png"))
	if err != nil {
		panic(err)
	}

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/png", res.Header.Get("Content-Type"),
		"Content-Type header should be image/png")
}

func (suite *FakeExternalServerTestSuite) TestResize() {
	//260x147
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "park-view-XS.png?res=130x0"))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	suite.Equal(130, img.Bounds().Dx(), "The width should be resized")
	suite.Equal(74, img.Bounds().Dy(), "The height should be resized")
	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/png", res.Header.Get("Content-Type"),
		"Content-Type header should be image/png")
}

func (suite *FakeExternalServerTestSuite) TestEnlargement() {
	//260x147
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "park-view-XS.png?res=520x0"))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	suite.Equal(520, img.Bounds().Dx(), "The width should be resized")
	suite.Equal(294, img.Bounds().Dy(), "The height should be resized")
	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/png", res.Header.Get("Content-Type"),
		"Content-Type header should be image/png")
}

func (suite *FakeExternalServerTestSuite) TearDownTest() {
	suite.subjectServer.Close()
}

func (suite *FakeExternalServerTestSuite) TearDownSuite() {
	if err := suite.fakeServer.Shutdown(nil); err != nil {
		panic(err)
	}
}
