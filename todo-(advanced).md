## Forum Advanced Features - Todo List

### Notification System (@Index.js, @postDetails.js, @Insert.go, @Init.go)
- [X] Creating a notification system DB for user interactions. 
    - [X] Create a new database table to store notifications (@Init.go).
    - [X] Modify the post liking/disliking logic to create a notification entry (@Insert.go, @Likes.go).
    - [X] Modify the comment creation logic to create a notification entry (@Insert.go, @Comments.go).
- [X] Retrieve and display notifications on the user's profile page (@Index.js, @postDetails.js).
- [ ] Implement clearing the notification when user has clicked on the notification, or when the user clicked on (mark all as read).
    - [ ] Create a new endpoint to clear notifications (@Get.go). 
    - [X] Create a new endpoint to mark all notifications as read (@Get.go).
    - [ ] Modify the notification retrieval logic to include the status of the notification (@Get.go).
    - [ ] Adjust the frontend to handle the notification clearing and marking as read.


### Activity Page (@activity.html, @Activity.js, @Activity.go, @Get.go)
- [ ] Create a new HTML page for the user activity (@activity.html).
    - [ ] Design the layout for displaying user-created posts, likes/dislikes, and comments.
- [ ] Implement the backend API endpoint to retrieve user activity data (@Activity.go).
    - [ ] Retrieve user-created posts from the database (@Get.go).
    - [ ] Retrieve posts where the user left a like or dislike (@Get.go).
    - [ ] Retrieve comments made by the user, along with the corresponding post information (@Get.go).
- [ ] Create a new JavaScript file to handle the activity page functionality (@Activity.js).
    - [ ] Fetch user activity data from the backend API.
    - [ ] Display user-created posts, likes/dislikes, and comments on the activity page.

### Edit/Remove Posts and Comments (@postDetails.html, @postDetails.js, @Posts.go, @Comments.go)
    - [ ] Add edit and remove options for posts and comments (@postDetails.html).
    - [ ] Display edit and remove buttons for posts and comments created by the user.
- [ ] Implement the backend API endpoints for editing and removing posts and comments (@Posts.go, @Comments.go).
    - [ ] Create an API endpoint to handle post editing.
    - [ ] Create an API endpoint to handle post removal.
    - [ ] Create an API endpoint to handle comment editing.
    - [ ] Create an API endpoint to handle comment removal.
- [ ] Update the post details page to support editing and removing posts and comments (@postDetails.js).
    - [ ] Add event listeners for edit and remove buttons.
    - [ ] Implement the logic to send edit and remove requests to the backend API.
    - [ ] Update the post and comment display after successful editing or removal.

### Additional Features (optional)
- [ ] Implement user profile pages (@profile.html, @Profile.js, @Profile.go).
    - [ ] Create a new HTML page for user profiles.
    - [ ] Implement the backend API endpoint to retrieve user profile data.
    - [ ] Display user profile information, including created posts and comments.
- [ ] Add search functionality (@search.html, @Search.js, @Search.go).
    - [ ] Create a new HTML page for search results.
    - [ ] Implement the backend API endpoint to handle search queries.
    - [ ] Display search results for posts and comments.
- [ ] Implement pagination for posts and comments (@Index.js, @postDetails.js, @Posts.go, @Comments.go).
    - [ ] Modify the backend API endpoints to support pagination.
    - [ ] Update the frontend to handle paginated results and display pagination controls.
