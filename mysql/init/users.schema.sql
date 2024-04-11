CREATE TABLE IF NOT EXISTS muzzmaindb.users(
	id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    email VARCHAR(320) NOT NULL,
    password_hash CHAR(60) NOT NULL,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(50) NOT NULL,
    age INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

INSERT INTO muzzmaindb.users (email, password_hash, name, gender, age) VALUES
('jonwick@example.com', '$2y$14$vAwYSgSsVPLN1OU/bYUT7e4n37EGNMRb7QjNuvCNf.qk/cIdxX7pG', 'John Wick', 'Male', 25),
('tonystark@example.com', '$2y$14$RExTk1W6b5j7XSRDywY4r.pcajHZ5UZMZlTQ2XYGCPTG/YHHlQXMO', 'Tony Stark', 'Male', 30),
('billieeilish@example.com', '$2y$14$RExTk1W6b5j7XSRDywY4r.pcajHZ5UZMZlTQ2XYGCPTG/YHHlQXMO', 'Billie Eilish', 'Female', 20);

-- {
--     "email": "billieeilish@example.com",
--     "password": "password-user-two"
-- }

-- {
--     "email": "jonwick@example.com",
--     "password": "password-user-one"
-- }

-- {
--     "email": "tonystark@example.com",
--     "password": "password-user-two"
-- }