# Image Upload Feature - Todo List

## Update Post Creation Form (@create-post.html)
- [X] Add an `<input type="file">` element to the form for image selection.
- [X] Set the `accept` attribute to allow only JPEG, PNG, and GIF file types.
- [X] Add necessary form fields for associated text (title, body, etc.).
## Client-Side Validation (@CreatePost.js)
- [X] Implement client-side validation to check the selected file type (JPEG, PNG, GIF).
- [X] Check the file size and display an error message if it exceeds 20 MB.
- [X] Prevent form submission if the validation fails.

## Server-Side Handling (@Posts.go)
- [X] Update the server-side handler for post creation to handle image uploads.
- [X] Parse the multipart form data to retrieve the uploaded image file.
- [X] Validate the file type and size on the server-side.
- [X] Move the uploaded image to a designated directory on the server for storage.
- [X] Generate a unique filename for the uploaded image to avoid conflicts.
- [X] Store the image file path or URL in the database along with the associated post details.

## Database Schema Update (@Init.go)
- [X] Modify the database schema to include a column for storing the image file path or URL associated with each post.
- [X] Update the necessary SQL queries and statements to handle the image file path or URL.

## Displaying Images (@postDetails.html, @postDetails.js)
- [X] Update the post details page to display the associated image.
- [X] Retrieve the image file path or URL from the database based on the post ID.
- [X] Use the `<img>` tag to display the image on the post details page.
- [X] Apply necessary styling or resizing to the image using CSS for consistent display.

## Error Handling (@CreatePost.js, @Posts.go)
- [ ] Implement error handling for file type mismatch, file size limit exceeded, or upload failures.
- [ ] Display appropriate error messages to the user, informing them about the specific issue.
- [ ] Log any errors or exceptions on the server-side for debugging and monitoring.

## Update Post Listing (@index.html, @Index.js)
- [ ] Modify the post listing page to display the associated image thumbnail for each post.
- [ ] Retrieve the image file path or URL along with other post details from the database.
- [ ] Use the `<img>` tag to display the image thumbnail for each post in the listing.

## Testing
- [ ] Test the image upload feature thoroughly with different scenarios.
- [ ] Upload valid images and verify successful storage and display on the post details page.
- [ ] Attempt to upload unsupported file types and ensure proper error handling.
- [ ] Test with images exceeding the file size limit and check for appropriate error messages.
- [ ] Verify the image thumbnails are displayed correctly on the post listing page.

