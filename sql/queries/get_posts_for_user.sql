-- name: GetPostsForUser :many 

SELECT * from posts
ORDER BY published_at DESC 
LIMIT $1; 
