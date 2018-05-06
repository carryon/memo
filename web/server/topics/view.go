package topics

import (
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/html-parser"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var viewRoute = web.Route{
	Pattern: res.UrlTopicView + "/" + urlTopicName.UrlPart(),
	Handler: func(r *web.Response) {
		preHandler(r)
		topicRaw := r.Request.GetUrlNamedQueryVariable(urlTopicName.Id)
		topic := html_parser.EscapeWithEmojis(topicRaw)
		var userPkHash []byte
		if auth.IsLoggedIn(r.Session.CookieId) {
			user, err := auth.GetSessionUser(r.Session.CookieId)
			if err != nil {
				r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
				return
			}
			key, err := db.GetKeyForUser(user.Id)
			if err != nil {
				r.Error(jerr.Get("error getting key for user", err), http.StatusInternalServerError)
				return
			}
			userPkHash = key.PkHash
		}
		topicPosts, err := profile.GetPostsForTopic(topic, userPkHash, 0)
		if err != nil {
			r.Error(jerr.Get("error getting topic posts from db", err), http.StatusInternalServerError)
			return
		}
		if len(topicPosts) == 0 {
			r.SetRedirect(res.UrlTopics)
			r.Error(jerr.New("no posts for topic"), http.StatusInternalServerError)
			return
		}
		if len(userPkHash) > 0 {
			err = profile.AttachReputationToPosts(topicPosts)
			if err != nil {
				r.Error(jerr.Get("error attaching reputation to posts", err), http.StatusInternalServerError)
				return
			}
		}
		err = profile.AttachLikesToPosts(topicPosts)
		if err != nil {
			r.Error(jerr.Get("error attaching likes to posts", err), http.StatusInternalServerError)
			return
		}
		var lastLikeId uint
		for _, topicPost := range topicPosts {
			for _, like := range topicPost.Likes {
				if like.Id > lastLikeId {
					lastLikeId = like.Id
				}
			}
		}
		r.Helper["Topic"] = topic
		r.Helper["Posts"] = topicPosts
		r.Helper["FirstPostId"] = topicPosts[0].Memo.Id
		r.Helper["LastPostId"] = topicPosts[len(topicPosts)-1].Memo.Id
		r.Helper["LastLikeId"] = lastLikeId
		r.RenderTemplate(res.TmplTopicView)
	},
}
