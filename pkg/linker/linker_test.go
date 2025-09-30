package linker

import (
	"testing"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/pkg/test"
)

var cfg = config.Config{
	ServerAddress: "localhost",
	ServerPort:    "8080",
}

func TestNewLinker(t *testing.T) {

	linker := NewVerifLinker(&cfg)

	test.NotEqual(t, linker, nil)

}

func TestNewLinkCreation(t *testing.T) {
	linker := NewVerifLinker(&cfg)

	links := make([]string, 0, 5)
	for range 5 {
		link, _ := linker.CreateVerificationLink()
		links = append(links, link)
	}

	for i := range links {
		test.NotEqual(t, links[i], "")
	}
}
