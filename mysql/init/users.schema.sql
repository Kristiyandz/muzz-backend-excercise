CREATE TABLE IF NOT EXISTS muzzmaindb.users(
	id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    email VARCHAR(320) NOT NULL,
    password_hash CHAR(60) NOT NULL,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(50) NOT NULL,
    age INT NOT NULL,
    latitude FLOAT(11, 7) NOT NULL DEFAULT 0.0,
    longitude FLOAT(11, 7) NOT NULL DEFAULT 0.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

INSERT INTO muzzmaindb.users (email, password_hash, name, gender, age, latitude, longitude) VALUES
('jonwick@example.com', '$2y$14$vAwYSgSsVPLN1OU/bYUT7e4n37EGNMRb7QjNuvCNf.qk/cIdxX7pG', 'John Wick', 'male', 25, 51.507351, -0.127758),
('tonystark@example.com', '$2y$14$RExTk1W6b5j7XSRDywY4r.pcajHZ5UZMZlTQ2XYGCPTG/YHHlQXMO', 'Tony Stark', 'male', 30, 50.715557, -3.530875),
('billieeilish@example.com', '$2y$14$RExTk1W6b5j7XSRDywY4r.pcajHZ5UZMZlTQ2XYGCPTG/YHHlQXMO', 'Billie Eilish', 'female', 20, 53.480759, -2.242631);