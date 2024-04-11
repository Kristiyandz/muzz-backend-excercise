CREATE TABLE IF NOT EXISTS muzzmaindb.interactions(
	id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    user_id INT,
    target_user_id INT,
    choice CHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;