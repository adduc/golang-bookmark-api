package main

/*

Synopsis:
- A bookmark database (in the vein of Del.icio.us or Pocket)

Goal:
- Create a simple web server that listens on port 8080

Entities:
- "Shared" entities used by mutiple users:
  - Bookmark: The URL, and possibly the title and description if we can
    get a scraper in place
  - Tag: A keyword associated with a bookmark
  - Auth Method: A method of authentication (e.g. email/password)
- User-specific entities:
  - User: A user of the system
  - UserAuth: One of the user's authentication methods (e.g. email/password)
  - UserBookmark: Bookmark metadata associated with a user (e.g. note, tags)
  - List: A collection of bookmarks created by a user
  - BookmarkTag: A tag associated with a bookmark by a user

Database (SQLite3):

- shared tables
   - auth_methods: id, name
   - bookmarks: id, url, title, description
   - tags: id, name
- user tables
   - users: id, username
   - user_auths: id, user_id, method, value
   - user_bookmarks: id, user_id, bookmark_id, note
   - lists: id, user_id, name
   - list_bookmarks: id, list_id, bookmark_id
   - user_bookmark_tags: id, bookmark_id, tag_id, user_id

Routes:
- GET /: Landing Page (Welcome to the Bookmark API)
- GET /me/bookmarks: List of Bookmarks for the user
- POST /me/bookmarks: Add a new bookmark
- GET /me/lists: List of Lists for the user
- GET /me/lists/:list_id: List of Bookmarks for a specific list
- GET /me/tags: List of Tags for the user
- GET /lists: List of all lists
- GET /bookmarks: List of all bookmarks
- GET /tags: List of all tags

Opportunities for Expansion:
- Add a scraper to get the title and description of a URL
- Add a search endpoint
- Social features like following lists or users
- Gamification features like badges or points to encourage usage

*/
