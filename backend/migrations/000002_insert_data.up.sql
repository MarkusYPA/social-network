-- All passwords are "pass"

-- Insert users (14 columns, all specified)
INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth, avatar_path, nickname, about_me, is_public, created_at, updated_at, updated_by, status) VALUES
('emma.bauer@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Emma', 'Bauer', '1993-03-22', 'data/uploads/avatars/1747642156208195000.jpeg', 'EmB', 'Loves painting and coffee shops.', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('liam.muller@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Liam', 'Müller', '1987-09-14', 'data/uploads/avatars/1747642184916489000.webp', 'LiamM', 'Avid cyclist and tech enthusiast.', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('sofia.vogel@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Sofia', 'Vogel', '1995-06-30', 'data/uploads/avatars/1747642206647198000.jpeg', 'Sofi', 'Enjoys yoga and travel.', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('noah.fischer@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Noah', 'Fischer', '1990-12-05', 'data/uploads/avatars/1747642227982390000.jpeg', 'NoahF', 'Music lover and amateur chef.', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('olivia.schreiber@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Olivia', 'Schreiber', '1992-08-17', 'data/uploads/avatars/1747642249397877000.webp', 'Liv', 'Bookworm and nature enthusiast.', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('luca.klein@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Luca', 'Klein', '1989-04-11', 'data/uploads/avatars/1747642275289046000.webp', 'LucaK', 'Gamer and coffee addict.', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('mia.wagner@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Mia', 'Wagner', '1997-01-25', 'data/uploads/avatars/1747642297774998000.jpeg', 'MiaW', 'Aspiring photographer.', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('elias.schmidt@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Elias', 'Schmidt', '1986-10-03', 'data/uploads/avatars/1747642317912843000.jpeg', 'Eli', 'Hiker and history buff.', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('lena.herrmann@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Lena', 'Herrmann', '1994-07-19', 'data/uploads/avatars/1747642338716395000.jpeg', 'LenaH', 'Loves baking and gardening.', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('finn.zimmermann@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Finn', 'Zimmermann', '1991-11-08', 'data/uploads/avatars/1747642362945352000.jpeg', 'FinnZ', 'Tech geek and sci-fi fan.', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('amina.ali@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Amina', 'Ali', '1990-05-15', 'data/uploads/avatars/1747642384091203000.jpeg', NULL, NULL, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('yusuf.kamau@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Yusuf', 'Kamau', '1988-11-20', 'data/uploads/avatars/1747642408504644000.jpeg', 'Yus', 'Loves hiking and photography.', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('aziz.bek@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Aziz', 'Bek', '1985-07-10', 'data/default/profile.svg', 'AzizB', NULL, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
('zara.khan@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Zara', 'Khan', '1996-12-25', 'data/default/profile.svg', NULL, 'Enjoys reading and learning new languages.', FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable');


-- Insert posts (8 columns, image_path always included)
INSERT INTO posts (user_id, content, image_path, privacy_level, created_at, updated_at, updated_by, status) VALUES
(1, 'Just enjoyed a beautiful sunset! #nature #peaceful', '/data/uploads/posts/b15c5404-20a3-4e84-a7f1-995687925af0.jpg', 'public', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 'Trying out a new recipe today. Fingers crossed! 🤞', '/data/uploads/posts/7650b44d-172c-4f19-9597-0b034254ec02.jpg', 'almost_private', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 'Thinking about life and everything in between.', '/data/uploads/posts/0c52f09a-525d-4f94-9700-2f0abbe14043.jpg', 'private', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 'Found a hidden gem of a cafe! The coffee is amazing. ☕', '/data/uploads/posts/ca166a1c-2ec1-4074-9dde-1aa7894d7745.jpg', 'public', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 'First attempt at watercolor painting! 🎨', '/data/uploads/posts/b94147a9-bd09-49c8-9e15-9125a370deef.jpg', 'public', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 'Captured a stunning view on my morning run.', '/data/uploads/posts/32e3577c-ad83-42b8-91fa-687ab0a93fa9.jpg', 'public', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 'Exploring a new hiking trail this weekend!', '/data/uploads/posts/f5dbeebd-33f0-4564-ac24-e917bdbfcfa0.jpg', 'almost_private', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 'Interesse im Fotoprojekt! Wie kann ich teilnehmen?', '/data/uploads/posts/12973660-24d1-4fbc-ac24-c1796ae07bf1.jpg', 'almost_private', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 'Gerade ein tolles Sci-Fi-Buch beendet. Empfehlungen?', '/data/uploads/posts/83b260eb-6399-4735-a4f8-73094a794544.jpg', 'public', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 'Abendstimmung mit Jazzmusik.', '/data/uploads/posts/627013d0-d8d9-4493-9bbc-8996160554b8.jpg', 'private', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 'Heute Sauerteigbrot gebacken. Riecht toll! 🍞', '/data/uploads/posts/54ab6009-3b23-4e9d-88db-337d0ba96c05.jpg', 'almost_private', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 'Neues Coding-Projekt läuft!', '/data/uploads/posts/ec5051be-8f3c-4ef0-9313-854ecbd90e83.jpg', 'public', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 'Plane ein Wochenendabenteuer. Vorschläge?', '/data/uploads/posts/4c5baa2a-d900-47eb-9b96-be97f680ddc4.jpg', 'public', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 'Neue Pflanze für meine Sammlung adoptiert! 🌿', '/data/uploads/posts/c6187ce7-e35d-4ee5-89ec-cc983fd3ab7e.jpg', 'almost_private', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, 'enable');

-- Insert comments (8 columns, image_path always included)
INSERT INTO comments (post_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(1, 2, 'Wow, that sunset sounds amazing!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 6, 'What’s the recipe? Looks tasty!', '/data/uploads/comments/7cc8201b-abc2-48d3-95d6-a4109e9c4b51.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 5, 'Love the colors in your painting!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 8, 'Share the cafe name, please!', '/data/uploads/comments/157149b2-d0cb-4451-a3a3-3c9dbdedf4e2.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 7, 'Which trail are you hiking?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 9, 'I recommend "Dune" if you haven’t read it!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 11, 'That bread looks delicious!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 13, 'Any tips for sourdough starters?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 1, 'Try the coast, it’s beautiful this time of year!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Insert follow requests (5 columns)
INSERT INTO follow_requests (follower_id, followed_id, approval_status, created_at) VALUES
(2, 1, 'pending', CURRENT_TIMESTAMP),
(3, 1, 'pending', CURRENT_TIMESTAMP),
(4, 1, 'pending', CURRENT_TIMESTAMP),
(5, 1, 'pending', CURRENT_TIMESTAMP),
(6, 1, 'pending', CURRENT_TIMESTAMP),
(7, 1, 'pending', CURRENT_TIMESTAMP),
(8, 1, 'pending', CURRENT_TIMESTAMP),
(7, 4, 'accepted', CURRENT_TIMESTAMP),
(9, 6, 'declined', CURRENT_TIMESTAMP),
(11, 8, 'pending', CURRENT_TIMESTAMP),
(13, 8, 'accepted', CURRENT_TIMESTAMP),
(1, 13, 'pending', CURRENT_TIMESTAMP),
(11, 13, 'accepted', CURRENT_TIMESTAMP),
(3, 11, 'pending', CURRENT_TIMESTAMP);

-- Insert 15 groups (7 columns)
INSERT INTO groups (creator_id, title, description, created_at, updated_by, status) VALUES
(1, 'Nature Lovers', 'A group for hiking, camping, and nature enthusiasts.', CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 'Art Enthusiasts', 'Share and discuss your artwork and inspirations.', CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 'Book Club', 'For book lovers to discuss their favorite reads.', CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 'Photography Hub', 'Connect with fellow photographers and share tips.', CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 'Fitness Fanatics', 'For those passionate about fitness and health.', CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 'Baking Buddies', 'Share recipes and baking tips.', CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 'Travel Explorers', 'Discuss travel destinations and experiences.', CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 'Music Makers', 'For musicians and music lovers to connect.', CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 'Tech Talk', 'Discuss the latest in tech and innovation.', CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 'Yoga & Meditation', 'A space for mindfulness and yoga enthusiasts.', CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 'Gamers Guild', 'For gamers to share tips and play together.', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 'Plant Parents', 'For plant lovers to share care tips and photos.', CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 'History Buffs', 'Discuss historical events and share knowledge.', CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 'Code Crafters', 'For programmers to share projects and ideas.', CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 'Coffee Connoisseurs', 'All about coffee brewing and cafe culture.', CURRENT_TIMESTAMP, NULL, 'enable');

-- Insert group members (7 columns)
INSERT INTO group_members (group_id, user_id, approval_status, created_at, updated_by, status) VALUES
(1, 2, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 7, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 13, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 5, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 8, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 1, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 9, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 4, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 11, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 2, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 6, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 12, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 7, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 10, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 3, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 11, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 14, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 5, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 4, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 9, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 2, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 8, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 13, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 6, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 13, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 1, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 7, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 6, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 11, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 4, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 10, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 5, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 8, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 14, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 2, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 9, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 3, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 12, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 1, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 12, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 7, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 10, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 1, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 4, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 6, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 1, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable');

-- Insert group invitations (8 columns)
INSERT INTO group_invitations (group_id, user_id, inviter_id, approval_status, created_at, updated_by, status) VALUES
(1, 3, 1, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 6, 5, 'accepted', CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 11, 9, 'declined', CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 13, 2, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 1, 7, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 8, 11, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 5, 4, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 9, 8, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 2, 13, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 3, 6, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 4, 10, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 7, 14, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 11, 3, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 1, 12, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 11, 12, 'pending', CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 8, 1, 'pending', CURRENT_TIMESTAMP, NULL, 'enable');

-- Insert group posts (8 columns, image_path always included, total 240 posts)
-- Group 1: Nature Lovers (16 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(1, 2, 'Planning a hike this weekend, anyone join?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 7, 'Best camping spots near the city?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 13, 'Just saw a rare bird on my last hike!', '/data/uploads/posts/2e9214de-645b-4c26-9bb1-78eaa1550761.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 2, 'Any tips for beginner hikers?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 7, 'Love the autumn colors in the forest.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 1, 'Group hike next month, who’s in?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 2, 'Found a great trail guide online.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 7, 'Anyone tried wild foraging?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 13, 'Need recommendations for hiking boots.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 1, 'Sunrise hikes are the best!', '/data/uploads/posts/4b5f8ac9-d5e3-40bb-b88e-6ba0f351fc2a.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 2, 'Anyone been to the national park recently?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 7, 'Tips for staying safe in bear country?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 13, 'Favorite nature podcasts?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 1, 'Planning a stargazing night, join us!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 2, 'Just joined a local conservation group.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(1, 7, 'What’s your favorite trail snack?', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 2: Art Enthusiasts (18 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(2, 5, 'Sharing my latest sketch. Feedback welcome!', '/data/uploads/posts/4e6cfacf-b209-4f72-90aa-c2cb9996d291.webp', CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 8, 'New painting in progress.', '/data/uploads/posts/5f50abfa-b572-4917-8cbe-6ed4c20cad4a.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 1, 'Trying oil paints for the first time.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 5, 'Any watercolor tutorials you recommend?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 8, 'Visited an art gallery today, so inspired!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 1, 'Working on a portrait, tips appreciated.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 5, 'Best brushes for acrylic painting?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 8, 'Just finished a digital art piece.', '/data/uploads/posts/6ea40dca-c490-4e1a-87e8-1cc789b98097.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 1, 'Anyone doing inktober this year?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 5, 'Art challenge: draw something inspired by nature.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 8, 'What’s your favorite art style?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 1, 'Planning an art meetup, who’s interested?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 5, 'Struggling with perspective, any tips?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 8, 'Loving these new colored pencils!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 1, 'Anyone tried sculpting?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 5, 'Best way to frame artwork?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 8, 'Art supply sale at the local store!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 1, 'Sketching at the park today, join me!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 3: Book Club (15 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(3, 9, 'Just finished "1984". Let’s discuss!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 4, 'What’s everyone reading this month?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 11, 'Loved "The Hobbit", any similar books?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 9, 'Book club meeting next week, RSVP!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 4, 'Anyone read historical fiction?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 11, 'Just started a new mystery novel.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 9, 'Favorite book-to-movie adaptations?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 4, 'Looking for sci-fi recommendations.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 11, 'Poetry fans, share your favorites!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 9, 'Just joined an audiobook subscription.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 4, 'What’s your go-to cozy read?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 11, 'Anyone read "Dune" yet?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 9, 'Planning a book swap event.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 4, 'Best bookstores in town?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 11, 'Started a reading challenge, join me!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 4: Photography Hub (17 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(4, 2, 'Captured a sunset yesterday.', '/data/uploads/posts/37f5688d-0bc0-419f-9786-35abad695ea9.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 6, 'Tips for low-light photography?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 12, 'Best lenses for portrait shots?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 2, 'Editing software recommendations?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 6, 'Sharing my latest street photography.', '/data/uploads/posts/49b992e7-da8a-4f05-937f-bb382a21ce92.webp', CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 12, 'Anyone tried drone photography?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 2, 'Photography workshop next month!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 6, 'How to get crisp macro shots?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 12, 'Favorite photo spots in the city?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 2, 'Just got a new tripod, love it!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 6, 'Black and white vs. color, thoughts?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 12, 'Sharing a time-lapse project.', '/data/uploads/posts/77ca9bf9-2150-4a70-9479-5a8c3591581c.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 2, 'Anyone doing 365 photo challenges?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 6, 'Best way to organize photos?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 12, 'Night sky photography tips?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 2, 'Photo contest coming up, join us!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 6, 'Just printed my first photo book!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 5: Fitness Fanatics (16 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(5, 7, 'New PR on my deadlift today!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 10, 'Best post-workout meals?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 3, 'Anyone doing HIIT workouts?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 7, 'Running group meetup this weekend.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 10, 'Tips for staying motivated?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 3, 'Just tried kickboxing, so fun!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 7, 'Favorite fitness apps?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 10, 'Sharing my gym playlist.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 3, 'How to improve flexibility?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 7, 'Group cycling event next month.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 10, 'Best shoes for cross-training?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 3, 'Anyone tried rock climbing?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 7, 'Post-workout stretches you love?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 10, 'Meal prep ideas for the week?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 3, 'Fitness challenge: 30 days of planks!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 7, 'Just hit my step goal for the month!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 6: Baking Buddies (18 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(6, 11, 'Sourdough loaf turned out great!', '/data/uploads/posts/239f5552-cefc-42e2-9032-cb88d498ad7a.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 14, 'Best cookie recipes?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 5, 'Tips for perfect pie crust?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 11, 'Baking a cake for a party, ideas?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 14, 'Sharing my muffin recipe.', '/data/uploads/posts/8376db0d-1770-46bb-9f14-e1c271da6c91.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 5, 'Gluten-free baking tips?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 11, 'Favorite baking tools?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 14, 'Just made croissants, so proud!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 5, 'How to prevent soggy bottoms?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 11, 'Baking class next weekend, join us!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 14, 'Best frosting techniques?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 5, 'Anyone tried vegan baking?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 11, 'Sharing my brownie recipe.', '/data/uploads/posts/38948d9f-d837-4a93-a8e2-5560d6793f9d.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 14, 'Bread baking tips for beginners?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 5, 'Favorite bakeware brands?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 11, 'Planning a bake sale, any tips?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 14, 'Just mastered macarons!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 5, 'What’s your go-to dessert?', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 7: Travel Explorers (15 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(7, 4, 'Best budget travel destinations?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 9, 'Sharing my Italy itinerary.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 2, 'Tips for solo travel?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 4, 'Favorite travel apps?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 9, 'Just booked a trip to Japan!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 2, 'How to pack light?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 4, 'Group trip planning, ideas?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 9, 'Best travel podcasts?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 2, 'Anyone been to Iceland?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 4, 'Sharing my travel journal.', '/data/uploads/posts/39606ad7-7a17-424c-ba7c-2626d1fa6e74.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 9, 'Tips for jet lag?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 2, 'Favorite travel blogs?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 4, 'Planning a road trip, suggestions?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 9, 'Best travel gear?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 2, 'Travel meetup next month!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 8: Music Makers (16 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(8, 8, 'New song I’m working on.', '/data/uploads/posts/990483d8-79be-48c9-b053-385ea18ace4f.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 13, 'Best guitars for beginners?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 6, 'Favorite music genres?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 8, 'Just learned a new chord progression.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 13, 'Sharing my latest cover.', '/data/uploads/posts/9584590f-8625-413c-9e77-c353624cbed4.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 6, 'Tips for stage performance?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 8, 'Anyone play the drums?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 13, 'Music theory resources?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 6, 'Planning a jam session, join us!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 8, 'Best music production software?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 13, 'Favorite concert memories?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 6, 'Just got new studio monitors!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 8, 'Tips for writing lyrics?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 13, 'Sharing my playlist.', '/data/uploads/posts/b939624b-3114-4883-adff-f10143a16b22.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 6, 'Anyone in a band?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 8, 'Music open mic next week!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 9: Tech Talk (17 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(9, 13, 'Latest AI news, thoughts?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 1, 'Best laptops for coding?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 7, 'Anyone tried VR gaming?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 13, 'Sharing my app idea.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 1, 'Tips for cybersecurity?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 7, 'Favorite tech podcasts?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 13, 'Just built a new PC!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 1, 'Cloud computing trends?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 7, 'Best coding bootcamps?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 13, 'Sharing my GitHub project.', '/data/uploads/posts/bbe6b7ca-3c67-4d20-9a2c-bdf0be7950d8.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 1, 'Anyone into robotics?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 7, 'Tips for tech interviews?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 13, 'Favorite tech YouTubers?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 1, 'Planning a hackathon, join us!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 7, 'Best smart home devices?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 13, 'Just got a Raspberry Pi!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 1, 'What’s next for 5G?', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 10: Yoga & Meditation (15 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(10, 6, 'Favorite morning yoga routine?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 11, 'Tips for meditation beginners?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 4, 'Best yoga mats?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 6, 'Sharing my mindfulness journal.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 11, 'Just tried a new yoga flow.', '/data/uploads/posts/c6c67abe-8c6f-42ca-a09b-14b285e39fea.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 4, 'Favorite meditation apps?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 6, 'Group yoga session this weekend!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 11, 'How to stay consistent with yoga?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 4, 'Best breathing exercises?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 6, 'Anyone tried hot yoga?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 11, 'Sharing my gratitude practice.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 4, 'Tips for yoga at home?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 6, 'Favorite yoga YouTubers?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 11, 'Just hit 30 days of meditation!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 4, 'Planning a mindfulness retreat.', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 11: Gamers Guild (16 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(11, 10, 'New game I’m playing, love it!', '/data/uploads/posts/c98cfb51-2a46-4a2d-bafa-6485562bb0f6.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 5, 'Best gaming headsets?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 8, 'Anyone play RPGs?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 10, 'Tips for streaming on Twitch?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 5, 'Sharing my gaming setup.', '/data/uploads/posts/c366c1cb-4e20-4ec6-99f8-def2fd2ade59.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 8, 'Favorite esports teams?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 10, 'Just beat a tough boss!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 5, 'Best co-op games?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 8, 'Planning a LAN party, join us!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 10, 'Tips for game modding?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 5, 'Favorite retro games?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 8, 'Just got a new GPU!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 10, 'Sharing my speedrun.', '/data/uploads/posts/cb7876f3-c216-4383-9f7e-f8ce2448ab8e.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 5, 'Best gaming keyboards?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 8, 'Anyone play survival games?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 10, 'Game night next week!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 12: Plant Parents (18 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(12, 14, 'New monstera cutting!', '/data/uploads/posts/ed4ead42-4ed0-48c2-a683-01ad8c94ef7a.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 2, 'Tips for succulent care?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 9, 'Best pots for indoor plants?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 14, 'Sharing my plant shelf.', '/data/uploads/posts/f7f67a2d-477c-469e-8170-890cbcf3371e.webp', CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 2, 'How to deal with pests?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 9, 'Favorite plant shops?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 14, 'Just repotted my fiddle leaf fig.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 2, 'Tips for low-light plants?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 9, 'Anyone propagate plants?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 14, 'Sharing my terrarium.', '/data/uploads/posts/f9e3f143-6702-46b2-b141-efce381b0747.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 2, 'Best fertilizers for houseplants?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 9, 'Just got a new orchid!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 14, 'Tips for plant watering schedules?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 2, 'Favorite plant Instagram accounts?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 9, 'Planning a plant swap event.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 14, 'How to revive a dying plant?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 2, 'Best grow lights?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 9, 'Just joined a plant care workshop!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 13: History Buffs (15 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(13, 3, 'Favorite historical documentaries?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 12, 'Just visited a museum, so cool!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 1, 'Tips for researching family history?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 3, 'Sharing my WWII book list.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 12, 'Anyone into medieval history?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 1, 'Best history podcasts?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 3, 'Just watched a documentary on ancient Egypt.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 12, 'Favorite historical figures?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 1, 'Planning a history trivia night.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 3, 'Tips for visiting historical sites?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 12, 'Sharing my Roman history notes.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 1, 'Best history YouTubers?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 3, 'Just joined a history book club.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 12, 'Favorite historical movies?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 1, 'Anyone study local history?', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 14: Code Crafters (16 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(14, 12, 'New Python project I’m working on.', '/data/uploads/posts/f74ada9b-5c91-4bcc-b6f6-de590e9c2dc7.png', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 7, 'Best IDEs for coding?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 10, 'Tips for debugging code?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 10, 'Sharing my web app.', '/data/uploads/posts/0d3e702f-8f6f-4110-8610-5f6fa4f034de.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 7, 'Favorite coding tutorials?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 10, 'Just learned a new algorithm.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 12, 'Planning a coding meetup.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 7, 'Best resources for learning JavaScript?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 10, 'Anyone use GitLab?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 12, 'Sharing my CLI tool.', '/data/uploads/posts/1deffb0c-6b26-48a9-b074-3e7be6a675b6.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 7, 'Tips for clean code?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 10, 'Favorite coding podcasts?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 12, 'Just deployed my first app!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 7, 'Best free coding courses?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 10, 'Anyone into open-source projects?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 12, 'Coding challenge next week!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Group 15: Coffee Connoisseurs (17 posts)
INSERT INTO group_posts (group_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(15, 1, 'Best coffee beans for espresso?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 4, 'Sharing my pour-over recipe.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 6, 'Favorite coffee shops in town?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 1, 'Just got a new coffee grinder!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 4, 'Tips for cold brew?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 6, 'Best coffee makers?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 1, 'Planning a coffee tasting event.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 4, 'Favorite coffee pairings?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 6, 'Just tried a new roast.', '/data/uploads/posts/b45d820c-9dd1-485b-88d1-dd61351bdd6c.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 1, 'Tips for latte art?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 4, 'Best decaf options?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 6, 'Sharing my coffee setup.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 1, 'Favorite coffee YouTubers?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 4, 'Just joined a coffee subscription.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 6, 'Best milk alternatives for coffee?', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 1, 'Planning a cafe crawl.', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 4, 'What’s your coffee order?', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Insert group comments (7 columns, image_path always included)
INSERT INTO group_comments (group_post_id, user_id, content, image_path, created_at, updated_by, status) VALUES
(1, 1, 'Count me in for the hike!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 7, 'Great camping tips!', '/data/uploads/comments/b2f70289-bd01-4c73-9ba1-6a6fe82373bd.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(17, 8, 'Love the shading in your sketch!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(18, 1, 'Awesome digital art!', '/data/uploads/comments/b57ee776-35d9-45a0-a3c2-edec51f98b0c.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(33, 9, '"1984" is so thought-provoking!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(49, 6, 'That sunset shot is stunning!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(50, 12, 'Love the street photography!', '/data/uploads/comments/b3119bcc-ce09-4e36-a0a2-eb6b9b9f53f3.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(65, 10, 'Congrats on the PR!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(81, 5, 'That loaf looks perfect!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(82, 14, 'Those muffins look delicious!', '/data/uploads/comments/f56d62d4-3c1d-46ad-8196-167f1e2af361.png', CURRENT_TIMESTAMP, NULL, 'enable'),
(97, 9, 'Italy sounds amazing!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(113, 6, 'Love that chord progression!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(129, 7, 'VR gaming is so immersive!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(130, 11, 'Cool app idea!', '/data/uploads/comments/fb0ac7e2-2bf7-4d13-9ca3-0291aed637d7.jpg', CURRENT_TIMESTAMP, NULL, 'enable'),
(145, 11, 'Love that yoga flow!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(161, 8, 'Great gaming setup!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(177, 9, 'That monstera is thriving!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(193, 1, 'That WWII list is great!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(209, 10, 'Love that Python project!', NULL, CURRENT_TIMESTAMP, NULL, 'enable'),
(225, 6, 'That pour-over recipe is spot-on!', NULL, CURRENT_TIMESTAMP, NULL, 'enable');

-- Insert events (9 columns)
INSERT INTO events (group_id, creator_id, title, description, event_datetime, created_at, updated_by, status) VALUES
(1, 1, 'Group Hike', 'A fun hike in the nearby mountains.', '2025-05-20 08:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(2, 5, 'Art Workshop', 'Learn new painting techniques.', '2025-05-25 14:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(3, 9, 'Book Discussion', 'Discuss "1984" by George Orwell.', '2025-05-30 18:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(4, 2, 'Photography Walk', 'Explore the city with cameras.', '2025-06-01 10:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(5, 7, 'Group Run', '5K run in the park.', '2025-06-05 07:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(6, 11, 'Baking Class', 'Learn to make sourdough.', '2025-06-10 15:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(7, 4, 'Travel Meetup', 'Share travel stories.', '2025-06-15 19:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(8, 8, 'Jam Session', 'Open mic for musicians.', '2025-06-20 20:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(9, 13, 'Tech Talk', 'Discuss AI advancements.', '2025-06-25 18:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(10, 6, 'Yoga Retreat', 'Day of yoga and meditation.', '2025-06-30 09:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(11, 10, 'Game Night', 'Board and video games.', '2025-07-05 17:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(12, 14, 'Plant Swap', 'Exchange plant cuttings.', '2025-07-10 11:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(13, 3, 'History Trivia', 'Test your history knowledge.', '2025-07-15 19:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(14, 12, 'Code Hackathon', 'Build a project in 24 hours.', '2025-07-20 10:00:00', CURRENT_TIMESTAMP, NULL, 'enable'),
(15, 1, 'Coffee Tasting', 'Sample local roasts.', '2025-07-25 16:00:00', CURRENT_TIMESTAMP, NULL, 'enable');

-- Insert event responses (5 columns)
INSERT INTO event_responses (event_id, user_id, response, created_at) VALUES
(1, 2, 'going', CURRENT_TIMESTAMP),
(1, 7, 'going', CURRENT_TIMESTAMP),
(1, 13, 'not_going', CURRENT_TIMESTAMP),
(2, 5, 'pending', CURRENT_TIMESTAMP),
(2, 8, 'going', CURRENT_TIMESTAMP),
(2, 1, 'not_going', CURRENT_TIMESTAMP),
(3, 9, 'pending', CURRENT_TIMESTAMP),
(3, 4, 'going', CURRENT_TIMESTAMP),
(3, 11, 'not_going', CURRENT_TIMESTAMP),
(4, 2, 'pending', CURRENT_TIMESTAMP),
(4, 6, 'not_going', CURRENT_TIMESTAMP),
(4, 12, 'going', CURRENT_TIMESTAMP),
(5, 7, 'pending', CURRENT_TIMESTAMP),
(5, 10, 'not_going', CURRENT_TIMESTAMP),
(5, 3, 'going', CURRENT_TIMESTAMP),
(6, 11, 'pending', CURRENT_TIMESTAMP),
(6, 14, 'going', CURRENT_TIMESTAMP),
(6, 5, 'going', CURRENT_TIMESTAMP),
(7, 4, 'pending', CURRENT_TIMESTAMP),
(7, 9, 'going', CURRENT_TIMESTAMP),
(7, 2, 'not_going', CURRENT_TIMESTAMP),
(8, 8, 'pending', CURRENT_TIMESTAMP),
(8, 13, 'going', CURRENT_TIMESTAMP),
(8, 6, 'going', CURRENT_TIMESTAMP),
(9, 13, 'pending', CURRENT_TIMESTAMP),
(9, 1, 'going', CURRENT_TIMESTAMP),
(9, 7, 'not_going', CURRENT_TIMESTAMP),
(10, 6, 'pending', CURRENT_TIMESTAMP),
(10, 11, 'going', CURRENT_TIMESTAMP),
(10, 4, 'going', CURRENT_TIMESTAMP),
(11, 10, 'pending', CURRENT_TIMESTAMP),
(11, 5, 'going', CURRENT_TIMESTAMP),
(11, 8, 'not_going', CURRENT_TIMESTAMP),
(12, 14, 'pending', CURRENT_TIMESTAMP),
(12, 2, 'going', CURRENT_TIMESTAMP),
(12, 9, 'going', CURRENT_TIMESTAMP),
(13, 3, 'pending', CURRENT_TIMESTAMP),
(13, 12, 'going', CURRENT_TIMESTAMP),
(13, 1, 'not_going', CURRENT_TIMESTAMP),
(14, 12, 'pending', CURRENT_TIMESTAMP),
(14, 7, 'going', CURRENT_TIMESTAMP),
(14, 10, 'going', CURRENT_TIMESTAMP),
(15, 1, 'pending', CURRENT_TIMESTAMP),
(15, 4, 'going', CURRENT_TIMESTAMP),
(15, 6, 'going', CURRENT_TIMESTAMP),
(1, 1, 'pending', CURRENT_TIMESTAMP);

-- Insert messages (8 columns, 20-message dialogue between Emma and Yusuf, plus 15 more)
INSERT INTO messages (sender_id, receiver_id, content, created_at, updated_at, updated_by, status) VALUES
-- Dialogue: Planning a group hike (Emma: user_id 5, Yusuf: user_id 2)
(5, 2, 'Hey Liam, saw your post about the group hike. I’m interested! Any details yet?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'Hi Olivia! Awesome, glad you’re in. Thinking about the Blue Ridge Trail, maybe 10am Saturday?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 2, 'Blue Ridge sounds great! How long is the trail?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'It’s about 6 miles round trip, moderate difficulty. Should take 3-4 hours.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 2, 'Perfect, I can handle that. Should we bring lunch or snacks?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'Snacks and water for sure. Maybe a sandwich if we want to picnic at the top.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 2, 'Nice, I’ll pack some granola bars and a PB&J. How many people are coming?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'So far, you, me, Amina, and Sofia. Still waiting on a few others.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 2, 'Cool, sounds like a fun group. Any carpool plans?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'Yeah, I can drive. Want to meet at my place at 9:30am?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 2, 'Works for me! I’ll bring some coffee for the road. 😄', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'You’re a lifesaver! I’ll send you my address.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 2, 'Got it. Should we bring anything else, like hiking poles or bug spray?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'Bug spray’s a good idea, it’s been humid. Poles if you use them, but the trail’s not too steep.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 2, 'Sweet, I’ll grab some spray. Weather looking good?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'Forecast says sunny, mid-70s. Perfect hiking weather!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 2, 'Awesome, I’m pumped! Let me know if anything changes.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'Will do. See you Saturday!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 2, 'See ya then! Thanks for organizing.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'No prob, it’s gonna be a blast!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
-- Additional messages
(9, 4, 'Hey Noah, loved your cafe post! What’s the name of that place?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(11, 8, 'Elias, your photography post was amazing. Any editing tips?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(7, 10, 'Finn, what game are you playing? Saw your gaming post.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(13, 1, 'Emma, your sunset post was beautiful! Where was that?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(3, 5, 'Olivia, how’s the painting going? Saw your art post.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(6, 12, 'Yusuf, your coding project looks cool. What’s it do?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(4, 9, 'Lena, read "1984" yet? Saw your book club post.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(8, 2, 'Liam, that recipe looked tasty! Can you share it?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 7, 'Mia, which trail are you hiking? Saw your post.', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(10, 11, 'Amina, your bread looks perfect. Sourdough tips?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(12, 6, 'Luca, saw your running post. Best shoes for trails?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(14, 13, 'Aziz, your getaway post got me thinking. Beach or mountains?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 14, 'Zara, your plant post was cool. How do you care for it?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 3, 'Sofia, saw your reflection post. Deep thoughts! Wanna chat?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(9, 1, 'Emma, joining the coffee tasting? Saw your event post.', CURRENT_TIMESTAMP, NULL, NULL, 'enable');

-- Insert group messages (8 columns)
INSERT INTO group_messages (group_id, sender_id, content, created_at, updated_at, updated_by, status) VALUES
(1, 1, 'Excited for the group hike! Who’s bringing snacks?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 5, 'Who’s joining the art workshop? It’s gonna be fun!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(3, 9, 'Looking forward to the book discussion!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(4, 2, 'Photography walk this weekend, who’s in?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 7, 'Group run tomorrow, meet at the park!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(6, 11, 'Baking class sign-ups open, grab a spot!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(7, 4, 'Travel meetup next week, share your stories!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(8, 8, 'Jam session Friday, bring your instruments!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(9, 13, 'Tech talk on AI, don’t miss it!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(10, 6, 'Yoga retreat spots filling up, RSVP!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(11, 10, 'Game night Saturday, what games to play?', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(12, 14, 'Plant swap Sunday, bring your cuttings!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(13, 3, 'History trivia night, test your knowledge!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(14, 12, 'Hackathon next month, form your teams!', CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(15, 1, 'Coffee tasting event, who’s coming?', CURRENT_TIMESTAMP, NULL, NULL, 'enable');

-- Insert notifications (12 columns: user_id, type, follow_req_id, group_invite_id, group_members_id, event_id, content, is_read, created_at, updated_at, updated_at, status)
INSERT INTO notifications (user_id, type, follow_req_id, group_invite_id, group_members_id, event_id, content, is_read, created_at, updated_at, updated_by, status) VALUES
-- Follow request notifications (aligned with follow_requests table)
(1, 'follow_request', 1, NULL, NULL, NULL, 'Liam wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 'follow_request', 2, NULL, NULL, NULL, 'Sofia wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 'follow_request', 3, NULL, NULL, NULL, 'Noah wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 'follow_request', 4, NULL, NULL, NULL, 'Olivia wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 'follow_request', 5, NULL, NULL, NULL, 'Luca wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 'follow_request', 6, NULL, NULL, NULL, 'Mia wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 'follow_request', 7, NULL, NULL, NULL, 'Elias wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(4, 'follow_request', 8, NULL, NULL, NULL, 'Mia wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'delete'), -- accepted
(6, 'follow_request', 9, NULL, NULL, NULL, 'Lena wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'delete'),    -- declined
(8, 'follow_request', 10, NULL, NULL, NULL, 'Amina wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(8, 'follow_request', 11, NULL, NULL, NULL, 'Amina wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'delete'), -- accepted
(13, 'follow_request', 12, NULL, NULL, NULL, 'Emma wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(13, 'follow_request', 13, NULL, NULL, NULL, 'Amina wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'delete'), -- accepted
(11, 'follow_request', 14, NULL, NULL, NULL, 'Sofia wants to follow you.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
-- Group invitation notifications (corrected to use group_invite_id)
(3, 'group_invitation', NULL, 1, NULL, NULL, 'You’ve been invited to Nature Lovers.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(6, 'group_invitation', NULL, 2, NULL, NULL, 'You’ve been invited to Art Enthusiasts.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'delete'), -- accepted
(11, 'group_invitation', NULL, 3, NULL, NULL, 'You’ve been invited to Book Club.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'delete'),      -- declined
(13, 'group_invitation', NULL, 4, NULL, NULL, 'You’ve been invited to Photography Hub.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 'group_invitation', NULL, 5, NULL, NULL, 'You’ve been invited to Fitness Fanatics.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(8, 'group_invitation', NULL, 6, NULL, NULL, 'You’ve been invited to Baking Buddies.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 'group_invitation', NULL, 7, NULL, NULL, 'You’ve been invited to Travel Explorers.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(9, 'group_invitation', NULL, 8, NULL, NULL, 'You’ve been invited to Music Makers.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 'group_invitation', NULL, 9, NULL, NULL, 'You’ve been invited to Tech Talk.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(3, 'group_invitation', NULL, 10, NULL, NULL, 'You’ve been invited to Yoga & Meditation.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(4, 'group_invitation', NULL, 11, NULL, NULL, 'You’ve been invited to Gamers Guild.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(7, 'group_invitation', NULL, 12, NULL, NULL, 'You’ve been invited to Plant Parents.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(11, 'group_invitation', NULL, 13, NULL, NULL, 'You’ve been invited to History Buffs.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 'group_invitation', NULL, 14, NULL, NULL, 'You’ve been invited to Code Crafters.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(11, 'group_invitation', NULL, 15, NULL, NULL, 'You’ve been invited to group 14.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(8, 'group_invitation', NULL, 16, NULL, NULL, 'You’ve been invited to Coffee Connoisseurs.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
-- Group join request notifications (aligned with group_members pending requests)
(1, 'group_join_request', NULL, NULL, 3, NULL, 'Amina wants to join Nature Lovers.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
-- (5, 'group_join_request', NULL, NULL, 2, NULL, 'Mia wants to join Art Enthusiasts.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'), -- cant find match to this
(9, 'group_join_request', NULL, NULL, 9, NULL, 'Amina wants to join Book Club.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(7, 'group_join_request', NULL, NULL, 15, NULL, 'Sofia wants to join Fitness Fanatics.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(4, 'group_join_request', NULL, NULL, 21, NULL, 'Liam wants to join Travel Explorers.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(13, 'group_join_request', NULL, NULL, 27, NULL, 'Mia wants to join Tech Talk.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(10, 'group_join_request', NULL, NULL, 33, NULL, 'Elias wants to join Gamers Guild.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(3, 'group_join_request', NULL, NULL, 39, NULL, 'Emma wants to join History Buffs.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
-- Event creation notifications (sent to group creators)
(1, 'event_creation', NULL, NULL, NULL, 1, 'New event: Group Hike in Nature Lovers.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(5, 'event_creation', NULL, NULL, NULL, 2, 'New event: Art Workshop in Art Enthusiasts.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(9, 'event_creation', NULL, NULL, NULL, 3, 'New event: Book Discussion in Book Club.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(2, 'event_creation', NULL, NULL, NULL, 4, 'New event: Photography Walk in Photography Hub.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(7, 'event_creation', NULL, NULL, NULL, 5, 'New event: Group Run in Fitness Fanatics.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(11, 'event_creation', NULL, NULL, NULL, 6, 'New event: Baking Class in Baking Buddies.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(4, 'event_creation', NULL, NULL, NULL, 7, 'New event: Travel Meetup in Travel Explorers.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(8, 'event_creation', NULL, NULL, NULL, 8, 'New event: Jam Session in Music Makers.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(13, 'event_creation', NULL, NULL, NULL, 9, 'New event: Tech Talk in Tech Talk.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(6, 'event_creation', NULL, NULL, NULL, 10, 'New event: Yoga Retreat in Yoga & Meditation.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(10, 'event_creation', NULL, NULL, NULL, 11, 'New event: Game Night in Gamers Guild.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(14, 'event_creation', NULL, NULL, NULL, 12, 'New event: Plant Swap in Plant Parents.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(3, 'event_creation', NULL, NULL, NULL, 13, 'New event: History Trivia in History Buffs.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(12, 'event_creation', NULL, NULL, NULL, 14, 'New event: Code Hackathon in Code Crafters.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(1, 'event_creation', NULL, NULL, NULL, 15, 'New event: Coffee Tasting in Coffee Connoisseurs.', FALSE, CURRENT_TIMESTAMP, NULL, NULL, 'enable');

-- Insert post privacy (7 columns)
INSERT INTO post_privacy (post_id, user_id, created_at, updated_at, updated_by, status) VALUES
(2, 1, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(3, 5, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(7, 2, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(8, 4, CURRENT_TIMESTAMP, NULL, NULL, 'enable'),
(10, 11, CURRENT_TIMESTAMP, NULL, NULL, 'enable');


-- Update tables with random times for created_at
UPDATE posts SET created_at = datetime('2025-04-21 13:39:00', '+' || CAST(ABS(RANDOM() % (60 * 60 * 24 * 28)) AS TEXT) || ' seconds');
UPDATE comments SET created_at = datetime('2025-04-21 13:39:00', '+' || CAST(ABS(RANDOM() % (60 * 60 * 24 * 28)) AS TEXT) || ' seconds');
UPDATE group_posts SET created_at = datetime('2025-04-21 13:39:00', '+' || CAST(ABS(RANDOM() % (60 * 60 * 24 * 28)) AS TEXT) || ' seconds');
UPDATE group_comments SET created_at = datetime('2025-04-21 13:39:00', '+' || CAST(ABS(RANDOM() % (60 * 60 * 24 * 28)) AS TEXT) || ' seconds');
UPDATE messages SET created_at = datetime('2025-04-21 13:39:00', '+' || CAST(ABS(RANDOM() % (60 * 60 * 24 * 28)) AS TEXT) || ' seconds');
UPDATE group_messages SET created_at = datetime('2025-04-21 13:39:00', '+' || CAST(ABS(RANDOM() % (60 * 60 * 24 * 28)) AS TEXT) || ' seconds');
UPDATE notifications SET created_at = datetime('2025-04-21 13:39:00', '+' || CAST(ABS(RANDOM() % (60 * 60 * 24 * 28)) AS TEXT) || ' seconds');
