package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Social features for the publishing platform

// Follow represents a user following another user or author
type Follow struct {
	ID        uuid.UUID `json:"id"`
	FollowerID uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Share represents a shared article
type Share struct {
	ID        uuid.UUID `json:"id"`
	ArticleID uuid.UUID `json:"article_id"`
	UserID    uuid.UUID `json:"user_id"`
	Platform  string    `json:"platform"` // twitter, linkedin, facebook, email, copy
	CreatedAt time.Time `json:"created_at"`
}

// Reaction represents various reactions to content
type Reaction struct {
	ID        uuid.UUID `json:"id"`
	ArticleID uuid.UUID `json:"article_id,omitempty"`
	CommentID uuid.UUID `json:"comment_id,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"` // like, love, insightful, agree, disagree
	CreatedAt time.Time `json:"created_at"`
}

// Notification represents a social notification
type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ActorID   uuid.UUID `json:"actor_id"` // Who triggered the notification
	Type      string    `json:"type"`    // comment, like, follow, reply, mention
	ArticleID uuid.UUID `json:"article_id,omitempty"`
	Message   string    `json:"message"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Mention represents a @mention in a comment
type Mention struct {
	ID        uuid.UUID `json:"id"`
	CommentID uuid.UUID `json:"comment_id"`
	MentionedUserID uuid.UUID `json:"mentioned_user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterSocialRoutes registers all social feature endpoints
func RegisterSocialRoutes(router *mux.Router) {
	// Following
	router.HandleFunc("/api/v1/social/users/{id}/follow", FollowUser).Methods("POST")
	router.HandleFunc("/api/v1/social/users/{id}/unfollow", UnfollowUser).Methods("DELETE")
	router.HandleFunc("/api/v1/social/users/{id}/followers", GetFollowers).Methods("GET")
	router.HandleFunc("/api/v1/social/users/{id}/following", GetFollowing).Methods("GET")
	router.HandleFunc("/api/v1/social/users/{id}/is-following", IsFollowing).Methods("GET")

	// Reactions (Extended: like, love, insightful, etc)
	router.HandleFunc("/api/v1/social/articles/{id}/react", ReactToArticle).Methods("POST")
	router.HandleFunc("/api/v1/social/articles/{id}/reactions", GetArticleReactions).Methods("GET")
	router.HandleFunc("/api/v1/social/comments/{id}/react", ReactToComment).Methods("POST")

	// Sharing
	router.HandleFunc("/api/v1/social/articles/{id}/share", ShareArticle).Methods("POST")
	router.HandleFunc("/api/v1/social/articles/{id}/shares", GetArticleShares).Methods("GET")
	router.HandleFunc("/api/v1/social/articles/{id}/share-url", GenerateShareURL).Methods("GET")

	// Mentions
	router.HandleFunc("/api/v1/social/mentions", GetMentions).Methods("GET")
	router.HandleFunc("/api/v1/social/mentions/{id}/read", MarkMentionAsRead).Methods("PUT")

	// Notifications
	router.HandleFunc("/api/v1/social/notifications", GetNotifications).Methods("GET")
	router.HandleFunc("/api/v1/social/notifications/{id}/read", MarkNotificationAsRead).Methods("PUT")
	router.HandleFunc("/api/v1/social/notifications/read-all", MarkAllNotificationsAsRead).Methods("PUT")

	// Network/Discovery
	router.HandleFunc("/api/v1/social/recommendations", GetFollowRecommendations).Methods("GET")
	router.HandleFunc("/api/v1/social/trending-authors", GetTrendingAuthors).Methods("GET")
}

// FollowUser follows another user/author
func FollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followingID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	if userID == followingID {
		respondError(w, http.StatusBadRequest, "Cannot follow yourself")
		return
	}

	followID := uuid.New()
	_, err := db.Exec(`
		INSERT INTO follows (id, follower_id, following_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`, followID, userID, followingID, time.Now())

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to follow user")
		return
	}

	// Create notification
	db.Exec(`
		INSERT INTO notifications (id, user_id, actor_id, type, message, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, uuid.New(), followingID, userID, "follow", fmt.Sprintf("User followed you"), time.Now())

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User followed",
	})
}

// UnfollowUser unfollows a user
func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followingID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	_, err := db.Exec(`
		DELETE FROM follows
		WHERE follower_id = $1 AND following_id = $2
	`, userID, followingID)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to unfollow user")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User unfollowed",
	})
}

// GetFollowers gets list of followers
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	limit := 20

	rows, err := db.Query(`
		SELECT follower_id, created_at
		FROM follows
		WHERE following_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch followers")
		return
	}
	defer rows.Close()

	followers := []map[string]interface{}{}
	for rows.Next() {
		var followerID string
		var createdAt time.Time
		if err := rows.Scan(&followerID, &createdAt); err == nil {
			followers = append(followers, map[string]interface{}{
				"follower_id": followerID,
				"created_at":  createdAt,
			})
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"followers": followers,
		"count":     len(followers),
	})
}

// GetFollowing gets list of users being followed
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	limit := 20

	rows, err := db.Query(`
		SELECT following_id, created_at
		FROM follows
		WHERE follower_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch following")
		return
	}
	defer rows.Close()

	following := []map[string]interface{}{}
	for rows.Next() {
		var followingID string
		var createdAt time.Time
		if err := rows.Scan(&followingID, &createdAt); err == nil {
			following = append(following, map[string]interface{}{
				"following_id": followingID,
				"created_at":   createdAt,
			})
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"following": following,
		"count":     len(following),
	})
}

// IsFollowing checks if user is following another user
func IsFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followingID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM follows
			WHERE follower_id = $1 AND following_id = $2
		)
	`, userID, followingID).Scan(&exists)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to check follow status")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"is_following": exists,
	})
}

// ReactToArticle adds a reaction to an article
func ReactToArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	var req struct {
		Type string `json:"type"` // like, love, insightful, agree, disagree
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Type == "" {
		req.Type = "like"
	}

	reactionID := uuid.New()
	_, err := db.Exec(`
		INSERT INTO reactions (id, article_id, user_id, type, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (article_id, user_id) DO UPDATE
		SET type = $4, created_at = $5
	`, reactionID, articleID, userID, req.Type, time.Now())

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to add reaction")
		return
	}

	// Update article likes count (for "like" type)
	if req.Type == "like" {
		db.Exec("UPDATE articles SET likes = likes + 1 WHERE id = $1 AND NOT EXISTS (SELECT 1 FROM reactions WHERE article_id = $1 AND user_id = $2 AND type = 'like')", articleID, userID)
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Reaction added",
		"type":    req.Type,
	})
}

// GetArticleReactions gets all reactions for an article
func GetArticleReactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]

	rows, err := db.Query(`
		SELECT type, COUNT(*) as count
		FROM reactions
		WHERE article_id = $1
		GROUP BY type
	`, articleID)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch reactions")
		return
	}
	defer rows.Close()

	reactions := map[string]int{}
	for rows.Next() {
		var reactionType string
		var count int
		if err := rows.Scan(&reactionType, &count); err == nil {
			reactions[reactionType] = count
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"reactions": reactions,
		"total":     getTotalReactions(reactions),
	})
}

// ReactToComment adds a reaction to a comment
func ReactToComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["id"]

	var req struct {
		Type string `json:"type"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Type == "" {
		req.Type = "like"
	}

	_, err := db.Exec(`
		UPDATE comments
		SET likes = likes + 1
		WHERE id = $1
	`, commentID)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to react to comment")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Reaction added",
	})
}

// ShareArticle records article share
func ShareArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	var req struct {
		Platform string `json:"platform"` // twitter, linkedin, facebook, email, copy
	}
	json.NewDecoder(r.Body).Decode(&req)

	shareID := uuid.New()
	_, err := db.Exec(`
		INSERT INTO shares (id, article_id, user_id, platform, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, shareID, articleID, userID, req.Platform, time.Now())

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to record share")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message":   "Article shared",
		"platform":  req.Platform,
		"share_id":  shareID.String(),
	})
}

// GetArticleShares gets share count by platform
func GetArticleShares(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]

	rows, err := db.Query(`
		SELECT platform, COUNT(*) as count
		FROM shares
		WHERE article_id = $1
		GROUP BY platform
	`, articleID)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch shares")
		return
	}
	defer rows.Close()

	shares := map[string]int{}
	for rows.Next() {
		var platform string
		var count int
		if err := rows.Scan(&platform, &count); err == nil {
			shares[platform] = count
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"shares": shares,
		"total":  getTotalShares(shares),
	})
}

// GenerateShareURL generates a shareable URL with tracking
func GenerateShareURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	var article struct {
		Slug string
	}
	err := db.QueryRow("SELECT slug FROM articles WHERE id = $1", articleID).Scan(&article.Slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "Article not found")
		return
	}

	shareURL := fmt.Sprintf("https://agennext.com/publishing/article/%s?shared_by=%s", article.Slug, userID)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"share_url": shareURL,
		"platforms": map[string]string{
			"twitter":   fmt.Sprintf("https://twitter.com/intent/tweet?url=%s", shareURL),
			"linkedin":  fmt.Sprintf("https://www.linkedin.com/sharing/share-offsite/?url=%s", shareURL),
			"facebook":  fmt.Sprintf("https://www.facebook.com/sharer/sharer.php?u=%s", shareURL),
		},
	})
}

// GetNotifications gets user notifications
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	limit := 20
	unreadOnly := r.URL.Query().Get("unread") == "true"

	query := `
		SELECT id, actor_id, type, message, read_at, created_at
		FROM notifications
		WHERE user_id = $1
	`
	args := []interface{}{userID}

	if unreadOnly {
		query += " AND read_at IS NULL"
	}

	query += " ORDER BY created_at DESC LIMIT $2"
	args = append(args, limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}
	defer rows.Close()

	notifications := []Notification{}
	for rows.Next() {
		var notif Notification
		rows.Scan(&notif.ID, &notif.ActorID, &notif.Type, &notif.Message, &notif.ReadAt, &notif.CreatedAt)
		notifications = append(notifications, notif)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"notifications": notifications,
		"count":         len(notifications),
		"unread":        countUnreadNotifications(userID),
	})
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["id"]

	now := time.Now()
	_, err := db.Exec(`
		UPDATE notifications
		SET read_at = $1
		WHERE id = $2
	`, now, notificationID)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Notification marked as read",
	})
}

// MarkAllNotificationsAsRead marks all user notifications as read
func MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	now := time.Now()
	_, err := db.Exec(`
		UPDATE notifications
		SET read_at = $1
		WHERE user_id = $2 AND read_at IS NULL
	`, now, userID)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to mark notifications as read")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "All notifications marked as read",
	})
}

// GetFollowRecommendations gets recommended authors to follow
func GetFollowRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	limit := 10

	// Recommend authors based on category subscriptions
	rows, err := db.Query(`
		SELECT DISTINCT a.author
		FROM articles a
		JOIN category_subscriptions cs ON a.category = cs.category_slug
		WHERE cs.user_id = $1
		AND a.author NOT IN (
			SELECT following_id FROM follows WHERE follower_id = $1
		)
		GROUP BY a.author
		ORDER BY COUNT(*) DESC
		LIMIT $2
	`, userID, limit)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch recommendations")
		return
	}
	defer rows.Close()

	recommendations := []string{}
	for rows.Next() {
		var author string
		if err := rows.Scan(&author); err == nil {
			recommendations = append(recommendations, author)
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetTrendingAuthors gets trending authors
func GetTrendingAuthors(w http.ResponseWriter, r *http.Request) {
	limit := 20

	rows, err := db.Query(`
		SELECT author, COUNT(*) as article_count, SUM(views) as total_views
		FROM articles
		WHERE status = 'published' AND published_at > now() - interval '7 days'
		GROUP BY author
		ORDER BY total_views DESC
		LIMIT $1
	`, limit)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch trending authors")
		return
	}
	defer rows.Close()

	authors := []map[string]interface{}{}
	for rows.Next() {
		var author string
		var articleCount int
		var totalViews int
		if err := rows.Scan(&author, &articleCount, &totalViews); err == nil {
			authors = append(authors, map[string]interface{}{
				"author":         author,
				"article_count":  articleCount,
				"total_views":    totalViews,
			})
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"authors": authors,
		"count":   len(authors),
	})
}

// Helper functions

func getTotalReactions(reactions map[string]int) int {
	total := 0
	for _, count := range reactions {
		total += count
	}
	return total
}

func getTotalShares(shares map[string]int) int {
	total := 0
	for _, count := range shares {
		total += count
	}
	return total
}

func countUnreadNotifications(userID string) int {
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM notifications
		WHERE user_id = $1 AND read_at IS NULL
	`, userID).Scan(&count)
	return count
}

// GetMentions gets @mentions for a user
func GetMentions(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	rows, err := db.Query(`
		SELECT m.id, m.comment_id, c.author, c.content, m.created_at
		FROM mentions m
		JOIN comments c ON m.comment_id = c.id
		WHERE m.mentioned_user_id = $1
		ORDER BY m.created_at DESC
		LIMIT 20
	`, userID)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch mentions")
		return
	}
	defer rows.Close()

	mentions := []map[string]interface{}{}
	for rows.Next() {
		var id, commentID, author, content string
		var createdAt time.Time
		if err := rows.Scan(&id, &commentID, &author, &content, &createdAt); err == nil {
			mentions = append(mentions, map[string]interface{}{
				"id":         id,
				"comment_id": commentID,
				"author":     author,
				"content":    content,
				"created_at": createdAt,
			})
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"mentions": mentions,
		"count":    len(mentions),
	})
}

// MarkMentionAsRead marks a mention as read
func MarkMentionAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mentionID := vars["id"]

	_, err := db.Exec(`
		UPDATE mentions SET read_at = $1 WHERE id = $2
	`, time.Now(), mentionID)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to mark mention as read")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Mention marked as read",
	})
}
