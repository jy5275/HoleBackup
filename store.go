package main

import (
	"math"
	"strconv"
)

type HoleStorage struct {
	//allPosts     []*Post
	//postMap      map[string]int
	commentsCnt  map[string]int
	deletedPosts map[string]bool
}

func NewHoleStorage() *HoleStorage {
	return &HoleStorage{
		//allPosts:     []*Post{},
		//postMap:      make(map[string]int),
		deletedPosts: make(map[string]bool),
		commentsCnt:  make(map[string]int),
	}
}

func GetLastPostID(newPosts []*Post) string {
	ret := math.MaxInt
	for _, post := range newPosts {
		pid, _ := strconv.Atoi(post.Pid)
		if pid < ret {
			ret = pid
		}
	}

	return strconv.Itoa(ret)
}

func (h *HoleStorage) GetAllDeleted() []string {
	var ans []string
	for k, _ := range h.deletedPosts {
		ans = append(ans, k)
	}

	return ans
}

func (h *HoleStorage) GetAndUpdateCommentsForPosts(newPosts []*Post) {
	var needToUpdateComments []*Post
	for _, p := range newPosts {
		getCnt, _ := strconv.Atoi(p.Reply)
		if cnt, ok := h.commentsCnt[p.Pid]; ok {
			if getCnt <= cnt { // no new comments
				continue
			}
		}

		needToUpdateComments = append(needToUpdateComments, p)
		h.commentsCnt[p.Pid] = getCnt
	}

	revIdxTable := make(map[string]*Post)
	for _, post := range needToUpdateComments {
		revIdxTable[post.Pid] = post
	}

	type Msg struct {
		Comments []*Comment
		Pid      string
	}
	ch := make(chan Msg, 10)
	for _, post := range needToUpdateComments {
		go func(pid string) {
			comments, _ := GetComment(pid)
			ch <- Msg{
				Comments: comments,
				Pid:      pid,
			}
		}(post.Pid)
	}

	for i := 0; i < len(needToUpdateComments); i++ {
		msg := <-ch
		revIdxTable[msg.Pid].Comments = msg.Comments
	}
}

func (h *HoleStorage) InsertAndCheck(newPosts []*Post) []string {
	// Updated comments in DB
	h.GetAndUpdateCommentsForPosts(newPosts)
	for _, p := range newPosts {
		if len(p.Comments) > 0 {
			if err := p.UpdateComments(); err != nil {
				logger.Fatalln("failed to update comments in DB: ", err)
				return []string{}
			}
		}
	}

	startPID := GetLastPostID(newPosts)
	//startIdx := 0

	//if len(h.postMap) > 0 {
	//	startIdx = h.postMap[startPID]
	//}
	var newDeleted []string

	// new:  | 10 | 11 | 13 | 16 | 17 || 20 | 21 | 22 | 24
	// old:  | 10 | 11 | 12 | 13 | 14 | 16 | 17 | 18 |
	pOld, pNew := 0, 0
	allPosts, err := SelectLatestPosts(startPID)
	if err != nil {
		logger.Fatalln("failed to select latest posts: ", err)
		return []string{}
	}

	for pOld < len(allPosts) {
		if pNew >= len(newPosts) || allPosts[pOld].Pid != newPosts[pNew].Pid {
			// Some posts are deleted
			oldPost := allPosts[pOld]
			if _, ok := h.deletedPosts[oldPost.Pid]; !ok {
				// need to process
				h.deletedPosts[oldPost.Pid] = true
				newDeleted = append(newDeleted, oldPost.Pid)
				if err := oldPost.MardAsDeleted(); err != nil {
					logger.Fatalln("failed to mark post as deleted: ", err)
				}
				logger.Printf("Deleted post: %+v, newPosts=%v, allPosts=%v, pNew=%v, pOld=%v\n",
					*oldPost, newPosts, allPosts, pNew, pOld)
			}
			pOld++
			continue
		}
		pNew++
		pOld++
	}

	// Append the new posts
	for pNew < len(newPosts) {
		logger.Println("append new posts: ", newPosts[pNew].Pid)
		if err := newPosts[pNew].FirstOrCreate(); err != nil {
			logger.Fatalln("failed to create post in DB: ", err)
		}

		//h.postMap[newPosts[pNew].Pid] = len(h.allPosts)
		//h.allPosts = append(h.allPosts, newPosts[pNew])
		pNew++
	}

	return newDeleted
}
