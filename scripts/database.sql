CREATE DATABASE IF NOT EXISTS image_hub;

USE image_hub;

CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       email VARCHAR(30) NULL,
                       password VARCHAR(255) NOT NULL,
                       role ENUM('USER', 'ADMIN') DEFAULT 'USER' NULL,
                       created_at DATETIME NOT NULL,
                       updated_at DATETIME NULL,
                       deleted_at DATETIME NULL
);

-- 어드민 더미 계정, 비밀번호: test1234
INSERT INTO users (email, password, role, created_at) VALUES
    ('admin@example.com', '$2a$10$hxuqlzruhSoYFh/5qrLEBONNDpNbsGAjRh3hmSuYwSAPMViZibLSW', 'ADMIN', NOW());

-- 유저 더미 계정, 비밀번호: test1234
INSERT INTO users (email, password, role, created_at) VALUES
    ('user1@example.com', '$2a$10$hxuqlzruhSoYFh/5qrLEBONNDpNbsGAjRh3hmSuYwSAPMViZibLSW', 'USER', NOW());
INSERT INTO users (email, password, role, created_at) VALUES
    ('user2@example.com', '$2a$10$hxuqlzruhSoYFh/5qrLEBONNDpNbsGAjRh3hmSuYwSAPMViZibLSW', 'USER', NOW());
INSERT INTO users (email, password, role, created_at) VALUES
    ('user3@example.com', '$2a$10$hxuqlzruhSoYFh/5qrLEBONNDpNbsGAjRh3hmSuYwSAPMViZibLSW', 'USER', NOW());
INSERT INTO users (email, password, role, created_at) VALUES
    ('user4@example.com', '$2a$10$hxuqlzruhSoYFh/5qrLEBONNDpNbsGAjRh3hmSuYwSAPMViZibLSW', 'USER', NOW());
INSERT INTO users (email, password, role, created_at) VALUES
    ('user5@example.com', '$2a$10$hxuqlzruhSoYFh/5qrLEBONNDpNbsGAjRh3hmSuYwSAPMViZibLSW', 'USER', NOW());


CREATE TABLE categories (
                            id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            created_at TIMESTAMP NULL,
                            updated_at TIMESTAMP NULL,
                            deleted_at TIMESTAMP NULL,
                            CONSTRAINT unique_name UNIQUE (name)
);

INSERT INTO categories (name) VALUES
                                  ('PERSON'),
                                  ('LANDSCAPE'),
                                  ('ANIMAL'),
                                  ('FOOD'),
                                  ('OTHERS');

CREATE TABLE images (
                        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                        file_name VARCHAR(255) NOT NULL,
                        file_path VARCHAR(255) NOT NULL,
                        upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
                        description TEXT NULL,
                        user_id INT NULL,
                        created_at TIMESTAMP NULL,
                        updated_at TIMESTAMP NULL,
                        deleted_at TIMESTAMP NULL,
                        thumbnail_path VARCHAR(255) NOT NULL,
                        CONSTRAINT unique_id UNIQUE (id)
);

-- 이미지 더미 데이터
INSERT INTO images (file_name, file_path, upload_date, description, user_id, created_at, updated_at, deleted_at, thumbnail_path) VALUES
                                                                                                                                     ('image1.jpg', 'uploads/1/image1.jpg', NOW(), 'Description for image1', 1, NOW(), NULL, NULL, 'uploads/1/thumb_image1.jpg'),
                                                                                                                                     ('image2.jpg', 'uploads/2/image2.jpg', NOW(), 'Description for image2', 2, NOW(), NULL, NULL, 'uploads/2/thumb_image2.jpg'),
                                                                                                                                     ('image3.jpg', 'uploads/3/image3.jpg', NOW(), 'Description for image3', 3, NOW(), NULL, NULL, 'uploads/3/thumb_image3.jpg'),
                                                                                                                                     ('image4.jpg', 'uploads/4/image4.jpg', NOW(), 'Description for image4', 4, NOW(), NULL, NULL, 'uploads/4/thumb_image4.jpg'),
                                                                                                                                     ('image5.jpg', 'uploads/5/image5.jpg', NOW(), 'Description for image5', 5, NOW(), NULL, NULL, 'uploads/5/thumb_image5.jpg');

CREATE TABLE image_categories (
                                  image_id INT NOT NULL,
                                  category_id INT NOT NULL,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
);

-- 이미지 카테고리 더미 데이터
INSERT INTO image_categories (image_id, category_id, created_at, updated_at) VALUES
                                                                                 (1, 1, NOW(), NOW()),
                                                                                 (1, 2, NOW(), NOW()), -- Image 1 has two categories
                                                                                 (2, 3, NOW(), NOW()),
                                                                                 (3, 4, NOW(), NOW()),
                                                                                 (4, 5, NOW(), NOW()),
                                                                                 (5, 1, NOW(), NOW()),
                                                                                 (5, 2, NOW(), NOW()); -- Image 5 has two categories