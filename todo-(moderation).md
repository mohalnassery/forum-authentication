## Forum Moderation - Todo List

### Update User Roles and Permissions (@Init.go, @Verify.go, @Insert.go)
- [ ] Add user roles for guests, users, moderators, and administrators in the database schema.
    - [ ] Create a new table `user_roles` to store user roles.
    - [ ] Add a foreign key constraint in the `users` table to reference the `user_roles` table.
- [ ] Implement functions to assign and verify user roles.
    - [ ] Create functions to promote or demote users to/from the moderator role.
    - [ ] Create functions to check user permissions based on their role.

### Implement Post Filtering and Moderation (@Posts.go, @Comments.go)
- [ ] Add fields in the `posts` and `comments` tables to track moderation status (e.g., pending, approved, rejected).
- [ ] Implement functions for moderators to review and approve/reject posts and comments.
    - [ ] Create endpoints for moderators to retrieve pending posts and comments.
    - [ ] Create endpoints for moderators to approve or reject posts and comments.
- [ ] Modify post and comment retrieval functions to filter based on moderation status.
    - [ ] Ensure only approved posts and comments are visible to users.
    - [ ] Allow moderators and administrators to view all posts and comments regardless of moderation status.

### Implement Reporting System (@Posts.go, @Comments.go, @Insert.go)
- [ ] Add a new table `reports` to store reported posts and comments.
    - [ ] Include fields for report ID, post/comment ID, reporter ID, reason, and status.
- [ ] Create endpoints for users to report posts and comments.
    - [ ] Allow users to provide a reason for reporting.
    - [ ] Insert reported posts and comments into the `reports` table.
- [ ] Create endpoints for moderators and administrators to view and manage reported content.
    - [ ] Retrieve reported posts and comments for review.
    - [ ] Allow moderators and administrators to take action on reported content (e.g., delete, approve, reject).

### Implement Administrator Functionalities (@Categories.go, @Posts.go, @Comments.go, @Auth.go)
- [ ] Create endpoints for administrators to manage categories.
    - [ ] Allow administrators to create and delete categories.
    - [ ] Modify category retrieval functions to include administrator-only categories.
- [ ] Create endpoints for administrators to manage user roles.
    - [ ] Allow administrators to promote or demote users to/from the moderator role.
- [ ] Create endpoints for administrators to manage reported content.
    - [ ] Allow administrators to view and respond to reports from moderators.
    - [ ] Implement functions for administrators to delete posts and comments.

### Update User Interface (@create-post.html, @index.html, @postDetails.html)
- [ ] Modify the user interface to display user roles and permissions.
    - [ ] Show appropriate actions based on user role (e.g., report button for users, approve/reject buttons for moderators).
- [ ] Add moderation status indicators for posts and comments.
    - [ ] Display pending, approved, or rejected status for posts and comments.
- [ ] Create a dedicated section for moderators and administrators to review and manage reported content.
    - [ ] Display a list of reported posts and comments.
    - [ ] Provide options to take action on reported content.

### Update Client-Side Scripts (@CreatePost.js, @Index.js, @postDetails.js)
- [ ] Modify client-side scripts to handle user roles and permissions.
    - [ ] Show/hide relevant actions based on user role.
    - [ ] Disable certain functionalities for guests and users without appropriate permissions.
- [ ] Implement client-side handling of moderation status.
    - [ ] Display moderation status indicators for posts and comments.
    - [ ] Prevent users from interacting with pending or rejected content.
- [ ] Add client-side functionality for reporting posts and comments.
    - [ ] Allow users to select a reason for reporting.
    - [ ] Send report details to the server-side endpoints.
