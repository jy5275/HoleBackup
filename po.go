package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

var (
	db *gorm.DB
)

func init() {
	dsn := "root:root@tcp(188.166.217.29:3306)/hole?charset=utf8mb4"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalln(err)
	}
}

type Post struct {
	ID        int        `gorm:"primaryKey;column:id" json:"-"`
	Pid       string     `gorm:"column:pid" json:"pid"`
	Extra     string     `gorm:"column:extra" json:"extra"`
	Hidden    string     `gorm:"column:hidden" json:"hidden"`
	Hot       string     `gorm:"column:hot" json:"hot"`
	Likenum   string     `gorm:"column:likenum" json:"likenum"`
	Reply     string     `gorm:"column:reply" json:"reply"`
	Tag       *string    `gorm:"column:tag" json:"tag"`
	Text      string     `gorm:"column:text" json:"text"`
	Timestamp string     `gorm:"column:timestamp" json:"timestamp"`
	Type      string     `gorm:"column:type" json:"type"`
	URL       string     `gorm:"column:url" json:"url"`
	Deleted   bool       `gorm:"column:deleted" json:"-"`
	Comments  []*Comment `gorm:"-" json:"-"`
}

func (p *Post) TableName() string {
	return "post_tab"
}

func (p *Post) FirstOrCreate() error {
	var post Post
	if err := db.Where(&Post{Pid: p.Pid}).FirstOrCreate(&post, p).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Fatalln(err)
			return err
		}
	}

	return nil
}

func (p *Post) UpdateComments() error {
	for _, c := range p.Comments {
		if err := c.FirstOrCreate(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Post) MardAsDeleted() error {
	if err := db.Model(&Post{}).Where("pid = ?", p.Pid).
		Updates(map[string]interface{}{"deleted": true}).Error; err != nil {
		logger.Fatalln(err)
		return err
	}

	return nil
}

func (p *Post) Print() string {
	timestamp, _ := strconv.Atoi(p.Timestamp)
	tm := time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")

	contentStr := fmt.Sprintf("pid=%v, timestamp=%v, text=%v", p.Pid, tm, p.Text)

	if p.Type == "image" {
		contentStr += fmt.Sprintf("\nimage_link: https://pkuhelper.pku.edu.cn/services/pkuhole/images/%v", p.URL)
	}
	contentStr += "\n\n"

	return contentStr
}

func BackUp() {
	var posts []*Post

	if err := db.Where("deleted = ?", true).Find(&posts).Error; err != nil {
		logger.Fatalln(err)
	}
	dump, err := os.OpenFile("deleted_dump.txt", os.O_WRONLY|os.O_CREATE, 0755)
	defer dump.Close()
	if err != nil {
		panic(err)
	}

	for _, p := range posts {
		_, err := dump.WriteString(p.Print())
		if err != nil {
			logger.Println(err)
		}
	}
}

func SelectLatestPosts(startPid string) ([]*Post, error) {
	var posts []*Post
	if err := db.Where("pid >= ?", startPid).Find(&posts).Error; err != nil {
		logger.Fatalln(err)
		return posts, err
	}

	return posts, nil
}

type PostList []*Post

func (pl PostList) Len() int {
	return len(pl)
}

func (pl PostList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func (pl PostList) Less(i, j int) bool {
	ii, err := strconv.Atoi(pl[i].Pid)
	if err != nil {
		panic(err)
	}
	jj, err := strconv.Atoi(pl[j].Pid)
	if err != nil {
		panic(err)
	}
	return ii < jj
}

type Comment struct {
	ID        int    `gorm:"primaryKey;column:id" json:"-"`
	Anonymous string `gorm:"column:anonymous" json:"anonymous"`
	Cid       string `gorm:"column:cid" json:"cid"`
	Islz      int    `gorm:"column:islz" json:"islz"`
	Name      string `gorm:"column:name" json:"name"`
	Pid       string `gorm:"column:pid" json:"pid"`
	Tag       string `gorm:"column:tag" json:"tag"`
	Text      string `gorm:"column:text" json:"text"`
	Timestamp string `gorm:"column:timestamp" json:"timestamp"`
}

func (c *Comment) TableName() string {
	return "comment_tab"
}

func (c *Comment) FirstOrCreate() error {
	var comment Comment
	if err := db.Where(&Comment{Cid: c.Cid}).FirstOrCreate(&comment, c).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Fatalln(err)
			return err
		}
	}

	return nil
}
