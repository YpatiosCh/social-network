-- Seed data for posts service database
-- User IDs are from 1-10 (references to user service)

------------------------------------------
-- Posts (5 posts from different users)
------------------------------------------
INSERT INTO posts (post_body, creator_id, group_id, audience) VALUES
('Just finished an amazing project! #excited #coding', 1, NULL, 'everyone'),
('Beautiful sunset at the beach today! 🌅', 2, NULL, 'everyone'),
('Working on some exciting features for our app', 3, 1, 'group'),
('Anyone want to grab coffee later?', 4, NULL, 'followers'),
('Finally deployed to production! 🚀', 5, NULL, 'everyone');

------------------------------------------
-- Post Audience (for selected audience posts)
------------------------------------------
INSERT INTO post_audience (post_id, allowed_user_id) VALUES
(4, 1),
(4, 2),
(4, 3),
(4, 6);

------------------------------------------
-- Comments (multiple comments on different posts)
------------------------------------------
INSERT INTO comments (comment_creator_id, parent_id, comment_body) VALUES
(2, 1, 'That sounds awesome! Congrats! 🎉'),
(3, 1, 'Nice work! What stack are you using?'),
(1, 2, 'Looks incredible! Where was this taken?'),
(4, 2, 'I love the colors in this photo'),
(5, 3, 'This is going to be great for the team'),
(1, 4, 'I''m in! What time works for you?'),
(6, 4, 'Me too! Been wanting to catch up'),
(2, 5, 'Amazing! Congrats on the successful deployment!'),
(3, 5, 'How long did this take to build?'),
(7, 5, 'This is incredible! Well done!');

------------------------------------------
-- Reactions (likes on posts and comments)
------------------------------------------
INSERT INTO reactions (content_id, user_id) VALUES
-- Reactions on post 1
(1, 2),
(1, 3),
(1, 4),
(1, 5),
(1, 6),
(1, 7),
-- Reactions on post 2
(2, 1),
(2, 3),
(2, 4),
(2, 5),
(2, 8),
-- Reactions on post 3
(3, 1),
(3, 2),
(3, 4),
(3, 5),
(3, 6),
(3, 9),
-- Reactions on post 4
(4, 1),
(4, 2),
(4, 5),
(4, 7),
-- Reactions on post 5
(5, 1),
(5, 2),
(5, 3),
(5, 4),
(5, 6),
(5, 7),
(5, 8),
(5, 9),
(5, 10),
-- Reactions on comments
(6, 1),
(6, 4),
(6, 5),
(7, 2),
(7, 3),
(8, 1),
(8, 5),
(9, 2),
(9, 3),
(9, 4),
(10, 1),
(10, 6);

------------------------------------------
-- Events (3 group events)
------------------------------------------
INSERT INTO events (event_title, event_body, event_creator_id, group_id, event_date) VALUES
('Team Lunch Meeting', 'Let''s gather for lunch and discuss the upcoming quarter goals', 1, 1, CURRENT_DATE + INTERVAL '7 days'),
('Project Planning Session', 'Sprint planning for next two weeks. Please bring your ideas!', 3, 1, CURRENT_DATE + INTERVAL '14 days'),
('Social Gathering', 'Casual meetup for team members to relax and socialize', 5, 1, CURRENT_DATE + INTERVAL '21 days');

------------------------------------------
-- Event Responses (members responding to events)
------------------------------------------
INSERT INTO event_responses (event_id, user_id, going) VALUES
-- Event 1 responses
(1, 1, true),
(1, 2, true),
(1, 3, true),
(1, 4, false),
(1, 5, true),
(1, 6, true),
(1, 7, false),
(1, 8, true),
(1, 9, true),
(1, 10, true),
-- Event 2 responses
(2, 1, true),
(2, 2, true),
(2, 3, true),
(2, 4, true),
(2, 5, true),
(2, 6, false),
(2, 7, true),
(2, 8, false),
(2, 9, true),
(2, 10, true),
-- Event 3 responses
(3, 1, true),
(3, 2, true),
(3, 3, false),
(3, 4, true),
(3, 5, true),
(3, 6, true),
(3, 7, true),
(3, 8, true),
(3, 9, true),
(3, 10, false);

------------------------------------------
-- Images (optional: add images to posts/comments if needed)
------------------------------------------
-- Images would be referenced from media service
-- Uncomment and adjust if images are needed
-- INSERT INTO images (id, parent_id, sort_order) VALUES
-- (1001, 1, 1),
-- (1002, 2, 1),
-- (1003, 3, 1);
