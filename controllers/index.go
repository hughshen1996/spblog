package controllers

import (
	"math"
	"net/http"
	"spblog/conf"
	"spblog/models"
	"spblog/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func IndexGet(c *gin.Context) {
	var (
		pageIndex int
		pageSize  = conf.Con.GetInt("page_size")
		total     int
		page      string
		err       error
		posts     []*models.Post
		policy    *bluemonday.Policy
	)
	page = c.Query("page")
	pageIndex, _ = strconv.Atoi(page)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	posts, err = models.ListPublishedPost("", pageIndex, pageSize)
	if err != nil {
		util.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	total, err = models.CountPostByTag("")
	if err != nil {
		util.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Tags, _ = models.ListTagByPostId(strconv.FormatUint(uint64(post.ID), 10))
		post.Body = policy.Sanitize(string(blackfriday.Run([]byte(post.Body))))
	}
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"posts":           posts,
		"tags":            models.MustListTag(),
		"archives":        models.MustListPostArchives(),
		"links":           models.MustListLinks(),
		"user":            user,
		"pageIndex":       pageIndex,
		"totalPage":       int(math.Ceil(float64(total) / float64(pageSize))),
		"path":            c.Request.URL.Path,
		"maxReadPosts":    models.MustListMaxReadPost(),
		"maxCommentPosts": models.MustListMaxCommentPost(),
	})
}

func AdminIndex(c *gin.Context) {
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "admin/index.html", gin.H{
		"userCount":    models.CountUser(),
		"postCount":    models.CountPost(),
		"tagCount":     models.CountTag(),
		"commentCount": models.CountComment(),
		"user":         user,
		"comments":     models.MustListUnreadComment(),
	})
}
