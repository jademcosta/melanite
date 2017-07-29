package main

import (
	"fmt"
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
		suite.subjectServer.URL, "http://localhost:8081/park-view.jpg"))
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
		"Content-Type header should be image/jpg")
	suite.Equal(strconv.Itoa(len(body)), res.Header.Get("Content-Length"),
		"Content-Length header should be")
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

// func (suite *FakeExternalServerTestSuite) TestAnswers400WheUrlIsInvalid() {
// 	res, err := http.Get(fmt.Sprintf("%s/%s",
// 		suite.subjectServer.URL, "asdaishdih"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	body, err := ioutil.ReadAll(res.Body)
// 	res.Body.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	suite.Equal(400, res.StatusCode, "status code should be 400")
// 	suite.Equal(string(body), "", "they should be equal")
// }

func (suite *FakeExternalServerTestSuite) TearDownTest() {
	suite.subjectServer.Close()
}

func (suite *FakeExternalServerTestSuite) TearDownSuite() {
	if err := suite.fakeServer.Shutdown(nil); err != nil {
		panic(err)
	}
}
