CREATE TABLE IF NOT EXISTS muzzmaindb.interactions(
	id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    swiper_id INT,
    swiped_id INT,
    swipe_direction ENUM('YES', 'NO'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;