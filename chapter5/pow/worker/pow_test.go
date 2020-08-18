package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func testHash(t *testing.T, prefix string, difficulty int) {
	r, _ := http.NewRequest("GET", "/start?prefix="+
		url.QueryEscape(prefix)+"&difficulty="+
		strconv.Itoa(difficulty), nil)
	w := httptest.NewRecorder()

	Router().ServeHTTP(w, r)
	bs, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !Hash(bs, 10) {
		t.Log(string(bs))
		t.Fail()
	}
}

func TestDifficulty10(t *testing.T) {
	for i := 0; i <= 100000; i++ {
		testHash(t, "test", 10)
	}
}

func TestDifficulty15(t *testing.T) {
	for i := 0; i <= 10000; i++ {
		testHash(t, "test", 15)
	}
}

func TestDifficulty20(t *testing.T) {
	for i := 0; i <= 2; i++ {
		testHash(t, "test", 20)
	}
}
