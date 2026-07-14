var page_limit = 10
var page_offset = 0

$(document).ready(() => {
  fetch_page_content = () => {
    $.ajax({
      url: "/data/all/info/?limit=" + page_limit + "&offset=" + page_offset,
      method: "GET",
      success: (data) => {
        console.log(data)
        for (let i = 0 ; i < data.Posts.length ; i++) {
          new_post = $('<div> <h4>' + data.Posts[i].Title + '</h4> <h5>' + data.Posts[i].Date + '</h5> </div>')
          new_post.on("click", () => {location.href = data.Posts[i].Url})
          new_post.addClass("recent-posting")
          $(".recent-post-links").append(new_post)
        }

        if (!data.More_content_available) {
          $(".recent-posts-button").hide()
        } else {
          page_offset++
        }
      }
    })
  }
  fetch_page_content()
  $(".recent-posts-button").on("click", fetch_page_content)
})