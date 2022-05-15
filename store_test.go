package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore(t *testing.T) {
	hs := NewHoleStorage()
	initPosts := []*Post{
		{Pid: "10"}, {Pid: "11"}, {Pid: "12"}, {Pid: "13"}, {Pid: "15"}, {Pid: "16"}, {Pid: "17"}, {Pid: "18"},
	}
	deleted := hs.InsertAndCheck(initPosts)
	assert.Empty(t, deleted)

	secondPosts := []*Post{
		{Pid: "11"}, {Pid: "13"}, {Pid: "15"}, {Pid: "16"}, {Pid: "17"},
	}
	deleted = hs.InsertAndCheck(secondPosts)
	assert.ElementsMatch(t, deleted, []string{"12", "18"})

	thirdPosts := []*Post{
		{Pid: "15"}, {Pid: "16"}, {Pid: "19"}, {Pid: "20"},
	}
	deleted = hs.InsertAndCheck(thirdPosts)
	assert.ElementsMatch(t, deleted, []string{"17"})

	assert.ElementsMatch(t, hs.GetAllDeleted(), []string{"12", "17", "18"})
}
