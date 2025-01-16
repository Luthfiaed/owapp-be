INSERT INTO products (name, description, price, is_draft, image) VALUES 
('Product One', 'Experience ultimate comfort with this product.', 11, FALSE, 'https://fakeimg.pl/500x500/ff6600'),
('Product Two', 'This product is stylish and comfortable.', 22, FALSE, 'https://fakeimg.pl/500x500/4466cc'),
('Product Three', 'This product is made from eco-friendly materials.', 33, FALSE, 'https://fakeimg.pl/500x500/ff3399'),
('Product Four', 'Our newest product.', 44, FALSE, 'https://fakeimg.pl/500x500/ff9933'),
('Product Five', 'DRAFT PRODUCT - NOT RELEASED YET', 55, TRUE, 'https://fakeimg.pl/500x500/33cc99'),
('Product Six', 'DRAFT PRODUCT - NOT RELEASED YET', 66, TRUE, 'https://fakeimg.pl/500x500/336699');

INSERT INTO users (username, password, role) VALUES
('minyong', '$2a$12$tDHk5g8R4j0FhnPmGy9gpeYpuly8NQ3TQmyeh8.np.le46ITYs0M6', "GUEST");

INSERT INTO reviews (review, user_id, product_id) VALUES
('This product is great!', 1, 1),
('I''ve purchased this three times.', 1, 2),
('Recommended for everyone reading this!', 1, 3),
('Worth the price.', 1, 4);