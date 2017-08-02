package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

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
	suite.subjectServer = httptest.NewServer(GetApp())
}

func (suite *FakeExternalServerTestSuite) TestAnswers404WhenImageNotFound() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "http://localhost:8081/oi.png"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	suite.Equal(404, res.StatusCode, "status code should be 404")
	suite.Equal(string(body), "\n", "they should be equal")
}

func (suite *FakeExternalServerTestSuite) TestAnswers200WhenImageExists() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "http://localhost:8081/park-view-XS.jpg"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XS.jpg"))
	if err != nil {
		panic(err)
	}

	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/jpeg", res.Header.Get("Content-Type"),
		"Content-Type header should be image/jpg")
	suite.Equal(strconv.Itoa(len(body)), res.Header.Get("Content-Length"),
		"Content-Length header should be")
	suite.Equal(len(body), len(dat),
		"The image should be the on we asked for")
}

func (suite *FakeExternalServerTestSuite) TestAnswers400WheNoUrlIsGiven() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, ""))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	suite.Equal(400, res.StatusCode, "status code should be 400")
	suite.Equal(string(body), "Invalid image url provided\n", "they should be equal")
}

func (suite *FakeExternalServerTestSuite) TestAnswers400WhenUrlIsInvalid() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "asdaishdih"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	suite.Equal(400, res.StatusCode, "status code should be 400")
	suite.Equal(string(body), "Invalid image url provided\n", "they should be equal")
}

func (suite *FakeExternalServerTestSuite) TestAnswers400WhenUrlDoNotStartWithHttp() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "google.com"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	suite.Equal(400, res.StatusCode, "status code should be 400")
	suite.Equal(string(body), "Invalid image url provided\n", "they should be equal")
}

func (suite *FakeExternalServerTestSuite) TestAnswers400WhenAskForUnsupportedConversionFormat() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "http://localhost:8081/park-view-XS.jpg?o=pngdas"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	suite.Equal(400, res.StatusCode, "status code should be 400")
	suite.Equal(string(body), "Image conversion format not supported\n",
		"they should be equal")
}

func (suite *FakeExternalServerTestSuite) TestConvertJpgToPng() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "http://localhost:8081/park-view-XS.jpg?o=png"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/png", res.Header.Get("Content-Type"),
		"Content-Type header should be image/png")

	_, format, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	suite.Equal("png", format, "The returned image should be a PNG")
}

func (suite *FakeExternalServerTestSuite) TestConvertPngToJpg() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "http://localhost:8081/park-view-XS.png?o=jpg"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/jpeg", res.Header.Get("Content-Type"),
		"Content-Type header should be image/jpeg")

	_, format, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	suite.Equal("jpeg", format, "The returned image should be a JPG")
}

func (suite *FakeExternalServerTestSuite) TestConvertPngToPng() {
	res, err := http.Get(fmt.Sprintf("%s/%s",
		suite.subjectServer.URL, "http://localhost:8081/park-view-XS.png?o=png"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	suite.Equal(200, res.StatusCode, "status code should be 200")
	suite.Equal("image/png", res.Header.Get("Content-Type"),
		"Content-Type header should be image/png")

	_, format, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	suite.Equal("png", format, "The returned image should be a PNG")
}

func (suite *FakeExternalServerTestSuite) TearDownTest() {
	suite.subjectServer.Close()
}

func (suite *FakeExternalServerTestSuite) TearDownSuite() {
	if err := suite.fakeServer.Shutdown(nil); err != nil {
		panic(err)
	}
}
