# Exercise: A Bookmarking API in Golang

This is a simple bookmarking API written in Go. It uses an SQLite
database to store bookmarks.

## Status

This project is still in development. The API is not yet complete.

## Bookmarklet

```javascript
javascript:(()=>{dest="http://127.0.0.1:8080/me/bookmarks";msg="Notes for link:";note=prompt(msg);url=window.location.href;note&&0!==note.length&&fetch(dest,{method:"POST",body:JSON.stringify({note:note,url:url})}).catch((t=>alert(`Error posting note: ${note}\n\nError: ${t}`)))})();
```